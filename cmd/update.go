package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/9Ashwin/spec-cli/internal/platform"
	"github.com/9Ashwin/spec-cli/internal/schema"
	"github.com/9Ashwin/spec-cli/internal/skill"
	"github.com/9Ashwin/spec-cli/internal/vfs"
	"github.com/spf13/cobra"
)

var (
	updateJSON     bool
	updateLanguage string
	updateScope    string
)

var updateCmd = &cobra.Command{
	Use:   "update [path]",
	Short: "Update packages and schemas",
	Long:  "Re-copy skill files and update schema bundles when versions differ.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().BoolVar(&updateJSON, "json", false, "Output structured JSON")
	updateCmd.Flags().StringVar(&updateLanguage, "language", "en", "Language: en | zh")
	updateCmd.Flags().StringVar(&updateScope, "scope", "project", "Update scope: project | global")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}
	projectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	log := func(format string, a ...interface{}) {
		if !updateJSON {
			fmt.Fprintf(os.Stderr, format, a...)
		}
	}

	baseDir := projectPath
	if updateScope == "global" {
		home, err := vfs.UserHomeDir()
		if err != nil {
			return err
		}
		baseDir = home
	}

	type updateResult struct {
		SkillsUpdated  int            `json:"skillsUpdated"`
		SchemasUpdated int            `json:"schemasUpdated"`
		Comet          map[string]int `json:"comet"`
	}

	result := updateResult{Comet: make(map[string]int)}

	// Update skills (overwrite mode)
	detected := platform.DetectPlatforms(projectPath)
	log("\n  Updating skills...\n")
	for _, p := range detected {
		skillsDir := p.SkillsDir
		if updateScope == "global" && p.GlobalSkillsDir != "" {
			skillsDir = p.GlobalSkillsDir
		}
		copied, _, err := skill.CopySkills(baseDir, skillsDir, updateLanguage, true)
		if err != nil {
			log("  %s: error — %v\n", p.Name, err)
		} else {
			result.Comet[p.ID] = copied
			result.SkillsUpdated += copied
			log("  %s: %d updated\n", p.Name, copied)
		}
	}

	// Update schemas
	schemas, err := schema.ListSchemas()
	if err == nil {
		log("\n  Updating schemas...\n")
		for _, s := range schemas {
			installed := schema.GetInstalledVersion(s.Name, projectPath)
			if installed != "" && installed == s.Version {
				log("  %s: up to date (v%s)\n", s.Name, s.Version)
				continue
			}
			if err := schema.InstallSchema(s.Name, projectPath); err != nil {
				log("  %s: failed — %v\n", s.Name, err)
			} else {
				result.SchemasUpdated++
				log("  %s: updated to v%s\n", s.Name, s.Version)
			}
		}
	}

	if updateJSON {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		log("\n  Update complete. %d skills, %d schemas updated.\n\n",
			result.SkillsUpdated, result.SchemasUpdated)
	}

	return nil
}
