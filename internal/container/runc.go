package container

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var containerRuntime = "runc"

func init() {
	if runtime.GOOS == "darwin" {
		containerRuntime = "podman"
	}
}

// StartContainer
func StartContainer(name, image string, command, envVars, ports, volumes []string) error {
	args := []string{"run", "-d", "--name", name}

	// Add environment variables
	for _, env := range envVars {
		args = append(args, "-e", env)
	}

	// Add port mappings
	for _, port := range ports {
		args = append(args, "-p", port)
	}

	// Add volume bindings
	for _, volume := range volumes {
		args = append(args, "-v", volume)
	}

	// Add image and command
	args = append(args, image)
	args = append(args, command...)

	// Debug: Print the full podman command
	fmt.Printf("[DEBUG] Executing: podman %s\n", strings.Join(args, " "))

	// Execute the command
	cmd := exec.Command("podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// StopContainer
func StopContainer(containerName string, force bool) error {
	cmdArgs := []string{"stop", containerName}
	if force {
		cmdArgs = append(cmdArgs, "-f")
	}

	cmd := exec.Command(containerRuntime, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stop container failed: %v", err)
	}

	// stop container
	cmdArgs = []string{"rm", "-v", containerName}
	cmd = exec.Command(containerRuntime, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ListContainers
func ListContainers(all bool, filters []string) error {
	cmdArgs := []string{"ps"}
	if all {
		cmdArgs = append(cmdArgs, "-a")
	}
	for _, filter := range filters {
		cmdArgs = append(cmdArgs, "--filter", filter)
	}

	cmd := exec.Command(containerRuntime, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetLogs
func GetLogs(containerName string, follow bool, tailLines int) error {
	cmdArgs := []string{"logs"}
	if follow {
		cmdArgs = append(cmdArgs, "-f")
	}
	if tailLines > 0 {
		cmdArgs = append(cmdArgs, "--tail", fmt.Sprintf("%d", tailLines))
	}
	cmdArgs = append(cmdArgs, containerName)

	cmd := exec.Command(containerRuntime, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
