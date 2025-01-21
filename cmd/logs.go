package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var follow bool
var tailLines int

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [container]",
	Short: "Fetch logs for a specific container",
	Long: `Retrieves logs for the given container. Supports options for 
following real-time logs (-f) and displaying the last N lines (--tail).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		containerName := args[0]
		return fetchContainerLogs(containerName, follow, tailLines)
	},
}

func fetchContainerLogs(containerName string, follow bool, tail int) error {
	cmdArgs := []string{"logs"}
	if follow {
		cmdArgs = append(cmdArgs, "-f")
	}
	if tail > 0 {
		cmdArgs = append(cmdArgs, "--tail", fmt.Sprintf("%d", tail))
	}
	cmdArgs = append(cmdArgs, containerName)

	if containerRuntime == "runc" {
		logPath := fmt.Sprintf("/var/lib/runc/containers/%s/log.json", containerName)
		fmt.Printf("Using runc: reading logs from %s\n", logPath)

		catArgs := []string{"cat", logPath}
		if follow {
			catArgs = []string{"tail", "-f", logPath}
		} else if tail > 0 {
			catArgs = []string{"tail", fmt.Sprintf("-%d", tail), logPath}
		}

		cmd := exec.Command(catArgs[0], catArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	cmd := exec.Command(containerRuntime, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Fetching logs for container: %s\n", containerName)
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().IntVar(&tailLines, "tail", 0, "Number of recent log lines to display")
}
