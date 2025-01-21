package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var force bool

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stops and removes all running containers",
	Long: `Stops all running containers and removes them, including their volumes.
If the --force flag is set, it will forcefully stop running containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := stopAndRemoveAllContainers(force)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func stopAndRemoveAllContainers(force bool) error {
	psCmd := exec.Command(containerRuntime, "ps", "-q")
	output, err := psCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list running containers: %v", err)
	}

	containers := strings.Fields(string(output))
	if len(containers) == 0 {
		fmt.Println("No running containers found.")
		return nil
	}

	fmt.Println("Stopping and removing the following containers:", containers)

	for _, container := range containers {
		stopArgs := []string{"stop", container}
		if force {
			stopArgs = append(stopArgs, "-f")
		}

		stopCmd := exec.Command(containerRuntime, stopArgs...)
		stopCmd.Stdout = exec.Command("echo").Stdout
		stopCmd.Stderr = exec.Command("echo").Stderr
		if err := stopCmd.Run(); err != nil {
			fmt.Printf("Failed to stop container %s: %v\n", container, err)
		} else {
			fmt.Printf("Stopped container %s\n", container)
		}

		rmCmd := exec.Command(containerRuntime, "rm", "-v", container)
		rmCmd.Stdout = exec.Command("echo").Stdout
		rmCmd.Stderr = exec.Command("echo").Stderr
		if err := rmCmd.Run(); err != nil {
			fmt.Printf("Failed to remove container %s: %v\n", container, err)
		} else {
			fmt.Printf("Removed container %s\n", container)
		}
	}

	fmt.Println("All containers stopped and removed.")
	return nil
}

func init() {
	rootCmd.AddCommand(downCmd)
	downCmd.Flags().BoolVarP(&force, "force", "f", false, "Force stop running containers")
}
