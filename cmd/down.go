package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// downCmd 關閉所有容器
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop all containers",
	Long:  `Stop all containers started by Rover and clean up resources`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rover: Stop all containers...")

		err := stopAllContainers()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func stopAllContainers() error {
	// 停止所有容器
	cmd := exec.Command("podman", "stop", "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	// 刪除所有容器
	cmd = exec.Command("podman", "rm", "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(downCmd)
}
