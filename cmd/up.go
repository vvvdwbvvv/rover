package cmd

import (
	"fmt"
	"github.com/vvvdwbvvv/rover/internal/config"
	"github.com/vvvdwbvvv/rover/internal/container"
	"os"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up [container-name]",
	Short: "Launch containers from rover-compose.yaml (or specific container)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// parse rover-compose.yaml / toml / json
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error loading config: %v\n", err)
			os.Exit(1)
		}

		startOrder, err := config.GetServiceStartupOrder(cfg.Services)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Dependency error: %v\n", err)
			os.Exit(1)
		}

		// parse CLI flags
		extraEnv, _ := cmd.Flags().GetStringArray("env")
		extraPorts, _ := cmd.Flags().GetStringArray("port")
		extraVolumes, _ := cmd.Flags().GetStringArray("volume")

		var containerName string
		if len(args) > 0 {
			containerName = args[0]
			if _, exists := cfg.Services[containerName]; !exists {
				fmt.Fprintf(os.Stderr, "❌ Container '%s' not found in rover-compose\n", containerName)
				os.Exit(1)
			}
			startOrder = []string{containerName}
		}

		fmt.Println("🚀 Starting containers in order:", startOrder)

		for _, serviceName := range startOrder {
			service := cfg.Services[serviceName]

			envVars := append(service.EnvironmentToArray(), extraEnv...)
			ports := append(service.Ports, extraPorts...)
			volumes := append(service.Volumes, extraVolumes...)

			fmt.Printf("[+] Starting %s (%s)...\n", service.Name, service.Image)
			if err := container.StartContainer(service.Name, service.Image, nil, envVars, ports, volumes); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Failed to start %s: %v\n", service.Name, err)
				os.Exit(1)
			}
		}

		fmt.Println("✅ All containers started successfully.")
	},
}

func (s Service) EnvironmentToArray() []string {
	env := []string{}
	for key, value := range s.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}

func init() {
	upCmd.Flags().StringArray("env", nil, "Extra environment variables (e.g., --env KEY=VALUE)")
	upCmd.Flags().StringArray("port", nil, "Extra port mappings (e.g., --port 8080:80)")
	upCmd.Flags().StringArray("volume", nil, "Extra volume bindings (e.g., --volume /host:/container)")
}
