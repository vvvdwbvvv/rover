package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// logsCmd 取得logs
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve container logs",
	Long:  `Display logs of the specified container, with support for streaming logs.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		container := args[0]
		fmt.Println("Rover: Retrieving container logs...")

		err := getContainerLogs(container)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func getContainerLogs(container string) error {
	cmd := exec.Command("podman", "logs", "-f", container)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
