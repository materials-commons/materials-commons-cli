package autoupdate

import (
	"bitbucket.org/kardianos/osext"
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/util"
	"os"
	"time"
)

var updater = NewUpdater()

// StartUpdateMonitor starts a back ground task that periodically
// checks for update to the materials command and website, downloads
// and deploys them. If the materials command is updated then the
// materials server is restarted.
func StartUpdateMonitor() {
	setLastUpdateServer()
	go updateMonitor()
}

func setLastUpdateServer() {
	binaryPath, err := osext.Executable()
	if err != nil {
		return
	}

	finfo, err := os.Stat(binaryPath)
	if err != nil {
		return
	}

	materials.Config.Server.LastServerUpdate = util.FormatTime(finfo.ModTime())
}

// updateMonitor is the back ground monitor that checks for
// updates to the materials command and website. It checks
// for updates every materials.Config.UpdateCheckInterval().
func updateMonitor() {
	for {
		materials.Config.Server.LastUpdateCheck = timeStrNow()
		materials.Config.Server.NextUpdateCheck = timeStrAfterUpdateInterval()
		if updater.UpdatesAvailable() {
			applyUpdates()
		}
		time.Sleep(materials.Config.Server.UpdateCheckInterval)
	}
}

func timeStrNow() string {
	n := time.Now()
	return util.FormatTime(n)
}

func timeStrAfterUpdateInterval() string {
	n := time.Now()
	n = n.Add(materials.Config.Server.UpdateCheckInterval)
	return util.FormatTime(n)
}

func applyUpdates() {
	updater.ApplyUpdates()
}
