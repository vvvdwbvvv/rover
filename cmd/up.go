package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// upCmd 啟動容器
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start containers",
	Long:  `Run rootless containers.`,

	Run: func(cmd *cobra.Command, args []string) {
		image := args[0]
		fmt.Println("Rover: Starting container...")

		err := runPodmanContainer(image)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func runPodmanContainer(image string) error {
	// 執行podman run
	cmd := exec.Command("podman", "run", "-d", "--name", image, image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(upCmd)
}
