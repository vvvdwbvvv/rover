package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rover",
	Short: "Rover: Lightweight container orchestrator",
	Long:  `Rover is a lightweight container orchestrator that allows you to deploy and manage containers on a single host.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(logsCmd)
}
