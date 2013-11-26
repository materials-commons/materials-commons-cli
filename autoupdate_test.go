package materials

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Keep compiler from aborting due to an used fmt package
var _ = fmt.Printf

func TestBinaryUrl(t *testing.T) {
	oss := []string{"windows", "darwin", "linux"}
	expected := map[string]string{
		"windows": "http://localhost/windows/materials.exe",
		"darwin":  "http://localhost/darwin/materials",
		"linux":   "http://localhost/linux/materials",
	}
	url := "http://localhost"
	for _, os := range oss {
		binaryUrl := binaryUrlForRuntime(url, os)
		expectedUrl, _ := expected[os]
		if binaryUrl != expectedUrl {
			t.Fatalf("Bad url %s, expected %s\n", binaryUrl, expectedUrl)
		}
	}

	expectedUrl, _ := expected[runtime.GOOS]
	binaryUrl := binaryUrl(url)
	if binaryUrl != expectedUrl {
		t.Fatalf("Bad url %s, expected %s\n", binaryUrl, expectedUrl)
	}
}

func TestDownloadNewBinary(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("test_data")))
	defer ts.Close()

	path, err := downloadNewBinary(binaryUrlForRuntime(ts.URL, "linux"))
	if err != nil {
		t.Fatalf("Unexpected error on download %s\n", err.Error())
	}

	expectedPath := filepath.Join(os.TempDir(), "materials.test")
	if path != expectedPath {
		t.Fatalf("Downloaded to unexpected name %s, expected %s\n", path, expectedPath)
	}

	// Update this sum if you change the file test_data/linux/materials
	// Computed this by doing:
	// dlChecksum := checksumFor(path)
	// fmt.Printf("checksum = %d", dlChecksum)
	expectedChecksum := uint32(1134331119)
	dlChecksum := checksumFor(path)
	if dlChecksum != expectedChecksum {
		t.Fatalf("Checksums don't match got: %d, expected: %d\n", dlChecksum, expectedChecksum)
	}
}
