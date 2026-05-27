package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/9Ashwin/spec-cli/internal/openspec"
	"github.com/9Ashwin/spec-cli/internal/platform"
	"github.com/9Ashwin/spec-cli/internal/schema"
	"github.com/9Ashwin/spec-cli/internal/skill"
	"github.com/9Ashwin/spec-cli/internal/vfs"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type initOptions struct {
	yes          bool
	skipExisting bool
	overwrite    bool
	jsonOutput   bool
	scope        string
}

var initOpts initOptions

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initialize workflow scaffolding",
	Long: `Initialize OpenSpec, Superpowers, opsx:super entry skill, and schema bundles.

Detects AI coding platforms and interactively installs all components.
Use --yes for non-interactive mode.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&initOpts.yes, "yes", false, "Non-interactive mode")
	initCmd.Flags().BoolVar(&initOpts.skipExisting, "skip-existing", false, "Skip already installed components")
	initCmd.Flags().BoolVar(&initOpts.overwrite, "overwrite", false, "Overwrite all existing components")
	initCmd.Flags().BoolVar(&initOpts.jsonOutput, "json", false, "Output structured JSON")
	initCmd.Flags().StringVar(&initOpts.scope, "scope", "", "Install scope: project | global")
}

type initResult struct {
	ProjectPath       string         `json:"projectPath"`
	Scope             string         `json:"scope"`
	Language          string         `json:"language"`
	SelectedPlatforms []string       `json:"selectedPlatforms"`
	OpenSpec          string         `json:"openspec"`
	Superpowers       string         `json:"superpowers"`
	OpsxSuper         map[string]int `json:"opsxSuper"`
	SchemasInstalled  int            `json:"schemasInstalled"`
}

func runInit(cmd *cobra.Command, args []string) error {
	projectPath, err := resolveProjectPath(args)
	if err != nil {
		return err
	}

	log := newPrinter(initOpts.jsonOutput)

	log.printf("\n")
	log.printf("   ███████╗██████╗ ███████╗ ██████╗\n")
	log.printf("   ██╔════╝██╔══██╗██╔════╝██╔════╝\n")
	log.printf("   ███████╗██████╔╝█████╗  ██║     \n")
	log.printf("   ╚════██║██╔═══╝ ██╔══╝  ██║     \n")
	log.printf("   ███████║██║     ███████╗╚██████╗\n")
	log.printf("   ╚══════╝╚═╝     ╚══════╝ ╚═════╝\n")
	log.printf("      OpenSpec + Superpowers Workflow\n\n")

	// Step 1: Detect platforms
	detected := platform.DetectPlatforms(projectPath)
	if len(detected) > 0 {
		names := make([]string, len(detected))
		for i, p := range detected {
			names[i] = p.Name
		}
		log.printf("  Detected platforms: %s\n", strings.Join(names, ", "))
	}

	// Step 2: Select scope
	scope := initOpts.scope
	if scope == "" {
		if initOpts.yes {
			scope = ScopeProject
		} else {
			scope = selectScope()
		}
	}
	log.printf("  Scope: %s\n", scope)

	// Step 3: Select language
	language := LangEN
	if !initOpts.yes {
		language = selectLanguage()
	}
	log.printf("  Language: %s\n", languageName(language))

	// Step 4: Select platforms
	selected := selectPlatforms(detected)
	if len(selected) == 0 {
		log.printf("\n  No platforms selected. Exiting.\n")
		if initOpts.jsonOutput {
			printJSON(initResult{ProjectPath: projectPath, Scope: scope, Language: language})
		}
		return nil
	}
	log.printf("  Selected: %s\n", strings.Join(platformNames(selected), ", "))

	// Step 5: Determine base directory
	baseDir := projectPath
	if scope == ScopeGlobal {
		home, err := vfs.UserHomeDir()
		if err != nil {
			return err
		}
		baseDir = home
	}

	// Step 6: Install OpenSpec
	var openSpecStatus string
	toolIDs := make([]string, len(selected))
	for i, p := range selected {
		toolIDs[i] = p.OpenSpecToolID
	}
	log.printf("\n  Installing OpenSpec for: %s\n", strings.Join(toolIDs, ", "))
	if err := openspec.InitOpenSpec(projectPath, toolIDs, scope); err != nil {
		log.printf("  OpenSpec: failed — %v\n", err)
		openSpecStatus = "failed"
	} else {
		log.printf("  OpenSpec: installed\n")
		openSpecStatus = "installed"
	}

	// Step 7: Detect Superpowers
	superpowersStatus := "skipped"
	log.printf("\n  Superpowers: checking...\n")
	if checkSuperpowers() {
		log.printf("  Superpowers: detected (plugin-installed)\n")
		superpowersStatus = "detected"
	} else {
		log.printf("  Superpowers: not detected. Install with: claude plugin install superpowers@claude-plugins-official\n")
	}

	// Step 8: Install opsx:super entry skill
	opsxSuperResults := make(map[string]int)
	for _, p := range selected {
		skillsDir := p.SkillsDir
		if scope == ScopeGlobal && p.GlobalSkillsDir != "" {
			skillsDir = p.GlobalSkillsDir
		}
		copied, _, err := skill.CopySkills(baseDir, skillsDir, language, initOpts.overwrite)
		if err != nil {
			log.printf("  opsx:super -> %s: error — %v\n", p.Name, err)
		} else {
			log.printf("  opsx:super -> %s: %d copied\n", p.Name, copied)
			opsxSuperResults[p.ID] = copied
		}
	}

	// Step 9: Create working directories
	specsDir := filepath.Join(projectPath, "docs", "superpowers", "specs")
	plansDir := filepath.Join(projectPath, "docs", "superpowers", "plans")
	if scope == ScopeProject {
		_ = vfs.MkdirAll(specsDir, 0o755)
		_ = vfs.MkdirAll(plansDir, 0o755)
		log.printf("\n  Working directories: docs/superpowers/specs/, docs/superpowers/plans/\n")
	}

	// Step 10: Install schemas
	schemasInstalled := 0
	schemas, err := schema.ListSchemas()
	if err == nil && len(schemas) > 0 {
		schemaNames := make([]string, len(schemas))
		for i, s := range schemas {
			schemaNames[i] = s.Name
		}

		selectedSchemas := schemaNames
		if !initOpts.yes && len(schemas) > 1 {
			selectedSchemas = selectSchemas(schemas)
		}

		for _, s := range schemas {
			if !contains(selectedSchemas, s.Name) {
				continue
			}
			if err := schema.InstallSchema(s.Name, projectPath); err != nil {
				log.printf("  Schema %s: failed — %v\n", s.Name, err)
			} else {
				schemasInstalled++
				log.printf("  Schema: %s installed -> openspec/schemas/%s/\n", s.Name, s.Name)

				if added, _ := schema.AppendClaudeMdFragment(s.Name, projectPath, language); added {
					log.printf("  CLAUDE.md: appended %s workflow fragment\n", s.Name)
				}
			}
		}
	}

	// Summary
	log.printf("\n  spec-cli setup complete! (scope: %s)\n\n", scope)
	log.printf("  Get started:\n")
	log.printf("    /opsx:super \"your idea\"\n")
	log.printf("    openspec new change \"your-idea\" --schema superpowers-bridge --description \"your idea\"\n")
	log.printf("\n  Quick commands (via AI agent skills):\n")
	log.printf("    /opsx:super \"your idea\" — OpenSpec + Superpowers workflow\n\n")

	if initOpts.jsonOutput {
		platformIDs := make([]string, len(selected))
		for i, p := range selected {
			platformIDs[i] = p.ID
		}
		printJSON(initResult{
			ProjectPath:       projectPath,
			Scope:             scope,
			Language:          language,
			SelectedPlatforms: platformIDs,
			OpenSpec:          openSpecStatus,
			Superpowers:       superpowersStatus,
			OpsxSuper:         opsxSuperResults,
			SchemasInstalled:  schemasInstalled,
		})
	}

	return nil
}

// --- Interactive helpers ---

func selectScope() string {
	var scope string
	field := huh.NewSelect[string]().
		Title("Install scope:").
		Options(
			huh.NewOption("Project (current directory)", "project"),
			huh.NewOption("Global (home directory)", "global"),
		).
		Value(&scope)
	_ = huh.NewForm(huh.NewGroup(field)).
		WithInput(defaultIO.In).
		WithOutput(defaultIO.ErrOut).
		Run()
	return scope
}

func selectLanguage() string {
	var lang string
	field := huh.NewSelect[string]().
		Title("Language for opsx:super skills:").
		Options(
			huh.NewOption("English", "en"),
			huh.NewOption("简体中文", "zh"),
		).
		Value(&lang)
	_ = huh.NewForm(huh.NewGroup(field)).
		WithInput(defaultIO.In).
		WithOutput(defaultIO.ErrOut).
		Run()
	return lang
}

func languageName(lang string) string {
	if lang == "zh" {
		return "简体中文"
	}
	return "English"
}

func selectPlatforms(detected []platform.Platform) []platform.Platform {
	if initOpts.yes {
		if len(detected) > 0 {
			return detected
		}
		return platform.AllPlatforms
	}

	detectedSet := make(map[string]bool)
	for _, p := range detected {
		detectedSet[p.ID] = true
	}

	// Pre-select detected platforms. If none detected, pre-select all.
	selected := make([]string, 0, len(platform.AllPlatforms))
	options := make([]huh.Option[string], 0, len(platform.AllPlatforms))
	for _, p := range platform.AllPlatforms {
		label := p.Name
		if detectedSet[p.ID] {
			label += " (detected)"
		}
		options = append(options, huh.NewOption(label, p.ID))
		if detectedSet[p.ID] {
			selected = append(selected, p.ID)
		}
	}

	field := huh.NewMultiSelect[string]().
		Title("Select AI coding platforms").
		Description("x: toggle • a: all • ↑↓: move • /: search • enter: confirm").
		Options(options...).
		Value(&selected).
		Filterable(true)
	_ = huh.NewForm(huh.NewGroup(field)).
		WithInput(defaultIO.In).
		WithOutput(defaultIO.ErrOut).
		Run()

	if len(selected) == 0 {
		return nil
	}

	var result []platform.Platform
	for _, id := range selected {
		if p := platform.ByID(id); p != nil {
			result = append(result, *p)
		}
	}
	return result
}

func selectSchemas(schemas []schema.Info) []string {
	if initOpts.yes || len(schemas) <= 1 {
		names := make([]string, len(schemas))
		for i, s := range schemas {
			names[i] = s.Name
		}
		return names
	}

	// Build multi-select options with all schemas pre-selected.
	options := make([]huh.Option[string], len(schemas))
	selected := make([]string, len(schemas))
	for i, s := range schemas {
		label := fmt.Sprintf("%s (v%s)", s.Name, s.Version)
		options[i] = huh.NewOption(label, s.Name)
		selected[i] = s.Name
	}

	field := huh.NewMultiSelect[string]().
		Title("Select schema bundles:").
		Options(options...).
		Value(&selected)
	_ = huh.NewForm(huh.NewGroup(field)).
		WithInput(defaultIO.In).
		WithOutput(defaultIO.ErrOut).
		Run()

	return selected
}

// --- Superpowers detection ---

var superpowersSkillNames = []string{
	"brainstorming",
	"using-superpowers",
	"writing-plans",
	"test-driven-development",
	"subagent-driven-development",
}

// checkSuperpowers checks if Superpowers is installed via Claude Code plugins.
func checkSuperpowers() bool {
	home, err := vfs.UserHomeDir()
	if err != nil {
		return false
	}

	claudeDir := os.Getenv("CLAUDE_CONFIG_DIR")
	if claudeDir == "" {
		claudeDir = filepath.Join(home, ".claude")
	}

	pluginsCacheDir := filepath.Join(claudeDir, "plugins", "cache")

	marketplaceEntries, err := vfs.ReadDir(pluginsCacheDir)
	if err != nil {
		return false
	}

	for _, marketplace := range marketplaceEntries {
		if !marketplace.IsDir() {
			continue
		}

		superpowersDir := filepath.Join(pluginsCacheDir, marketplace.Name(), "superpowers")
		if _, err := vfs.Stat(superpowersDir); err != nil {
			continue
		}

		versionEntries, err := vfs.ReadDir(superpowersDir)
		if err != nil {
			continue
		}

		for _, version := range versionEntries {
			if !version.IsDir() {
				continue
			}

			skillsDir := filepath.Join(superpowersDir, version.Name(), "skills")
			skillEntries, err := vfs.ReadDir(skillsDir)
			if err != nil {
				continue
			}

			for _, entry := range skillEntries {
				for _, name := range superpowersSkillNames {
					if entry.Name() == name {
						return true
					}
				}
			}
		}
	}

	return false
}

// --- Utilities ---

func platformNames(platforms []platform.Platform) []string {
	names := make([]string, len(platforms))
	for i, p := range platforms {
		names[i] = p.Name
	}
	return names
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
