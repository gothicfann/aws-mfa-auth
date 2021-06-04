package cmd

import (
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Placeholder for all kind of configuration subcommands",
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
