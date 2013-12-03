package autoupdate

import (
	"github.com/materials-commons/materials"
	"time"
)

var updater = NewUpdater()

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
		time.Sleep(materials.Config.Server.UpdateCheckInterval)
		if updater.UpdatesAvailable() {
			updater.ApplyUpdates()
		}
	}
}
