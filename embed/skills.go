package embed

import "embed"

//go:embed all:assets/skills
var SkillsFS embed.FS

//go:embed all:assets/skills-zh
var SkillsZHFS embed.FS
