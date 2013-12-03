package autoupdate

import (
	"github.com/materials-commons/materials"
	"github.com/materials-commons/materials/site"
)

type Updater struct {
	downloaded     string
	websiteUpdated bool
	binaryUpdated  bool
}

func NewUpdater() *Updater {
	return &Updater{
		downloaded:     "",
		websiteUpdated: false,
		binaryUpdated:  false,
	}
}

func (u *Updater) UpdatesAvailable() bool {
	updateAvailable := false
	if downloaded, err := site.Download(); err == nil {
		if site.IsNew(u.downloaded) {
			u.downloaded = downloaded
			u.websiteUpdated = true
			updateAvailable = true
		}
	}

	if materials.Update(materials.Config.MCDownload()) {
		updateAvailable = true
		u.binaryUpdated = true
	}

	return updateAvailable

}

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

func (u *Updater) WebsiteUpdate() bool {
	return u.websiteUpdated
}

func (u *Updater) BinaryUpdate() bool {
	return u.binaryUpdated
}
