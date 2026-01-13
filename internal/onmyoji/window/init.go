package window

import "runtime"

var SystemOS string

func init() {
	SystemOS = detectOS()
	InitPlatform(SystemOS)
}

func detectOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "mac"
	case "windows":
		return "windows"
	default:
		return "unknown"
	}
}
