package materials

import (
	"os/user"
	"path/filepath"
	"testing"
)

func TestForCurrentUser(t *testing.T) {
	dm := dotmaterialsForCurrentUser()
	u, _ := user.Current()
	if filepath.Join(u.HomeDir, ".materials") != dm {
		t.Fatalf("Expected paths to match")
	}
}

func TestFromPath(t *testing.T) {
	dm := dotmaterialsFrom("/tmp")
	if filepath.Join("/tmp", ".materials") != dm {
		t.Fatalf("Expected paths to match")
	}
}
