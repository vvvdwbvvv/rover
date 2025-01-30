package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"rover/pkg/storage"

	"github.com/spf13/cobra"
)

// psCmd 顯示狀態
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running containers",
	Long:  `Show all running containers (podman ps) or only Rover-managed containers (rover ps -l)`,
	Run: func(cmd *cobra.Command, args []string) {
		listRoverContainers, _ := cmd.Flags().GetBool("last")

		if listRoverContainers {
			listRoverManagedContainers()
		} else {
			listAllPodmanContainers()
		}
	},
}

// 列出所有 Podman 容器（與 `podman ps` 一致）
func listAllPodmanContainers() {
	fmt.Println("📌 Listing all running Podman containers...")
	cmd := exec.Command("podman", "ps", "--format",
		"table {{.ID}}\t{{.Image}}\t{{.Status}}\t{{.Names}}")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Error retrieving Podman containers:", err)
		os.Exit(1)
	}
}

// 列出 Rover 啟動的容器
func listRoverManagedContainers() {
	db, err := storage.NewBoltDB("rover.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	containers, err := db.GetContainers()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🚀 Rover-managed containers:")
	if len(containers) == 0 {
		fmt.Println("🔹 No containers were started by Rover.")
		return
	}

	for _, c := range containers {
		fmt.Printf("🟢 %s (ID: %s) - Status: %s\n", c.Name, c.ID, c.Status)
	}
}

func init() {
	psCmd.Flags().BoolP("last", "l", false, "Show only Rover-managed containers")
	rootCmd.AddCommand(psCmd)
}
