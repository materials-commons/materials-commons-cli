// This code is from tebeka:
//https://bitbucket.org/tebeka/go-wise/src/151647b160ed257f6c05e327215d13c50ebb3856/desktop.go?at=default
package desktop

import (
	"fmt"
	"os/exec"
	"runtime"
)

var commands = map[string][]string{
	"windows": []string{"cmd", "/c", "start"},
	"darwin":  []string{"open"},
	"linux":   []string{"xdg-open"},
}

// Open calls the OS default program for uri
// e.g. Open("http://www.google.com") will open the default browser on www.google.com
func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open uri on %s platform", runtime.GOOS)
	}

	run = append(run, uri)
	cmd := exec.Command(run[0], run[1:]...)

	return cmd.Start()
}
