package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "iptv-cli",
	Short: "A CLI tool for managing IPTV services",
}

// InitializeRootCmd initializes the root command with necessary dependencies
func InitializeRootCmd() *cobra.Command {
	// Add subcommands to rootCmd
	rootCmd.AddCommand(updateVODCmd)
	return rootCmd
}

