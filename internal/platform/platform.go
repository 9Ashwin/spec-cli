package platform

// Platform represents an AI coding platform supported by OpenSpec.
type Platform struct {
	ID              string   // "claude", "cursor", "roocode", ...
	Name            string   // "Claude Code", "Cursor", "RooCode", ...
	SkillsDir       string   // e.g. ".claude", ".cursor"
	GlobalSkillsDir string   // optional, e.g. ".gemini/antigravity" for antigravity
	DetectionPaths  []string // paths checked for detection; nil means fall back to SkillsDir
	OpenSpecToolID  string   // tool ID passed to openspec init --tools
}

// AllPlatforms lists all 29 supported platforms.
// Keep these IDs aligned with OpenSpec's tool registry.
var AllPlatforms = []Platform{
	{ID: "claude", Name: "Claude Code", SkillsDir: ".claude", OpenSpecToolID: "claude"},
	{ID: "cursor", Name: "Cursor", SkillsDir: ".cursor", OpenSpecToolID: "cursor"},
	{ID: "codex", Name: "Codex", SkillsDir: ".codex", OpenSpecToolID: "codex"},
	{ID: "opencode", Name: "OpenCode", SkillsDir: ".opencode", OpenSpecToolID: "opencode"},
	{ID: "windsurf", Name: "Windsurf", SkillsDir: ".windsurf", OpenSpecToolID: "windsurf"},
	{ID: "cline", Name: "Cline", SkillsDir: ".cline", OpenSpecToolID: "cline"},
	{ID: "roocode", Name: "RooCode", SkillsDir: ".roo", OpenSpecToolID: "roocode"},
	{ID: "continue", Name: "Continue", SkillsDir: ".continue", OpenSpecToolID: "continue"},
	{
		ID: "github-copilot", Name: "GitHub Copilot", SkillsDir: ".github",
		DetectionPaths: []string{
			".github/copilot-instructions.md", ".github/instructions",
			".github/prompts", ".github/skills",
		},
		OpenSpecToolID: "github-copilot",
	},
	{ID: "gemini", Name: "Gemini CLI", SkillsDir: ".gemini", OpenSpecToolID: "gemini"},
	{ID: "amazon-q", Name: "Amazon Q Developer", SkillsDir: ".amazonq", OpenSpecToolID: "amazon-q"},
	{ID: "qwen", Name: "Qwen Code", SkillsDir: ".qwen", OpenSpecToolID: "qwen"},
	{ID: "kilocode", Name: "Kilo Code", SkillsDir: ".kilocode", OpenSpecToolID: "kilocode"},
	{ID: "auggie", Name: "Auggie (Augment CLI)", SkillsDir: ".augment", OpenSpecToolID: "auggie"},
	{ID: "kiro", Name: "Kiro", SkillsDir: ".kiro", OpenSpecToolID: "kiro"},
	{ID: "lingma", Name: "Lingma", SkillsDir: ".lingma", OpenSpecToolID: "lingma"},
	{ID: "junie", Name: "Junie", SkillsDir: ".junie", OpenSpecToolID: "junie"},
	{ID: "codebuddy", Name: "CodeBuddy Code", SkillsDir: ".codebuddy", OpenSpecToolID: "codebuddy"},
	{ID: "costrict", Name: "CoStrict", SkillsDir: ".cospec", OpenSpecToolID: "costrict"},
	{ID: "crush", Name: "Crush", SkillsDir: ".crush", OpenSpecToolID: "crush"},
	{ID: "factory", Name: "Factory Droid", SkillsDir: ".factory", OpenSpecToolID: "factory"},
	{ID: "iflow", Name: "iFlow", SkillsDir: ".iflow", OpenSpecToolID: "iflow"},
	{ID: "pi", Name: "Pi", SkillsDir: ".pi", OpenSpecToolID: "pi"},
	{ID: "qoder", Name: "Qoder", SkillsDir: ".qoder", OpenSpecToolID: "qoder"},
	{
		ID: "antigravity", Name: "Antigravity",
		SkillsDir: ".agents", GlobalSkillsDir: ".gemini/antigravity",
		OpenSpecToolID: "antigravity",
	},
	{ID: "bob", Name: "Bob Shell", SkillsDir: ".bob", OpenSpecToolID: "bob"},
	{ID: "forgecode", Name: "ForgeCode", SkillsDir: ".forge", OpenSpecToolID: "forgecode"},
	{ID: "trae", Name: "Trae", SkillsDir: ".trae", OpenSpecToolID: "trae"},
}

// ByID returns the platform with the given ID, or nil if not found.
func ByID(id string) *Platform {
	for i := range AllPlatforms {
		if AllPlatforms[i].ID == id {
			return &AllPlatforms[i]
		}
	}
	return nil
}
