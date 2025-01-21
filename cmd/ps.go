package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var all bool

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running containers",
	Long: `Lists running containers with details such as ID, status, and image.
Use -a to show all containers, including stopped ones.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listContainers(all)
	},
}

func listContainers(showAll bool) error {
	cmdArgs := []string{"ps"}
	if showAll {
		cmdArgs = append(cmdArgs, "-a")
	}

	cmd := exec.Command(containerRuntime, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Fetching container list...")
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(psCmd)

	//
	psCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers, including stopped ones")
}
