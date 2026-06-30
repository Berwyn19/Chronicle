package cmd

import (
	"fmt"
	"os"

	"github.com/berwyn/chronicle/internal/paths"
	"github.com/berwyn/chronicle/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a .chronicle directory in the current directory",
	Long: `Initialize Chronicle in the current directory.

Creates a .chronicle directory containing the metadata database and an events
directory for Git patches. Running init again is safe: it leaves an existing
.chronicle directory untouched.`,
	Args: cobra.NoArgs,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	p, err := paths.ForCWD()
	if err != nil {
		return err
	}

	if p.Exists() {
		fmt.Fprintf(cmd.OutOrStdout(), "Chronicle already initialized at %s\n", p.Root)
		return nil
	}

	// Create .chronicle/ and .chronicle/events/. MkdirAll on the events dir
	// creates the parent too.
	if err := os.MkdirAll(p.EventsDir(), 0o755); err != nil {
		return fmt.Errorf("create chronicle directory: %w", err)
	}

	// Opening the store creates and initializes the SQLite database.
	s, err := store.Open(p.DB())
	if err != nil {
		return err
	}
	if err := s.Close(); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Initialized Chronicle in %s\n", p.Root)
	return nil
}
