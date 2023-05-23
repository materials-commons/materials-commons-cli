package cmd

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/materials-commons/hydra/pkg/mcft/protocol"
	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/model"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/schollz/progressbar/v3"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes (uploads) added files to the server.",
	Long:  `Pushes (uploads) added files to the server.`,
	Run:   runPushCmd,
}

func runPushCmd(cmd *cobra.Command, args []string) {
	remoteStor := stor.MustLoadJsonRemoteStor()
	defaultRemote, err := remoteStor.GetDefaultRemote()
	if err != nil {
		log.Fatalf("No default remote set: %s", err)
	}

	db := mcdb.MustConnectToDB()
	projectStor := stor.NewGormProjectStor(db)
	p, err := projectStor.GetProject()
	if err != nil {
		log.Fatalf("Unable to retrieve project: %s", err)
	}

	addedFileStor := stor.NewGormAddedFileStor(db)
	bar := progressbar.DefaultBytes(-1, "Uploading")
	threadPool := pool.New().WithMaxGoroutines(config.GetMaxThreads())
	err = addedFileStor.ListPaged(func(f *model.AddedFile) error {
		threadPool.Go(func() {
			uploader := newUploader(defaultRemote.MCUrl, defaultRemote.MCAPIKey, p.ID, bar, addedFileStor)
			_ = uploader.uploadFile(f.Path)
		})
		return nil
	})

	if err != nil {
		fmt.Printf("Error during upload: %s\n", err)
	}
}

type uploader struct {
	mcurl         string
	apikey        string
	projectID     uint
	bar           *progressbar.ProgressBar
	addedFileStor stor.AddedFileStor
}

func newUploader(mcurl, apikey string, projectID uint, bar *progressbar.ProgressBar, addedFileStor stor.AddedFileStor) *uploader {
	return &uploader{mcurl: mcurl, apikey: apikey, projectID: projectID, bar: bar, addedFileStor: addedFileStor}
}

func (up *uploader) uploadFile(pathToFile string) error {
	u := url.URL{Scheme: config.GetWSScheme(), Host: up.mcurl, Path: "/ws"}
	websocket.DefaultDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Unable to connect to %s: %s", u.String(), err)
	}
	defer c.Close()

	f, err := os.Open(pathToFile)

	if err != nil {
		log.Fatalf("Unable to open %s: %s", pathToFile, err)
	}
	defer f.Close()

	var incomingReq protocol.IncomingRequestType

	if !up.authenticate(c) {
		log.Fatalf("Unable to authenticate")
	}

	incomingReq.RequestType = protocol.UploadFileReq
	if err := c.WriteJSON(incomingReq); err != nil {
		//log.Errorf("Unable to initiate upload: %s", err)
		return err
	}

	uploadToPath := mcc.ToProjectPath(pathToFile)
	// First send notice of upload
	uploadMsg := protocol.UploadFileRequest{
		Path: uploadToPath,
	}

	if err := c.WriteJSON(uploadMsg); err != nil {
		//log.Errorf("Unable to initiate upload: %s", err)
		return err
	}

	var status protocol.StatusResponse
	if err := c.ReadJSON(&status); err != nil {
		log.Printf("Unable to read upload status: %s", err)
		return err
	}

	if status.IsError {
		log.Printf("Error starting file transfer: %s", status.Status)
		return errors.New("failed to start transfer")
	}

	data := make([]byte, 32*1024*1024)
	fb := protocol.FileBlockRequest{}
	hasher := md5.New()
	for {

		n, err := f.Read(data)
		if err != nil {
			if err != io.EOF {
				//log.Errorf("Read returned error: %s", err)
				return err
			}
			break
		}

		incomingReq.RequestType = protocol.FileBlockReq
		if err := c.WriteJSON(incomingReq); err != nil {
			log.Printf("Error during upload: %s", err)
			return err
		}

		fb.Block = data[:n]
		if err := c.WriteJSON(fb); err != nil {
			//log.Errorf("WriteJSON failed: %s", err)
			return err
		}

		_, _ = io.Copy(io.MultiWriter(hasher, up.bar), bytes.NewBuffer(data[:n]))

		var status protocol.StatusResponse
		if err := c.ReadJSON(&status); err != nil {
			log.Printf("Unable to read upload status: %s", err)
			return err
		}

		if status.IsError {
			log.Printf("Error uploading file: %s", status.Status)
			return errors.New("failed upload")
		}
	}

	// compute checksum and check that they match by sending to the server
	var finishUploadRequest protocol.FinishUploadRequest
	finishUploadRequest.FileChecksum = fmt.Sprintf("%x", hasher.Sum(nil))
	finishUploadRequest.Path = uploadToPath
	incomingReq.RequestType = protocol.FinishUploadReq

	if err := c.WriteJSON(incomingReq); err != nil {
		log.Printf("Error during upload: %s", err)
		return err
	}

	if err := c.WriteJSON(&finishUploadRequest); err != nil {
		log.Printf("Error during upload: %s", err)
		return err
	}

	if err := c.ReadJSON(&status); err != nil {
		log.Printf("Unable to read upload status: %s", err)
		return err
	}

	// Uh oh the checksums didn't match
	if status.IsError {
		log.Printf("Error uploading file: %s", status.Status)
		return errors.New("failed upload - checksums didn't match")
	}

	// If we are here then the upload was successful, so delete the AddedFile entry
	if err := up.addedFileStor.RemoveByPath(pathToFile); err != nil {
		log.Printf("Error removing added file %q: %s\n", pathToFile, err)
	}

	return nil
}

func (up *uploader) authenticate(c *websocket.Conn) bool {
	var req protocol.IncomingRequestType
	req.RequestType = protocol.AuthenticateReq
	if err := c.WriteJSON(req); err != nil {
		return false
	}

	auth := protocol.AuthenticateRequest{
		APIToken:  up.apikey,
		ProjectID: int(up.projectID),
	}

	if err := c.WriteJSON(auth); err != nil {
		return false
	}

	return true
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
