// Package cmd wires up Chronicle's Cobra commands. It contains no business
// logic of its own; each command delegates to the internal packages.
package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd is the base "chronicle" command.
var rootCmd = &cobra.Command{
	Use:   "chronicle",
	Short: "Record AI-assisted development history",
	Long: `Chronicle records the intent behind AI-assisted development.

Git remains the source of truth for your source code. Chronicle records the
prompt, model, and metadata for each AI interaction and stores a Git patch of
the resulting changes under .chronicle/.`,
	SilenceUsage: true,
}

// Execute runs the root command. It is the single entrypoint called by main.
func Execute() error {
	return rootCmd.Execute()
}
