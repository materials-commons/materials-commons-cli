package website

import ()

type WebsiteInfo struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

func Download() error {
	return nil
}
