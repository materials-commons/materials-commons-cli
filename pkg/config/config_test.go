package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestJson(t *testing.T) {
	b, err := os.ReadFile("/home/gtarcea/.materialscommons/config.json")
	if err != nil {
		t.Fatalf("Failed reading file: %s", err)
	}

	fmt.Printf("%s\n", string(b))

	var c ConfigRemote

	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatalf("Error unmarshalling: %s", err)
	}

	fmt.Printf("%#v\n", c)
}

func TestGetRemote(t *testing.T) {
	r, err := GetRemote()
	if err != nil {
		t.Fatalf("Failure calling GetRemote: %s", err)
	}

	fmt.Printf("%#v\n", r)
}
