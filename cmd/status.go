package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}
	projectPath, err := filepath.Abs(projectPath)
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
		data, _ := json.MarshalIndent(changes, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	if len(changes) == 0 {
		fmt.Println("No active changes.")
		return nil
	}

	fmt.Fprintf(os.Stderr, "\n  Active Changes (%d):\n\n", len(changes))
	for _, c := range changes {
		schemaLabel := ""
		if c.Schema != "" {
			schemaLabel = fmt.Sprintf(" [schema: %s]", c.Schema)
		}
		fmt.Fprintf(os.Stderr, "  • %s%s\n", c.Name, schemaLabel)
	}
	fmt.Fprintln(os.Stderr)

	return nil
}
