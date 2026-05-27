package cmd

import (
	"fmt"

	"github.com/9Ashwin/spec-cli/internal/build"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spec-cli",
	Short: "Install OpenSpec, Superpowers, and schema bundles",
	Long: `spec-cli — OpenSpec + Superpowers workflow scaffolding tool.

spec-cli detects AI coding platforms and installs:
  - OpenSpec skills (spec lifecycle management)
  - Superpowers skills (brainstorming, TDD, code review)
  - opsx:super entry skill (thin workflow guide)
  - Schema bundles (workflow definitions for openspec/schemas/)

Commands:
  spec-cli init [path]     Initialize workflow scaffolding
  spec-cli status [path]   Show active changes
  spec-cli update [path]   Update packages and schemas
  spec-cli doctor [path]   Diagnose installation health

Examples:
  spec-cli init              # Interactive setup in current directory
  spec-cli init --yes        # Non-interactive, auto-detect platforms
  spec-cli status            # Show active workflow changes
  spec-cli doctor            # Check installation health`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and returns the process exit code.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(defaultIO.ErrOut, "Error:", err)
		return 1
	}
	return 0
}

func init() {
	rootCmd.Version = build.GetInfo()
	rootCmd.SetVersionTemplate("spec-cli version {{.Version}}\n")
	applyCommandIO()
	rootCmd.AddCommand(initCmd, statusCmd, updateCmd, doctorCmd, completionCmd)
}
