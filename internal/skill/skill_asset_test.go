package skill

import (
	"strings"
	"testing"

	specfs "github.com/9Ashwin/spec-cli/embed"
)

func TestOpsxSuperSkillHasFrontDoorGuardrails(t *testing.T) {
	data, err := specfs.SkillsFS.ReadFile("assets/skills/opsx-super/SKILL.md")
	if err != nil {
		t.Fatalf("read English opsx:super skill: %v", err)
	}

	content := string(data)
	required := []string{
		"description: \"Use when",
		"Instruction Priority",
		"STOP",
		"openspec list --json",
		"openspec new change",
		"Do not silently",
		"Do not write to `docs/superpowers/specs/`",
		"Do not write to `docs/superpowers/plans/`",
		"Red Flags",
	}
	for _, want := range required {
		if !strings.Contains(content, want) {
			t.Fatalf("English opsx:super skill missing %q", want)
		}
	}
}

func TestOpsxSuperChineseSkillHasSameGuardrails(t *testing.T) {
	data, err := specfs.SkillsZHFS.ReadFile("assets/skills-zh/opsx-super/SKILL.md")
	if err != nil {
		t.Fatalf("read Chinese opsx:super skill: %v", err)
	}

	content := string(data)
	required := []string{
		"指令优先级",
		"停止",
		"openspec list --json",
		"openspec new change",
		"不要静默",
		"不要写入 `docs/superpowers/specs/`",
		"不要写入 `docs/superpowers/plans/`",
		"危险信号",
	}
	for _, want := range required {
		if !strings.Contains(content, want) {
			t.Fatalf("Chinese opsx:super skill missing %q", want)
		}
	}
}
