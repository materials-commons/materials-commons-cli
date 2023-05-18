package config

type Remote struct {
	MCUrl    string `json:"mcurl"`
	EMail    string `json:"email"`
	MCAPIKey string `json:"mcapikey"`
}

type ConfigRemote struct {
	DefaultRemote Remote   `json:"default_remote"`
	Remotes       []Remote `json:"remotes"`
}
