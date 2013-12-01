package materials

import (
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
	if Config.materialscommons.api != "https://api.materialscommons.org" {
		t.Fatalf("api value incorrect %s\n", Config.materialscommons.api)
	}

	if Config.materialscommons.url != "https://materialscommons.org" {
		t.Fatalf("api value incorrect %s\n", Config.materialscommons.url)
	}

	if Config.materialscommons.download != "https://download.materialscommons.org" {
		t.Fatalf("api value incorrect %s\n", Config.materialscommons.download)
	}

	if Config.user.defaultProject != "" {
		t.Fatalf("defaultProject incorrect %s\n", Config.user.defaultProject)
	}

	expectedWebdir := filepath.Join(u.DotMaterialsPath(), "website")
	if Config.server.webdir != expectedWebdir {
		t.Fatalf("webdir incorrect %s\n", Config.server.webdir)
	}

	if Config.server.port != 8081 {
		t.Fatalf("port incorrect %d\n", Config.server.port)
	}

	if Config.server.address != "localhost" {
		t.Fatalf("address incorrect %s\n", Config.server.address)
	}

	if Config.server.updateCheckInterval != 4*time.Hour {
		t.Fatalf("address incorrect %d\n", Config.server.updateCheckInterval)
	}
}

func TestWithEnvSetting(t *testing.T) {
	u, _ := NewUserFrom("test_data/noconfig")
	os.Setenv("MCURL", "http://localhost")
	ConfigInitialize(u)
	if Config.materialscommons.url != "http://localhost" {
		t.Fatalf("url expected http://localhost, got %s\n", Config.materialscommons.url)
	}
}
