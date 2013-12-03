package materials

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var _ = fmt.Printf

func TestNoConfigNoEnv(t *testing.T) {
	u, _ := NewUserFrom("test_data/noconfig")
	ConfigInitialize(u)
	if Config.Materialscommons.Api != "https://api.materialscommons.org" {
		t.Fatalf("api value incorrect %s\n", Config.Materialscommons.Api)
	}

	if Config.Materialscommons.Url != "https://materialscommons.org" {
		t.Fatalf("api value incorrect %s\n", Config.Materialscommons.Url)
	}

	if Config.Materialscommons.Download != "https://download.materialscommons.org" {
		t.Fatalf("api value incorrect %s\n", Config.Materialscommons.Download)
	}

	if Config.User.DefaultProject != "" {
		t.Fatalf("defaultProject incorrect %s\n", Config.User.DefaultProject)
	}

	expectedWebdir := filepath.Join(u.DotMaterialsPath(), "website")
	if Config.Server.Webdir != expectedWebdir {
		t.Fatalf("webdir incorrect %s\n", Config.Server.Webdir)
	}

	if Config.Server.Port != 8081 {
		t.Fatalf("port incorrect %d\n", Config.Server.Port)
	}

	if Config.Server.Address != "localhost" {
		t.Fatalf("address incorrect %s\n", Config.Server.Address)
	}

	if Config.Server.UpdateCheckInterval != 4*time.Hour {
		t.Fatalf("address incorrect %d\n", Config.Server.UpdateCheckInterval)
	}
}

func TestWithEnvSetting(t *testing.T) {
	u, _ := NewUserFrom("test_data/noconfig")
	os.Setenv("MCURL", "http://localhost")
	ConfigInitialize(u)
	if Config.Materialscommons.Url != "http://localhost" {
		t.Fatalf("url expected http://localhost, got %s\n", Config.Materialscommons.Url)
	}
}

func TestJson(t *testing.T) {
	u, _ := NewUserFrom("test_data/noconfig")
	ConfigInitialize(u)
	b, err := json.MarshalIndent(Config, "", "   ")
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout.Write(b)
}
