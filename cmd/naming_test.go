package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestPublicJSONUsesOpsxSuperName(t *testing.T) {
	data, err := json.Marshal(initResult{
		OpsxSuper: map[string]int{"claude": 1},
	})
	if err != nil {
		t.Fatalf("marshal initResult: %v", err)
	}

	got := string(data)
	if !strings.Contains(got, `"opsxSuper"`) {
		t.Fatalf("expected JSON to include opsxSuper, got %s", got)
	}
	forbidden := "co" + "met"
	if strings.Contains(strings.ToLower(got), forbidden) {
		t.Fatalf("JSON contains old workflow name: %s", got)
	}
}

func TestRepositoryDoesNotUseOldWorkflowName(t *testing.T) {
	root := ".."
	forbidden := "co" + "met"

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules":
				return filepath.SkipDir
			}
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !utf8.Valid(data) {
			return nil
		}
		if strings.Contains(strings.ToLower(string(data)), forbidden) {
			t.Errorf("%s contains old workflow name", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan repository names: %v", err)
	}
}

func TestRepositoryDoesNotUseUnsupportedChineseLocale(t *testing.T) {
	root := ".."
	locale := "zh-" + string([]byte{'T', 'W'})
	forbidden := []string{
		locale,
		"README." + locale,
		"fragment." + locale,
		string([]byte{'T', 'r', 'a', 'd', 'i', 't', 'i', 'o', 'n', 'a', 'l'}) + " Chinese",
		string([]rune{'\u7e41', '\u9ad4'}),
		string([]rune{'\u7e41', '\u4f53'}),
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "node_modules":
				return filepath.SkipDir
			}
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !utf8.Valid(data) {
			return nil
		}

		content := string(data)
		for _, value := range forbidden {
			if strings.Contains(content, value) {
				t.Errorf("%s contains unsupported Chinese locale marker %q", path, value)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan repository locales: %v", err)
	}
}
