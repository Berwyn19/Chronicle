package cmd

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/berwyn/chronicle/internal/store"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <event-id>",
	Short: "Show the metadata for a single event",
	Long: `Show the recorded metadata for a single AI interaction: its prompt,
model, timestamp, commit, and the path to its Git patch.

Use 'chronicle diff <event-id>' to view the patch itself.`,
	Args: cobra.ExactArgs(1),
	RunE: runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	id := args[0]

	s, _, err := openStore()
	if err != nil {
		return err
	}
	defer s.Close()

	e, err := s.Get(id)
	if errors.Is(err, store.ErrNotFound) {
		return fmt.Errorf("no event with id %q (try 'chronicle list')", id)
	}
	if err != nil {
		return err
	}

	commit := e.CommitHash
	if commit == "" {
		commit = "(none)"
	}

	// Label/value pairs aligned with tabwriter (standard library).
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", e.ID)
	fmt.Fprintf(w, "Timestamp:\t%s\n", e.Timestamp)
	fmt.Fprintf(w, "Model:\t%s\n", e.Model)
	fmt.Fprintf(w, "Commit:\t%s\n", commit)
	fmt.Fprintf(w, "Patch:\t%s\n", e.PatchPath)
	fmt.Fprintf(w, "Prompt:\t%s\n", e.Prompt)
	return w.Flush()
}
