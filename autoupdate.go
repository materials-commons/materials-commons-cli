package materials

import (
	"bitbucket.org/kardianos/osext"
	"fmt"
	"os"
	"os/exec"
	//"github.com/materials-commons/gohandy/ezhttp"
	"hash/crc32"
	"io/ioutil"
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

func Download(url string) {
	path, _ := osext.Executable()
	myChecksum := checksumFor(path)
	path = downloadNewBinary(url)
	fmt.Printf("%d", myChecksum)
}

func checksumFor(path string) uint32 {
	file, _ := os.Open(path)
	defer file.Close()
	c := crc32.NewIEEE()
	bytes, _ := ioutil.ReadAll(file)
	withcrc := c.Sum(bytes)
	return crc32.ChecksumIEEE(withcrc)
}

func downloadNewBinary(url string) string {
	return ""
}
