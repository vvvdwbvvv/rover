package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vvvdwbvvv/rover/internal/container"
)

var upCmd = &cobra.Command{
	Use:   "up [container-name] [image]",
	Short: "Launch a container",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		containerName := args[0]
		image := args[1]

		// Parse CLI flags
		command, _ := cmd.Flags().GetStringArray("command")
		envVars, _ := cmd.Flags().GetStringArray("env")
		ports, _ := cmd.Flags().GetStringArray("port")
		volumes, _ := cmd.Flags().GetStringArray("volume")

		fmt.Printf("🚀 Launching container: %s with image: %s\n", containerName, image)

		// Launch container
		if err := container.StartContainer(containerName, image, command, envVars, ports, volumes); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Launch failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Container started successfully")
	},
}

func init() {
	upCmd.Flags().StringArray("command", nil, "Command to run inside the container (e.g., --command sh -c 'echo hello')")
	upCmd.Flags().StringArray("env", nil, "Environment variables (e.g., --env KEY=VALUE)")
	upCmd.Flags().StringArray("port", nil, "Port mappings (e.g., --port 8080:80)")
	upCmd.Flags().StringArray("volume", nil, "Volume bindings (e.g., --volume /host:/container)")
}
