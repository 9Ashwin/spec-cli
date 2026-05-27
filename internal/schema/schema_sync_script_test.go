package schema

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckSchemaSyncScriptDetectsDrift(t *testing.T) {
	tmp := t.TempDir()
	sourceRoot := filepath.Join(tmp, "openspec-schemas")
	source := filepath.Join(sourceRoot, "superpowers-bridge")
	target := filepath.Join(tmp, "embed", "assets", "schemas", "superpowers-bridge")

	writeTestFile(t, filepath.Join(source, "schema.yaml"), "name: superpowers-bridge\nversion: 2\n")
	writeTestFile(t, filepath.Join(target, "schema.yaml"), "name: superpowers-bridge\nversion: 1\n")

	cmd := exec.Command("bash", filepath.Join("..", "..", "scripts", "check-schema-sync.sh"), sourceRoot)
	cmd.Env = append(os.Environ(), "SPEC_CLI_SCHEMA_EMBED_DIR="+target)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected drift check to fail")
	}
	if got := string(output); !strings.Contains(got, "schema.yaml") {
		t.Fatalf("expected output to mention schema.yaml, got %q", got)
	}
}

func TestCheckSchemaSyncScriptPassesWhenCanonicalFilesMatch(t *testing.T) {
	tmp := t.TempDir()
	sourceRoot := filepath.Join(tmp, "openspec-schemas")
	source := filepath.Join(sourceRoot, "superpowers-bridge")
	target := filepath.Join(tmp, "embed", "assets", "schemas", "superpowers-bridge")

	writeTestFile(t, filepath.Join(source, "schema.yaml"), "name: superpowers-bridge\nversion: 1\n")
	writeTestFile(t, filepath.Join(target, "schema.yaml"), "name: superpowers-bridge\nversion: 1\n")
	locale := unsupportedLocale()
	traditionalLabel := string([]rune{'\u7e41', '\u9ad4'}) + "中文"
	writeTestFile(t, filepath.Join(source, "README.md"), "["+traditionalLabel+"](./README."+locale+".md)\n")
	writeTestFile(t, filepath.Join(target, "README.md"), "[简体中文](./README.zh.md)\n")
	writeTestFile(t, filepath.Join(source, "README."+locale+".md"), "upstream locale skipped\n")
	writeTestFile(t, filepath.Join(target, "README.zh.md"), "local locale\n")

	cmd := exec.Command("bash", filepath.Join("..", "..", "scripts", "check-schema-sync.sh"), sourceRoot)
	cmd.Env = append(os.Environ(), "SPEC_CLI_SCHEMA_EMBED_DIR="+target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected sync check to pass, err=%v output=%s", err, output)
	}
}

func TestCheckSchemaSyncScriptDoesNotRequireUnsupportedLocaleFiles(t *testing.T) {
	tmp := t.TempDir()
	sourceRoot := filepath.Join(tmp, "openspec-schemas")
	source := filepath.Join(sourceRoot, "superpowers-bridge")
	target := filepath.Join(tmp, "embed", "assets", "schemas", "superpowers-bridge")

	writeTestFile(t, filepath.Join(source, "schema.yaml"), "name: superpowers-bridge\nversion: 1\n")
	writeTestFile(t, filepath.Join(target, "schema.yaml"), "name: superpowers-bridge\nversion: 1\n")
	writeTestFile(t, filepath.Join(source, "templates", "adopters", "CLAUDE.md.fragment."+unsupportedLocale()+".md"), "upstream locale skipped\n")

	cmd := exec.Command("bash", filepath.Join("..", "..", "scripts", "check-schema-sync.sh"), sourceRoot)
	cmd.Env = append(os.Environ(), "SPEC_CLI_SCHEMA_EMBED_DIR="+target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected sync check to pass without unsupported locale targets, err=%v output=%s", err, output)
	}
}

func unsupportedLocale() string {
	return "zh-" + string([]byte{'T', 'W'})
}

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
