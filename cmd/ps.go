package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// psCmd 顯示狀態
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Check running containers",
	Long:  `Listed all running containers under rover`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rover: Check running containers...")

		err := listRunningContainers()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func listRunningContainers() error {
	// `podman ps` 的 `--format` 用於自定義輸出格式
	cmd := exec.Command("podman", "ps", "--format",
		"table {{.ID}}\t{{.Image}}\t{{.Command}}\t{{.CreatedAt}}\t{{.Status}}\t{{.Ports}}\t{{.Names}}")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(psCmd)
}
