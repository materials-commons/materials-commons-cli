package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/materials-commons/materials-commons-cli/pkg/config"
	"github.com/materials-commons/materials-commons-cli/pkg/mcc"
	"github.com/materials-commons/materials-commons-cli/pkg/mcdb"
	"github.com/materials-commons/materials-commons-cli/pkg/project"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Gets the status of local files and identifies what should be uploaded.",
	Long:  `Gets the status of local files and identifies what should be uploaded.`,
	Run:   runStatusCmd,
}

// runStatusCmd runs the status command. It traverses the project tree and prints out
// the status of files if they have changed or are unknown. For directories, it will
// only print out an unknown directory, it will not print out the status of the files
// in the directory because they are, by definition, unknown as well.
func runStatusCmd(cmd *cobra.Command, args []string) {
	db := mcdb.MustConnectToDB()

	// Create a statusReceiver to send file status messages to.
	sr := newStatusReceiver()

	// Create a status walker with callbacks to handle unknown and changed files.
	sw := newStatusWalkerState(sr)

	// The statusReceiver run thread will run in the background. The context provides a
	// cancel function we can use to tell it to exit.
	ctx, cancel := context.WithCancel(context.Background())
	go sr.run(ctx) // start the statusReceiver thread

	// Walk the project determining and printing file status.
	projectWalker := project.NewWalker(db).
		WithChangedFileHandler(sw.changedFileHandler).
		WithUnknownFileHandler(sw.unknownFileHandler)
	if err := projectWalker.Walk(config.GetProjectRootPath()); err != nil {
		log.Fatalf("Unable to get status: %s", err)
	}

	// Tell the status receiver go routine to exit
	cancel()
}

// fileStatus represents the status of a file. It is sent along a channel to the statusReceiver
type fileStatus struct {
	Path        string // The file system path
	ProjectPath string // The path in the project where project root is /
	Status      string // StatusUnknown or StatusChanged
	FType       string // FTypeDirectory or FTypeFile
}

// statusReceiver is a background thread that receives the status for files and prints
// them out to the user. This increases the perceived speed of status checks as the
// threads that are determining file status can send the status to the receiver which
// will immediately print it.
type statusReceiver struct {
	in chan *fileStatus
}

func newStatusReceiver() *statusReceiver {
	return &statusReceiver{
		in: make(chan *fileStatus, 1),
	}
}

// run should be started with go. It will run in the background receiving fileStatus messages
// and printing them to standard out.
func (r *statusReceiver) run(c context.Context) {
	for {
		select {
		case fstatus := <-r.in:
			fmt.Printf("%s %s (%s)\n", fstatus.Status, fstatus.Path, fstatus.FType)
		case <-c.Done():
			return
		case <-time.After(10 * time.Second):
		}
	}
}

// sendStatus will send a fileStatus to the status receiver (run method).
func (r *statusReceiver) sendStatus(fstatus *fileStatus) {
	r.in <- fstatus
}

// statusWalkerState is used to handle the different status's for files that the Walker finds.
// It has two callbacks (changedFileHandler and unknownFileHandler) that are called by the project walker.
// It constructs and sends the status to the statusReceiver.
type statusWalkerState struct {
	sr *statusReceiver
}

func newStatusWalkerState(sr *statusReceiver) *statusWalkerState {
	return &statusWalkerState{sr: sr}
}

// changedFileHandler is called when a file's mtime has been changed. It sends a fileStatus message
// with status StatusChanged.
func (w *statusWalkerState) changedFileHandler(projectPath, path string, finfo os.FileInfo) error {
	w.sendStatus(projectPath, path, mcc.FileChanged, finfo)
	return nil
}

// unknownFileHandler is called when a file is not known (in the project local database). It sends
// a fileStatus with StatusUnknown.
func (w *statusWalkerState) unknownFileHandler(projectPath, path string, finfo os.FileInfo) error {
	w.sendStatus(projectPath, path, mcc.FileUnknown, finfo)
	return nil
}

// sendStatus constructs the fileStatus, determines the file type and sends the message.
func (w *statusWalkerState) sendStatus(projectPath, path, status string, finfo os.FileInfo) {
	ftype := mcc.FTypeFile

	if finfo.IsDir() {
		ftype = mcc.FTypeDirectory
	}

	fstatus := &fileStatus{
		Path:        path,
		ProjectPath: projectPath,
		Status:      status,
		FType:       ftype,
	}

	w.sr.sendStatus(fstatus)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
