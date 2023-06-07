package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/materials-commons/hydra/pkg/mcdb/mcmodel"
	"github.com/materials-commons/hydra/pkg/mcft/protocol"
	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcapi"
	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/project"
	"github.com/materials-commons/materials-commons-cli/pkg/stor"
	"github.com/schollz/progressbar/v3"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls (downloads) files from the server.",
	Long:  `Pulls (downloads) files from the server.`,
	Run:   runPullCmd,
}

func runPullCmd(cmd *cobra.Command, args []string) {
	remoteStor := stor.MustLoadJsonRemoteStor()
	defaultRemote, err := remoteStor.GetDefaultRemote()
	if err != nil {
		log.Fatalf("No default remote set: %s", err)
	}

	_ = defaultRemote

	db := mcdb.MustConnectToDB()
	projectStor := stor.NewGormProjectStor(db)
	p, err := projectStor.GetProject()
	if err != nil {
		log.Fatalf("Unable to retrieve project: %s", err)
	}

	_ = p

	if len(args) != 0 {
		pullSpecificFiles()
		return
	}

	pullDownloadedDirs()
}

func pullSpecificFiles() {

}

var mu sync.Mutex

func pullDownloadedDirs() {

	db := mcdb.MustConnectToDB()

	//projectStor := stor.NewGormProjectStor(db)
	//p, err := projectStor.GetProject()
	//c := mcapi.NewClient("", "")
	//files, err := c.ListDirectoryByPath(int(p.ID), "/")
	//_ = files
	//
	//if err != nil {
	//	log.Fatalf("Unable")
	//}

	// There are two steps to doing a pull. The first step is determining the status of local files
	// and identifying files that can or cannot be overwritten. This includes any directories that
	// we would pull from.
	//
	// The second step querying the server to determine what needs to be downloaded and initiating
	// the downloads. We only pull from existing directories, and don't create directories that
	// aren't already downloaded. However, a user can ask to explicitly download new directories
	// or can give exact paths to pull.

	gatherStatus(db, 1)
}

func gatherStatus(db *gorm.DB, projectID int) {
	fsHandler := newPullFileStatusHandler()
	fsStatusReciever := mcc.NewFileStatusReceiver(fsHandler.handler)
	ctx, cancel := context.WithCancel(context.Background())
	go fsStatusReciever.Run(ctx)

	changedFileHandler := func(projectPath, path string, finfo os.FileInfo) error {
		fmt.Println("changed")
		fs := &mcc.FileStatus{
			Path:        path,
			ProjectPath: projectPath,
			Status:      mcc.FileChanged,
			FInfo:       finfo,
		}

		fsStatusReciever.SendStatus(fs)
		fmt.Println("past SendStatus")
		return nil
	}

	unknownFileHandler := func(projectPath, path string, finfo os.FileInfo) error {
		fmt.Println("unknown")
		fs := &mcc.FileStatus{
			Path:        path,
			ProjectPath: projectPath,
			Status:      mcc.FileUnknown,
			FInfo:       finfo,
		}

		fsStatusReciever.SendStatus(fs)
		fmt.Println("past SendStatus")
		return nil
	}

	unchangedFileHandler := func(projectPath, path string, finfo os.FileInfo) error {
		fmt.Println("unchanged")
		fs := &mcc.FileStatus{
			Path:        path,
			ProjectPath: projectPath,
			Status:      mcc.FileKnownAndUnchanged,
			FInfo:       finfo,
		}

		fsStatusReciever.SendStatus(fs)
		fmt.Println("past SendStatus")
		return nil
	}

	projectWalker := project.NewWalker(db).
		WithChangedFileHandler(changedFileHandler).
		WithUnknownFileHandler(unknownFileHandler).
		WithUnchangedFileHandler(unchangedFileHandler).
		WithSkipUnknownDirs(false)

	fmt.Println("start walking")
	if err := projectWalker.Walk(config.GetProjectRootPath()); err != nil {
		log.Fatalf("Unable to add files: %s", err)
	}

	cancel()

	// Status of existing files is gathered.

	for _, status := range fsHandler.knownFiles {
		fmt.Println("Known File:", status.Path)
	}

	for _, status := range fsHandler.unknownFiles {
		fmt.Println("Known File:", status.Path)
	}

	for _, status := range fsHandler.changedFiles {
		fmt.Println("Known File:", status.Path)
	}

	// For known files and directories download changes. We will make parallel
	// calls to get directory contents.
	threadPool := pool.New()
	c := mcapi.NewClient("", "")
	var mu sync.Mutex
	var allFiles []mcmodel.File

	for _, path := range fsHandler.knownDirectories {
		threadPool.Go(func() {
			files, err := c.ListDirectoryByPath(projectID, path)
			if err != nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			allFiles = append(allFiles, files...)
		})
	}

	downloader := newDownloader()

	// Now that we've collected all file status and
	// all the files in the directories we care about
	// from Materials Commons, we can begin downloading
	// files that should be downloaded.
	for _, file := range allFiles {
		if fsHandler.isDownloadable(file.Path) {
			threadPool.Go(func() {
				f := file
				if err := downloader.downloadFile(f.Path); err != nil {
					log.Printf("Failure downloading file: %s\n", err)
				}
			})
		}
	}
}

type downloader struct {
	mcurl     string
	apikey    string
	projectID uint
	bar       *progressbar.ProgressBar
	fileStor  stor.FileStor
}

func newDownloader() *downloader {
	return &downloader{}
}

func (d *downloader) withMCUrl(mcurl string) *downloader {
	d.mcurl = mcurl
	return d
}

func (d *downloader) withAPIKey(apikey string) *downloader {
	d.apikey = apikey
	return d
}

func (d *downloader) withProjectID(projectID uint) *downloader {
	d.projectID = projectID
	return d
}

func (d *downloader) withProgressBar(bar *progressbar.ProgressBar) *downloader {
	d.bar = bar
	return d
}

func (d *downloader) withFileStor(fstor stor.FileStor) *downloader {
	d.fileStor = fstor
	return d
}

func (d *downloader) downloadFile(path string) error {
	u := url.URL{Scheme: config.GetWSScheme(), Host: d.mcurl, Path: "/ws"}
	websocket.DefaultDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Unable to connect to %s: %s", u.String(), err)
	}
	defer c.Close()

	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to open %s: %s", path, err)
	}

	defer f.Close()

	if !d.authenticate(c) {
		log.Fatalf("Unable to authenticate")
	}

	req := protocol.IncomingRequestType{
		RequestType: protocol.DownloadFileReq,
	}

	if err := c.WriteJSON(req); err != nil {
		return err
	}

	downloadReq := protocol.DownloadRequest{
		Path: path,
	}

	if err := c.WriteJSON(downloadReq); err != nil {
		return err
	}

	// Should add a read here to send size and checksum
	// and then just keep on writing file blocks until
	// no more to send.

	var incomingRequest protocol.IncomingRequestType
	for {
		if err := c.ReadJSON(&incomingRequest); err != nil {
			break
		}

		switch incomingRequest.RequestType {
		case protocol.FileBlockReq:
			err = d.writeFileBlock(f, c)
		case protocol.FinishDownloadReq:
			err = d.finishDownload(c)
		}
	}

	return nil
}

func (d *downloader) authenticate(c *websocket.Conn) bool {
	var req protocol.IncomingRequestType
	req.RequestType = protocol.AuthenticateReq
	if err := c.WriteJSON(req); err != nil {
		return false
	}

	auth := protocol.AuthenticateRequest{
		APIToken:  d.apikey,
		ProjectID: int(d.projectID),
	}

	if err := c.WriteJSON(auth); err != nil {
		return false
	}

	return true
}

func (d *downloader) writeFileBlock(f *os.File, c *websocket.Conn) error {
	return nil
}

func (d *downloader) finishDownload(c *websocket.Conn) error {
	return nil
}

type pullFileStatusHandler struct {
	knownFiles       map[string]*mcc.FileStatus
	unknownFiles     map[string]*mcc.FileStatus
	changedFiles     map[string]*mcc.FileStatus
	knownDirectories []string
}

func newPullFileStatusHandler() *pullFileStatusHandler {
	h := &pullFileStatusHandler{
		knownFiles:   make(map[string]*mcc.FileStatus),
		unknownFiles: make(map[string]*mcc.FileStatus),
		changedFiles: make(map[string]*mcc.FileStatus),
	}

	h.knownDirectories = append(h.knownDirectories, "/")

	return h
}

func (h *pullFileStatusHandler) handler(fs *mcc.FileStatus) {
	switch fs.Status {
	case mcc.FileUnknown:
		h.unknownFiles[fs.Path] = fs
	case mcc.FileChanged:
		h.changedFiles[fs.Path] = fs
	default:
		// File is known
		h.knownFiles[fs.Path] = fs
	}
}

func (h *pullFileStatusHandler) isDownloadable(path string) bool {
	if _, ok := h.unknownFiles[path]; ok {
		// File is unknown, we can't overwrite it
		return false
	}

	if _, ok := h.changedFiles[path]; ok {
		// Known file that has changed
		return false
	}

	// If we are here then the file is either updating an
	// existing file that has already been uploaded, or
	// the file doesn't exist locally and thus can be
	// downloaded.

	return true
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().BoolP("recursive", "r", false, "Recursively download all files and directories")
	pullCmd.Flags().BoolP("new", "n", false, "Download new directories and files")
	pullCmd.Flags().BoolP("dry-run", "d", false, "Do not actually run downloads, just show what would happen")

	// Is this the same as dry-run?
	pullCmd.Flags().BoolP("differences", "x", false, "Show differences between existing and what is on server")
}
