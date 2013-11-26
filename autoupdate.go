package materials

import (
	"bitbucket.org/kardianos/osext"
	"fmt"
	"os"
	"os/exec"
)

func Restart() {
	commandPath, err := osext.Executable()
	if err != nil {
		fmt.Printf("Unable to determine my executable path: %s\n", err.Error())
		return
	}

	command := exec.Command("nohup", commandPath, "--server", "--retry=10")
	command.Start()
	os.Exit(0)
}
