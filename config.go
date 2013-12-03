package materials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type MaterialscommonsConfig struct {
	Api      string
	Url      string
	Download string
}

type ServerConfig struct {
	Port                uint
	Address             string
	Webdir              string
	UpdateCheckInterval time.Duration
}

type UserConfig struct {
	*User
	DefaultProject string
}

type ConfigSettings struct {
	Materialscommons MaterialscommonsConfig
	Server           ServerConfig
	User             UserConfig
}

type configFile map[string]interface{}

var defaultSettings = map[string]interface{}{
	"server_address":        "localhost",
	"server_port":           uint(8081),
	"update_check_interval": 4 * time.Hour,
	"MCURL":                 "https://materialscommons.org",
	"MCAPIURL":              "https://api.materialscommons.org",
	"MCDOWNLOADURL":         "https://download.materialscommons.org",
}

var Config ConfigSettings

//*********************************************************
// Create on Initialize() for the materials package
// that encompasses all the other initialization, such
// as projects, and .user
//*********************************************************
func ConfigInitialize(user *User) {
	Config.User.User = user
	Config.setConfigOverrides()
}

func (c *ConfigSettings) setConfigOverrides() {
	configFromFile, _ := readConfigFile(c.User.DotMaterialsPath())
	c.Server.Port = getConfigUint("server_port", "MATERIALS_PORT", configFromFile)
	c.Server.Address = getConfigStr("server_address", "MATERIALS_ADDRESS", configFromFile)
	updateCheckInterval := getConfigDuration("update_check_interval", "MATERIALS_UPDATE_CHECK_INTERVAL", configFromFile)
	c.Server.UpdateCheckInterval = updateCheckInterval
	c.Materialscommons.Api = getDefaultedConfigStr("MCAPIURL", "MCAPIURL")
	c.Materialscommons.Url = getDefaultedConfigStr("MCURL", "MCURL")
	c.Materialscommons.Download = getDefaultedConfigStr("MCDOWNLOADURL", "MCDOWNLOADURL")

	webdir := os.Getenv("MATERIALS_WEBDIR")
	if webdir == "" {
		webdir = filepath.Join(c.User.DotMaterialsPath(), "website")
	}

	c.Server.Webdir = webdir

	cf := configFromFile
	defaultProject, ok := cf["default_project"].(string)
	if ok {
		c.User.DefaultProject = defaultProject
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

func getConfigDuration(jsonName, envName string, c configFile) time.Duration {
	envVal, err := strconv.ParseUint(os.Getenv(envName), 0, 32)
	jsonVal, ok := c[jsonName].(time.Duration)

	switch {
	case err == nil:
		return time.Duration(envVal)
	case ok && jsonVal != 0:
		return jsonVal
	default:
		val, _ := defaultSettings[jsonName].(time.Duration)
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

func readConfigFile(dotmaterialsPath string) (cf configFile, err error) {
	configPath := configPath(dotmaterialsPath)
	bytes, err := ioutil.ReadFile(configPath)
	var config configFile
	cf = config

	if err != nil {
		return cf, err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return cf, err
	}

	return config, nil
}

func configPath(path string) string {
	return filepath.Join(path, ".config")
}

func writeConfigFile(config configFile, dotmaterialsPath string) error {
	return nil
}

func (c ConfigSettings) MCUrl() string {
	return c.Materialscommons.Url
}

func (c ConfigSettings) MCApi() string {
	return c.Materialscommons.Api
}

func (c ConfigSettings) MCDownload() string {
	return c.Materialscommons.Download
}

func (c ConfigSettings) ServerPort() uint {
	return c.Server.Port
}

func (c *ConfigSettings) SetServerPort(port uint) {
	c.Server.Port = port
}

func (c ConfigSettings) ServerAddress() string {
	return c.Server.Address
}

func (c *ConfigSettings) SetServerAddress(address string) {
	c.Server.Address = address
}

func (c ConfigSettings) UpdateCheckInterval() time.Duration {
	return c.Server.UpdateCheckInterval
}

func (c ConfigSettings) WebDir() string {
	return c.Server.Webdir
}

func (c ConfigSettings) DotMaterials() string {
	return c.User.DotMaterialsPath()
}

func (c ConfigSettings) DefaultProject() string {
	return c.User.DefaultProject
}

func (c ConfigSettings) Json() []byte {
	b, err := json.MarshalIndent(&c, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	return b
}

// Constructs the url to access an api service. Includes the
// apikey. Prepends a "/" if needed.
func (c ConfigSettings) ApiUrlPath(service string) string {
	if string(service[0]) != "/" {
		service = "/" + service
	}
	uri := c.Materialscommons.Api + service + "?apikey=" + c.User.Apikey
	return uri
}
