package site

import (
	"fmt"
	"github.com/materials-commons/gohandy/ezhttp"
	"github.com/materials-commons/gohandy/file"
	"github.com/materials-commons/materials"
	"os"
	"path/filepath"
)

// Name of archive file for the materials website.
const materialsArchive = "materials.tar.gz"

// Download will attempt to download the materials.tar.gz file from the
// MCDownload site. If it downloads the file it will return the path to
// the downloaded file. It downloads to the OS TempDir.
func Download() (to string, err error) {
	client := ezhttp.NewClient()
	url := fmt.Sprintf("%s/%s", materials.Config.Materialscommons.Download, materialsArchive)
	to = filepath.Join(os.TempDir(), materialsArchive)
	status, err := client.FileGet(url, to)
	switch {
	case err != nil:
		return
	case status != 200:
		return to, fmt.Errorf("download failed with HTTP status code %d", status)
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
	currentArchivePath := filepath.Join(materials.Config.User.DotMaterialsPath(), materialsArchive)
	if !file.Exists(currentArchivePath) {
		return true
	}

	downloadedChecksum := file.Checksum32(downloaded)
	currentArchiveChecksum := file.Checksum32(currentArchivePath)
	if currentArchiveChecksum != downloadedChecksum {
		return true
	}
	return false
}

// Deploy attempts to deploy the new materials website archive. It will replace
// the current archive with newly downloaded one and unpack it.
func Deploy(downloaded string) bool {
	currentArchivePath := filepath.Join(materials.Config.User.DotMaterialsPath(), materialsArchive)

	err := moveFile(downloaded, currentArchivePath)
	if err != nil {
		return false
	}

	tr, err := file.NewTarGz(currentArchivePath)
	if err != nil {
		return false
	}

	if err := tr.Unpack(materials.Config.Server.Webdir); err != nil {
		return false
	}
	return true
}

func moveFile(src, dest string) error {
	err := file.Copy(src, dest)
	if err != nil {
		return err
	}
	os.Remove(src)
	return nil
}
