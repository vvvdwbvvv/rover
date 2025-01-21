package cmd

import (
	"fmt"
	"github.com/vvvdwbvvv/rover/internal/config"
	"github.com/vvvdwbvvv/rover/internal/container"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [container-name]",
	Short: "Launch containers from rover-compose.yaml (or specific container)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Parse rover-compose.yaml / toml / json
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Step 2: Compute startup order
		startOrder, err := config.GetServiceStartupOrder(cfg.Services)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Dependency error: %v\n", err)
			os.Exit(1)
		}

		// Step 3: Handle CLI arguments
		extraEnv, _ := cmd.Flags().GetStringArray("env")
		extraPorts, _ := cmd.Flags().GetStringArray("port")
		extraVolumes, _ := cmd.Flags().GetStringArray("volume")
		customCommand, _ := cmd.Flags().GetStringArray("command")

		// Check if a specific container is targeted
		var containerName string
		if len(args) > 0 {
			containerName = args[0]
			if _, exists := cfg.Services[containerName]; !exists {
				fmt.Fprintf(os.Stderr, "❌ Container '%s' not found in rover-compose\n", containerName)
				os.Exit(1)
			}
			startOrder = []string{containerName}
		}

		// Step 4: Start containers in order
		fmt.Println("🚀 Starting containers in order:", startOrder)
		failedContainers := []string{}

		for _, serviceName := range startOrder {
			service := cfg.Services[serviceName]

			envVars := append(service.EnvironmentToArray(), extraEnv...)
			ports := append(service.Ports, extraPorts...)
			volumes := append(service.Volumes, extraVolumes...)
			command := customCommand
			if len(command) == 0 {
				command = service.Command
			}

			// Check if the container is already running
			if isContainerRunning(service.Name) {
				fmt.Printf("🔄 Container %s is already running. Skipping...\n", service.Name)
				continue
			}

			// Stop any conflicting containers
			if stopConflictingContainer(service.Name) {
				fmt.Printf("🛑 Stopped conflicting container: %s\n", service.Name)
			}

			fmt.Printf("[+] Starting %s (%s)...\n", service.Name, service.Image)

			// Start the container with retries
			if err := StartContainerWithRetry(service, 3); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to start %s after retries: %v\n", service.Name, err)
				failedContainers = append(failedContainers, service.Name)
				continue
			}

			// Verify the container
			if err := PostLaunchVerify(service); err != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Container %s started but verification failed: %v\n", service.Name, err)
				failedContainers = append(failedContainers, service.Name)
				continue
			}

			fmt.Printf("✅ Successfully started and verified container: %s\n", service.Name)
		}

		// Step 5: Display summary
		if len(failedContainers) > 0 {
			fmt.Fprintf(os.Stderr, "❌ The following containers failed to start: %v\n", failedContainers)
			os.Exit(1)
		}

		fmt.Println("✅ All containers started successfully.")
	},
}

func (s config.Service) EnvironmentToArray() []string {
	env := []string{}
	for key, value := range s.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}

func isContainerRunning(containerName string) bool {
	cmd := exec.Command("podman", "ps", "--filter", fmt.Sprintf("name=%s", containerName), "--format", "{{.ID}}")
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) != ""
}

func stopConflictingContainer(containerName string) bool {
	cmd := exec.Command("podman", "ps", "--filter", fmt.Sprintf("name=%s", containerName), "--format", "{{.ID}}")
	output, _ := cmd.Output()
	containerID := strings.TrimSpace(string(output))
	if containerID != "" {
		stopCmd := exec.Command("podman", "stop", containerID)
		if err := stopCmd.Run(); err == nil {
			return true
		}
	}
	return false
}

func PostLaunchVerify(service config.Service) error {
	cmd := exec.Command("podman", "ps", "--filter", fmt.Sprintf("name=%s", service.Name), "--format", "{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to verify container %s: %v", service.Name, err)
	}

	status := strings.TrimSpace(string(output))
	if !strings.HasPrefix(status, "Up") {
		return fmt.Errorf("container %s is not running, current status: %s", service.Name, status)
	}

	return nil
}

func StartContainerWithRetry(service config.Service, retries int) error {
	for attempt := 1; attempt <= retries; attempt++ {
		fmt.Printf("[DEBUG] [Attempt %d] Trying to start container %s...\n", attempt, service.Name)

		// Log container details for debugging
		fmt.Printf("[DEBUG] Image: %s, Ports: %v, Volumes: %v, Env: %v, Command: %v\n",
			service.Image, service.Ports, service.Volumes, service.EnvironmentToArray(), service.Command)

		// Attempt to start the container
		err := container.StartContainer(service.Name, service.Image, service.Command, service.EnvironmentToArray(), service.Ports, service.Volumes)
		if err == nil {
			fmt.Printf("[DEBUG] [Attempt %d] Successfully started container: %s\n", attempt, service.Name)
			return nil // Success
		}

		// Log the failure
		fmt.Printf("[DEBUG] [Attempt %d] Failed to start container %s: %v\n", attempt, service.Name, err)
	}
	return fmt.Errorf("[DEBUG] All %d attempts to start container %s failed", retries, service.Name)
}

func init() {
	upCmd.Flags().StringArray("env", nil, "Extra environment variables (e.g., --env KEY=VALUE)")
	upCmd.Flags().StringArray("port", nil, "Extra port mappings (e.g., --port 8080:80)")
	upCmd.Flags().StringArray("volume", nil, "Extra volume bindings (e.g., --volume /host:/container)")
	upCmd.Flags().StringArray("command", nil, "Custom command to run inside the container (e.g., --command sh -c 'echo hello')")
}
