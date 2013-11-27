package materials

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type materialscommonsConfig struct {
	api      string
	url      string
	download string
}

type serverConfig struct {
	port           uint
	address        string
	webdir         string
	updateInterval time.Duration
}

type userConfig struct {
	*User
	defaultProject string
}

type config struct {
	materialscommons materialscommonsConfig
	server           serverConfig
	user             userConfig
}

type configFile map[string]interface{}

var defaultSettings = map[string]interface{}{
	"server_address":        "localhost",
	"server_port":           uint(8081),
	"update_check_interval": uint(4 * time.Hour),
	"MCURL":                 "https://materialscommons.org",
	"MCAPIURL":              "https://api.materialscommons.org",
	"MCDOWNLOADURL":         "https://download.materialscommons.org",
}

var Config config

//*********************************************************
// Create on Initialize() for the materials package
// that encompasses all the other initialization, such
// as projects, and .user
//*********************************************************
func ConfigInitialize(user *User) {
	Config.user.User = user
	Config.setConfigOverrides()
}

func (c *config) setConfigOverrides() {
	configFromFile, _ := readConfigFile(c.user.DotMaterialsPath())
	c.server.port = getConfigUint("server_port", "MATERIALS_PORT", *configFromFile)
	c.server.address = getConfigStr("server_address", "MATERIALS_ADDRESS", *configFromFile)
	updateInterval := time.Duration(getConfigUint("update_check_interval",
		"MATERIALS_UPDATE_CHECK_INTERVAL", *configFromFile))
	c.server.updateInterval = updateInterval
	c.materialscommons.api = getDefaultedConfigStr("MCAPIURL", "MCAPIURL")
	c.materialscommons.url = getDefaultedConfigStr("MCURL", "MCURL")
	c.materialscommons.download = getDefaultedConfigStr("MCDOWNLOADURL", "MCDOWNLOADURL")

	webdir := os.Getenv("MATERIALS_WEBDIR")
	if webdir == "" {
		webdir = filepath.Join(c.user.DotMaterialsPath(), "website")
	}

	c.server.webdir = webdir

	cf := *configFromFile
	defaultProject, ok := cf["default_project"].(string)
	if ok {
		c.user.defaultProject = defaultProject
	}
}

func getConfigUint(jsonName, envName string, c configFile) uint {
	envVal, err := strconv.ParseUint(os.Getenv(envName), 0, 32)
	jsonVal, ok := c[jsonName].(uint)

	switch {
	case err == nil:
		return uint(envVal)
	case ok && jsonVal != 0:
		return jsonVal
	default:
		val, _ := defaultSettings[jsonName].(uint)
		return val
	}
}

func getConfigStr(jsonName, envName string, c configFile) string {
	envVal := os.Getenv(envName)
	jsonVal, ok := c[jsonName].(string)

	switch {
	case envVal != "":
		return envVal
	case ok && jsonVal != "":
		return jsonVal
	default:
		val, _ := defaultSettings[jsonName].(string)
		return val
	}
}

func getDefaultedConfigStr(envName, settingsName string) string {
	envVal := os.Getenv(envName)
	if envVal == "" {
		return defaultSettings[settingsName].(string)
	}

	return envVal
}

func readConfigFile(dotmaterialsPath string) (cf *configFile, err error) {
	configPath := configPath(dotmaterialsPath)
	bytes, err := ioutil.ReadFile(configPath)
	var config configFile
	cf = &config

	if err != nil {
		return cf, err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return cf, err
	}

	return &config, nil
}

func configPath(path string) string {
	return filepath.Join(path, ".config")
}

func writeConfigFile(config configFile, dotmaterialsPath string) error {
	return nil
}

func (c config) MCUrl() string {
	return c.materialscommons.url
}

func (c config) MCApi() string {
	return c.materialscommons.api
}

func (c config) MCDownload() string {
	return c.materialscommons.download
}

func (c config) ServerPort() uint {
	return c.server.port
}

func (c config) ServerAddress() string {
	return c.server.address
}

func (c config) WebDir() string {
	return c.server.webdir
}

func (c config) DotMaterials() string {
	return c.user.DotMaterialsPath()
}

func (c config) DefaultProject() string {
	return c.user.defaultProject
}

// Constructs the url to access an api service. Includes the
// apikey. Prepends a "/" if needed.
func (c config) ApiUrlPath(service string) string {
	if string(service[0]) != "/" {
		service = "/" + service
	}
	uri := c.materialscommons.api + service + "?apikey=" + c.user.Apikey
	return uri
}
