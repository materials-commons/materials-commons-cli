package site

import (
	"fmt"
	"github.com/materials-commons/gohandy/ezhttp"
	"github.com/materials-commons/gohandy/handyfile"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/wsmaterials"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Name of archive file for the materials website.
const materialsArchive = "materials.tar.gz"

// Start starts up all the webservices and the webserver.
func Start() {
	addr := setupSite()
	fmt.Println(http.ListenAndServe(addr, nil))
}

// StartRetry attempts a number of times to try connecting to the port address.
// This is useful when the server is restarting and the old server hasn't exited yet.
func StartRetry(retryCount int) {
	addr := setupSite()
	for i := 0; i < retryCount; i++ {
		fmt.Println(http.ListenAndServe(addr, nil))
		time.Sleep(1000 * time.Millisecond)
	}
	os.Exit(1)
}

// setupSite creates all the different web services for the http server.
// It returns the address and port the http server should use.
func setupSite() string {
	container := wsmaterials.NewRegisteredServicesContainer()
	http.Handle("/", container)
	dir := http.Dir(materials.Config.WebDir())
	http.Handle("/materials/", http.StripPrefix("/materials/", http.FileServer(dir)))
	addr := fmt.Sprintf("%s:%d", materials.Config.ServerAddress(), materials.Config.ServerPort())
	return addr
}

// Download will attempt to download the materials.tar.gz file from the
// MCDownload site. If it downloads the file it will return the path to
// the downloaded file. It downloads to the OS TempDir.
func Download() (to string, err error) {
	client := ezhttp.NewClient()
	url := fmt.Sprintf("%s/%s", materials.Config.MCDownload(), materialsArchive)
	to = filepath.Join(os.TempDir(), materialsArchive)
	status, err := client.FileGet(url, to)
	switch {
	case err != nil:
		return
	case status != 200:
		return to, fmt.Errorf("Download failed with HTTP status code %d", status)
	default:
		return to, nil
	}
}

// IsNew checks if a downloaded file is newer than the current file. It does
// this by comparing the checksum of the currently downloaded file to the
// newly downloaded one. If they are different then it assumes the new one
// is more recent. If there is no currently downloaded file then by default
// the new one is more recent.
func IsNew(downloaded string) bool {
	currentArchivePath := filepath.Join(materials.Config.DotMaterials(), materialsArchive)
	if !handyfile.Exists(currentArchivePath) {
		return true
	}

	downloadedChecksum := handyfile.Checksum32(downloaded)
	currentArchiveChecksum := handyfile.Checksum32(currentArchivePath)
	if currentArchiveChecksum != downloadedChecksum {
		return true
	}
	return false
}

// Deploy attempts to deploy the new materials website archive. It will replace
// the current archive with newly downloaded one and unpack it.
func Deploy(downloaded string) bool {
	currentArchivePath := filepath.Join(materials.Config.DotMaterials(), materialsArchive)
	os.Rename(downloaded, currentArchivePath)

	tr, err := handyfile.NewTarGz(currentArchivePath)
	if err != nil {
		return false
	}

	if err := tr.Unpack(materials.Config.WebDir()); err != nil {
		return false
	}

	return true
}
