package open

import (
	"fmt"
	"os/exec"
	"runtime"
)

func getCommand() (string, []string) {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "windows":
		cName = "cmd"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		cName = "xdg-open"
	}

	return cName, cParams
}

func Open(address string) {
	cName, cParams := getCommand()
	cParams = append(cParams, fmt.Sprintf("http://%s", address))
	out, err := exec.Command(cName, cParams...).CombinedOutput()

	if err != nil {
		fmt.Println(err, string(out))
	}
}
