package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recorded events, newest first",
	Long: `List recorded AI interactions for this project, newest first.

Shows the event ID, timestamp, model, and prompt for each recorded event.`,
	Args: cobra.NoArgs,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	s, _, err := openStore()
	if err != nil {
		return err
	}
	defer s.Close()

	events, err := s.List()
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if len(events) == 0 {
		fmt.Fprintln(out, "No events recorded yet.")
		return nil
	}

	// Aligned columns via tabwriter (standard library).
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTIMESTAMP\tMODEL\tPROMPT")
	for _, e := range events {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.ID, e.Timestamp, e.Model, truncate(e.Prompt, 60))
	}
	return w.Flush()
}

// truncate shortens s to at most max runes, appending an ellipsis when cut.
func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}
