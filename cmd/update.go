package cmd

import (
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
	updateCmd.Flags().StringVar(&updateLanguage, "language", LangEN, "Language: en | zh")
	updateCmd.Flags().StringVar(&updateScope, "scope", ScopeProject, "Update scope: project | global")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	projectPath, err := resolveProjectPath(args)
	if err != nil {
		return err
	}

	log := newPrinter(updateJSON)

	baseDir := projectPath
	if updateScope == ScopeGlobal {
		home, err := vfs.UserHomeDir()
		if err != nil {
			return err
		}
		baseDir = home
	}

	type updateResult struct {
		SkillsUpdated  int            `json:"skillsUpdated"`
		SchemasUpdated int            `json:"schemasUpdated"`
		OpsxSuper      map[string]int `json:"opsxSuper"`
	}

	result := updateResult{OpsxSuper: make(map[string]int)}

	// Update skills (overwrite mode)
	detected := platform.DetectPlatforms(projectPath)
	log.printf("\n  Updating skills...\n")
	for _, p := range detected {
		skillsDir := p.SkillsDir
		if updateScope == ScopeGlobal && p.GlobalSkillsDir != "" {
			skillsDir = p.GlobalSkillsDir
		}
		copied, _, err := skill.CopySkills(baseDir, skillsDir, updateLanguage, true)
		if err != nil {
			log.printf("  %s: error — %v\n", p.Name, err)
		} else {
			result.OpsxSuper[p.ID] = copied
			result.SkillsUpdated += copied
			log.printf("  %s: %d updated\n", p.Name, copied)
		}
	}

	// Update schemas
	schemas, err := schema.ListSchemas()
	if err == nil {
		log.printf("\n  Updating schemas...\n")
		for _, s := range schemas {
			installed := schema.GetInstalledVersion(s.Name, projectPath)
			if installed != "" && installed == s.Version {
				log.printf("  %s: up to date (v%s)\n", s.Name, s.Version)
				continue
			}
			if err := schema.InstallSchema(s.Name, projectPath); err != nil {
				log.printf("  %s: failed — %v\n", s.Name, err)
			} else {
				result.SchemasUpdated++
				log.printf("  %s: updated to v%s\n", s.Name, s.Version)
			}
		}
	}

	if updateJSON {
		printJSON(result)
	} else {
		log.printf("\n  Update complete. %d skills, %d schemas updated.\n\n",
			result.SkillsUpdated, result.SchemasUpdated)
	}

	return nil
}
