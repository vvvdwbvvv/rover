package cmd

import (
	"runtime"
)

var containerRuntime string

func init() {
	if runtime.GOOS == "darwin" {
		containerRuntime = "podman"
	} else {
		containerRuntime = "runc"
	}
}
