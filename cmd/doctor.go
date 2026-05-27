package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/9Ashwin/spec-cli/internal/openspec"
	"github.com/9Ashwin/spec-cli/internal/platform"
	"github.com/9Ashwin/spec-cli/internal/schema"
	"github.com/9Ashwin/spec-cli/internal/vfs"
	"github.com/spf13/cobra"
)

var (
	doctorJSON  bool
	doctorScope string
)

var doctorCmd = &cobra.Command{
	Use:   "doctor [path]",
	Short: "Diagnose installation health",
	Long:  "Check OpenSpec CLI, working directories, schema bundles, and skill files.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runDoctor,
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Output structured JSON")
	doctorCmd.Flags().StringVar(&doctorScope, "scope", "auto", "Check scope: auto | project | global")
}

type doctorCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "ok", "warning", "error"
	Detail string `json:"detail,omitempty"`
}

func runDoctor(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}
	projectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	var checks []doctorCheck

	// Check 1: openspec CLI
	version, err := openspec.Version()
	if err != nil {
		checks = append(checks, doctorCheck{
			Name: "OpenSpec CLI", Status: "error",
			Detail: fmt.Sprintf("not found on PATH: %v", err),
		})
	} else {
		checks = append(checks, doctorCheck{
			Name: "OpenSpec CLI", Status: "ok",
			Detail: fmt.Sprintf("version %s", version),
		})
	}

	// Check 2: Schema bundles
	schemas, err := schema.ListSchemas()
	if err != nil {
		checks = append(checks, doctorCheck{
			Name: "Schemas", Status: "error", Detail: err.Error(),
		})
	} else {
		for _, s := range schemas {
			if schema.IsInstalled(s.Name, projectPath) {
				installed := schema.GetInstalledVersion(s.Name, projectPath)
				if installed != s.Version {
					checks = append(checks, doctorCheck{
						Name: fmt.Sprintf("Schema: %s", s.Name), Status: "warning",
						Detail: fmt.Sprintf("installed v%s, available v%s", installed, s.Version),
					})
				} else {
					checks = append(checks, doctorCheck{
						Name: fmt.Sprintf("Schema: %s", s.Name), Status: "ok",
						Detail: fmt.Sprintf("v%s", installed),
					})
				}
			} else {
				checks = append(checks, doctorCheck{
					Name: fmt.Sprintf("Schema: %s", s.Name), Status: "warning",
					Detail: "not installed — run spec-cli init",
				})
			}
		}
	}

	// Check 4: Skill files for detected platforms
	detected := platform.DetectPlatforms(projectPath)
	if len(detected) == 0 {
		checks = append(checks, doctorCheck{
			Name: "Platform Skills", Status: "warning",
			Detail: "no platforms detected",
		})
	}
	for _, p := range detected {
		skillPath := filepath.Join(projectPath, p.SkillsDir, "skills", "opsx-super", "SKILL.md")
		if _, err := vfs.Stat(skillPath); err != nil {
			checks = append(checks, doctorCheck{
				Name: fmt.Sprintf("Skills: %s", p.Name), Status: "warning",
				Detail: fmt.Sprintf("%s/skills/opsx-super/SKILL.md not found", p.SkillsDir),
			})
		} else {
			checks = append(checks, doctorCheck{
				Name: fmt.Sprintf("Skills: %s", p.Name), Status: "ok",
				Detail: fmt.Sprintf("%s/skills/opsx-super/SKILL.md present", p.SkillsDir),
			})
		}
	}

	// Output
	if doctorJSON {
		data, _ := json.MarshalIndent(checks, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Fprintf(os.Stderr, "\n  spec-cli Doctor\n\n")
	okCount, warnCount, errCount := 0, 0, 0
	for _, c := range checks {
		icon := "✓"
		switch c.Status {
		case "warning":
			icon = "!"
			warnCount++
		case "error":
			icon = "✗"
			errCount++
		default:
			okCount++
		}
		detail := ""
		if c.Detail != "" {
			detail = fmt.Sprintf(" — %s", c.Detail)
		}
		fmt.Fprintf(os.Stderr, "  %s %s%s\n", icon, c.Name, detail)
	}
	fmt.Fprintf(os.Stderr, "\n  %d ok, %d warnings, %d errors\n\n", okCount, warnCount, errCount)

	if errCount > 0 {
		return fmt.Errorf("doctor found %d error(s)", errCount)
	}
	return nil
}
