package autoupdate

import (
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/site"
)

// A Updater keeps track of the status of binary and website updates
// and downloads updates when they are avaiable.
type Updater struct {
	downloaded     string
	websiteUpdated bool
	binaryUpdated  bool
}

// NewUpdater creates a new Updater instance.
func NewUpdater() *Updater {
	return &Updater{
		downloaded:     "",
		websiteUpdated: false,
		binaryUpdated:  false,
	}
}

// UpdatesAvailable checks if updates are available for either the website
// or the materials binary. If updates are available it will download them.
func (u *Updater) UpdatesAvailable() bool {
	updateAvailable := false
	if downloaded, err := site.Download(); err == nil {
		if site.IsNew(u.downloaded) {
			u.downloaded = downloaded
			u.websiteUpdated = true
			updateAvailable = true
		}
	}

	if materials.Update(materials.Config.Materialscommons.Download) {
		updateAvailable = true
		u.binaryUpdated = true
	}

	return updateAvailable
}

// ApplyUpdates deploys updates that have been downloaded. If the materials
// binary has been updated then it restarts the server.
func (u *Updater) ApplyUpdates() {
	if u.websiteUpdated {
		site.Deploy(u.downloaded)
		u.websiteUpdated = false
		u.downloaded = ""
	}

	if u.binaryUpdated {
		materials.Restart()
	}
}

// WebsiteUpdate returns true if the website has been updated.
func (u *Updater) WebsiteUpdate() bool {
	return u.websiteUpdated
}

// BinaryUpdate returns true if the materials binary has been updated.
func (u *Updater) BinaryUpdate() bool {
	return u.binaryUpdated
}
