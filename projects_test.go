package materials

import (
	"testing"
)

func TestNonExistantUser(t *testing.T) {
	username := "no-such-user-xxx"
	_, err := ProjectsForUser(username)
	if err == nil {
		t.Errorf("Should not have found user '%s'\n", username)
	}
}

//func Test
