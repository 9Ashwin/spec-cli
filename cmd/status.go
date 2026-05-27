package cmd

import (
	"fmt"

	"github.com/9Ashwin/spec-cli/internal/openspec"
	"github.com/spf13/cobra"
)

var statusJSON bool

var statusCmd = &cobra.Command{
	Use:   "status [path]",
	Short: "Show active changes",
	Long:  "Display active workflow changes from openspec list --json.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStatus,
}

func init() {
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "Output structured JSON")
}

func runStatus(cmd *cobra.Command, args []string) error {
	projectPath, err := resolveProjectPath(args)
	if err != nil {
		return err
	}

	changes, err := openspec.ListChanges(projectPath)
	if err != nil {
		return err
	}

	if statusJSON {
		if changes == nil {
			changes = []openspec.ChangeInfo{}
		}
		printJSON(changes)
		return nil
	}

	if len(changes) == 0 {
		fmt.Fprintln(defaultIO.Out, "No active changes.")
		return nil
	}

	fmt.Fprintf(defaultIO.ErrOut, "\n  Active Changes (%d):\n\n", len(changes))
	for _, c := range changes {
		schemaLabel := ""
		if c.Schema != "" {
			schemaLabel = fmt.Sprintf(" [schema: %s]", c.Schema)
		}
		fmt.Fprintf(defaultIO.ErrOut, "  • %s%s\n", c.Name, schemaLabel)
	}
	fmt.Fprintln(defaultIO.ErrOut)

	return nil
}
