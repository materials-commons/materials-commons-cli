package site

import (
	"fmt"
	"github.com/materials-commons/gohandy/file"
	"github.com/materials-commons/materials"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var _ = fmt.Println

func TestDeployingSite(t *testing.T) {
	os.Remove(filepath.Join(os.TempDir(), "materials.tar.gz"))
	testData := filepath.Join("..", "test_data")
	ts := httptest.NewServer(http.FileServer(http.Dir(testData)))
	os.Setenv("MCDOWNLOADURL", ts.URL)
	u, _ := materials.NewUserFrom(testData)
	materials.ConfigInitialize(u)

	downloadedTo, err := Download()
	if err != nil {
		t.Fatalf("Unable to download materials.tar.gz %s", err.Error())
	}

	if !file.Exists(downloadedTo) {
		t.Fatalf("Download failed, unable to locate materials.tar.gz at %s\n", downloadedTo)
	}

	if !IsNew(downloadedTo) {
		t.Fatalf("Expected downloaded to be new")
	}

	if !Deploy(downloadedTo) {
		t.Fatalf("Expected site to deploy")
	}

	w := materials.Config.Server.Webdir
	checkContents(filepath.Join(w, "a"), "Hello a", 8, t)
	checkContents(filepath.Join(w, "b"), "Hello b", 8, t)
	checkContents(filepath.Join(w, "c"), "Hello c", 8, t)
}

func checkContents(fpath, expectedContents string, expectedLength int, t *testing.T) {
	contents, err := ioutil.ReadFile(fpath)

	if err != nil {
		t.Fatalf("Unable to read file %s\n", fpath)
	}

	if len(contents) != expectedLength {
		t.Fatalf("Expected length %d, got length %d\n", expectedLength, len(contents))
	}
	if strings.TrimSpace(string(contents)) != expectedContents {
		t.Fatalf("Unexpected file contents '%s'\n", string(contents))
	}

}
