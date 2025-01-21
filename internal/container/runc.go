package container

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var containerRuntime = "runc"

func init() {
	if runtime.GOOS == "darwin" {
		containerRuntime = "podman"
	}
}

// StartContainer
func StartContainer(containerName, image string, command, envVars, ports, volumes []string) error {
	cmdArgs := []string{"run", "-d", "--name", containerName}

	// Add environment variables
	for _, env := range envVars {
		cmdArgs = append(cmdArgs, "-e", env)
	}

	// Add port mappings
	for _, port := range ports {
		cmdArgs = append(cmdArgs, "-p", port)
	}

	// Add volume mounts
	for _, volume := range volumes {
		cmdArgs = append(cmdArgs, "-v", volume)
	}

	cmdArgs = append(cmdArgs, image)
	cmdArgs = append(cmdArgs, command...)

	cmd := exec.Command(containerRuntime, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// StopContainer
func StopContainer(containerName string) error {
	cmd := exec.Command(containerRuntime, "stop", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Stop container failed: %v", err)
	}

	// DeleteContainer
	cmd = exec.Command(containerRuntime, "rm", containerName)
	return cmd.Run()
}

// ListContainers
func ListContainers() error {
	cmd := exec.Command(containerRuntime, "ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetLogs
func GetLogs(containerName string) error {
	cmd := exec.Command(containerRuntime, "logs", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
