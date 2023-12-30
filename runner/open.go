package runner

import "runtime"

func GetOpen() (string, []string) {
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
