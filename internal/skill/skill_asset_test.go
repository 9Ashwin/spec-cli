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
		"<EXTREMELY-IMPORTANT>",
		"This is a GATED workflow",
		"Instruction Priority",
		"User instructions always take precedence",
		"STOP",
		"openspec list --json",
		"openspec new change",
		"Continuous Execution",
		"After creating or selecting a change, enter the Continuous Execution loop below",
		"openspec status --change \"<name>\" --json",
		"openspec instructions <artifact-id> --change \"<name>\" --json",
		"openspec instructions apply --change \"<name>\" --json",
		"next incomplete schema step",
		"Do not delegate Continuous Execution to `openspec-continue-change`",
		"apply action",
		"Invoke the relevant Superpowers skill before acting on that schema step.",
		"Skill output path overrides",
		"Before advancing to the next phase, verify the current phase's EXIT GATE",
		"After creation, enter Continuous Execution",
		"Do not silently",
		"Do not write to `docs/superpowers/specs/`",
		"Do not write to `docs/superpowers/plans/`",
		"Red Flags",
		"These thoughts mean STOP",
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
		"<EXTREMELY-IMPORTANT>",
		"这是一个 GATED 工作流",
		"用户指令始终优先",
		"停止",
		"openspec list --json",
		"openspec new change",
		"连续执行",
		"创建或选择 change 后，进入下方的 Continuous Execution 循环",
		"openspec status --change \"<name>\" --json",
		"openspec instructions <artifact-id> --change \"<name>\" --json",
		"openspec instructions apply --change \"<name>\" --json",
		"下一个未完成的 schema 步骤",
		"不要把连续执行委托给 `openspec-continue-change`",
		"apply action",
		"先调用相关 Superpowers skill，再执行该 schema step。",
		"Skill 输出路径覆盖",
		"推进到下一阶段前，验证当前阶段的 EXIT GATE",
		"创建后进入连续执行",
		"不要静默",
		"不要写入 `docs/superpowers/specs/`",
		"不要写入 `docs/superpowers/plans/`",
		"危险信号",
		"这些想法表示你正在自我合理化",
	}
	for _, want := range required {
		if !strings.Contains(content, want) {
			t.Fatalf("Chinese opsx:super skill missing %q", want)
		}
	}
}
