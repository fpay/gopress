package cmd

import (
	"github.com/spf13/cobra"
)

var moduleRootPath string

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:     "make",
	Short:   "Generate module files",
	Long:    `Use subcommands to generate controllers, services, middlewares, models and views files quickly.`,
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"mk"},
}

func init() {
	RootCmd.AddCommand(makeCmd)

	makeCmd.PersistentFlags().StringVar(&moduleRootPath, "dir", "", "Where the module files located")
}
