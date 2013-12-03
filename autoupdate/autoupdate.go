package autoupdate

import (
	"fmt"
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
		materials.Config.Server.LastUpdateCheck = timeStrNow()
		materials.Config.Server.NextUpdateCheck = timeStrAfterUpdateInterval()
		if updater.UpdatesAvailable() {
			updater.ApplyUpdates()
		}
		time.Sleep(materials.Config.Server.UpdateCheckInterval)
	}
}

func timeStrNow() string {
	n := time.Now()
	return formatTime(n)
}

func timeStrAfterUpdateInterval() string {
	n := time.Now()
	n = n.Add(materials.Config.Server.UpdateCheckInterval)
	return formatTime(n)
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
