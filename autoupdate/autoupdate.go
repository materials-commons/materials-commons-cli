package autoupdate

import (
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/site"
	"time"
)

// StartUpdateMonitor starts a back ground task that periodically
// checks for update to the materials command and website, downloads
// and deploys them. If the materials command is updated then the
// materials server is restarted.
func StartUpdateMonitor() {
	go updateMonitor()
}

// updateMonitor is the back ground monitor that checks for
// updates to the materials command and website. It checks
// for updates every materials.Config.UpdateCheckInterval().
func updateMonitor() {
	for {
		time.Sleep(materials.Config.UpdateCheckInterval())
		updateWebsite()
		updateBinary()
	}
}

// updateWebsite downloads and deploys new versions of the website.
func updateWebsite() {
	if downloaded, err := site.Download(); err == nil {
		if site.IsNew(downloaded) {
			site.Deploy(downloaded)
		}
	}
}

// updateBinary downloads new versions of the materials server.
// It restarts the server when there is a new version.
func updateBinary() {
	if materials.Update(materials.Config.MCDownload()) {
		materials.Restart()
	}
}
