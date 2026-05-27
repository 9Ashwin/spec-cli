package cmd

import "github.com/spf13/cobra"

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for spec-cli.

Examples:
  # Bash (add to ~/.bashrc):
  eval "$(spec-cli completion bash)"

  # Zsh (add to ~/.zshrc):
  eval "$(spec-cli completion zsh)"

  # Fish:
  spec-cli completion fish | source

  # PowerShell:
  spec-cli completion powershell | Out-String | Invoke-Expression`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(defaultIO.Out)
		case "zsh":
			return rootCmd.GenZshCompletion(defaultIO.Out)
		case "fish":
			return rootCmd.GenFishCompletion(defaultIO.Out, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(defaultIO.Out)
		}
		return nil
	},
}
