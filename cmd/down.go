package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/vvvdwbvvv/rover/pkg/storage"

	"github.com/spf13/cobra"
)

// downCmd 停用容器
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove containers",
	Long:  `Stop and remove all Podman containers or only Rover-managed containers (rover down -l)`,
	Run: func(cmd *cobra.Command, args []string) {
		stopRoverContainers, _ := cmd.Flags().GetBool("last")

		if stopRoverContainers {
			stopRoverManagedContainers()
		} else {
			stopAllContainers()
		}
	},
}

// 停止並刪除所有 Podman 容器
func stopAllContainers() {
	fmt.Println("🛑 Stopping all Podman containers...")
	cmd := exec.Command("podman", "stop", "-a") // 停止所有容器
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Error stopping containers:", err)
		os.Exit(1)
	}

	cmd = exec.Command("podman", "rm", "-a") // 刪除所有容器
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Error removing containers:", err)
		os.Exit(1)
	}
	fmt.Println("✅ All containers have been stopped and removed.")
}

// 停止並刪除 Rover 啟動的容器
func stopRoverManagedContainers() {
	db, err := storage.NewBoltDB("rover.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	containers, err := db.GetContainers()
	if err != nil {
		log.Fatal(err)
	}

	if len(containers) == 0 {
		fmt.Println("🔹 No Rover-managed containers found.")
		return
	}

	for _, c := range containers {
		fmt.Printf("🛑 Stopping container %s...\n", c.Name)
		exec.Command("podman", "stop", c.Name).Run()
		exec.Command("podman", "rm", c.Name).Run()
		db.DeleteContainer(c.Name) // 從 BoltDB 刪除記錄
	}

	fmt.Println("✅ Rover-managed containers have been stopped and removed.")
}

func init() {
	downCmd.Flags().BoolP("last", "l", false, "Stop only Rover-managed containers")
	rootCmd.AddCommand(downCmd)
}
