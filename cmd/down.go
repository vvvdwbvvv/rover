package cmd

import (
	"fmt"
	"github.com/vvvdwbvvv/rover/internal/config"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var force bool

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stops and removes all containers defined in rover-compose.yaml",
	Long: `Stops and removes containers defined in rover-compose.yaml (or toml/json),
ensuring that they are stopped in the correct order based on depends_on.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := stopAndRemoveContainers(force)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func stopAndRemoveContainers(force bool) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load rover-compose: %v", err)
	}

	stopOrder, err := config.GetServiceShutdownOrder(cfg.Services)
	if err != nil {
		return fmt.Errorf("failed to calculate shutdown order: %v", err)
	}

	if len(stopOrder) == 0 {
		fmt.Println("No containers found in rover-compose. Nothing to stop.")
		return nil
	}

	fmt.Println("Stopping containers in reverse order:", stopOrder)

	for _, containerName := range stopOrder {
		if !isContainerRunning(containerName) {
			fmt.Printf("⚠️  Container %s is not running, skipping...\n", containerName)
			continue
		}

		stopArgs := []string{"stop", containerName}
		if force {
			stopArgs = append(stopArgs, "-f")
		}

		stopCmd := exec.Command(containerRuntime, stopArgs...)
		if err := stopCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to stop %s: %v\n", containerName, err)
			continue
		}
		fmt.Printf("✅ Stopped container: %s\n", containerName)

		rmCmd := exec.Command(containerRuntime, "rm", "-v", containerName)
		if err := rmCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to remove %s: %v\n", containerName, err)
		} else {
			fmt.Printf("🗑️  Removed container: %s\n", containerName)
		}
	}

	fmt.Println("✅ All containers stopped and removed.")
	return nil
}

func isContainerRunning(containerName string) bool {
	psCmd := exec.Command(containerRuntime, "ps", "--format", "{{.Names}}")
	output, err := psCmd.Output()
	if err != nil {
		fmt.Printf("⚠️  Failed to check running containers: %v\n", err)
		return false
	}

	containers := strings.Fields(string(output))
	for _, name := range containers {
		if name == containerName {
			return true
		}
	}
	return false
}

func GetServiceShutdownOrder(services map[string]config.Service) ([]string, error) {
	startOrder, err := config.GetServiceStartupOrder(services)
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(startOrder)-1; i < j; i, j = i+1, j-1 {
		startOrder[i], startOrder[j] = startOrder[j], startOrder[i]
	}
	return startOrder, nil
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().BoolVarP(&force, "force", "f", false, "Force stop running containers")
}
