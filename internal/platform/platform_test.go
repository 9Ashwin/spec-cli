package platform

import "testing"

func TestByID_found(t *testing.T) {
	p := ByID("claude")
	if p == nil {
		t.Fatal("expected to find platform 'claude'")
	}
	if p.Name != "Claude Code" {
		t.Errorf("expected Name='Claude Code', got %q", p.Name)
	}
	if p.SkillsDir != ".claude" {
		t.Errorf("expected SkillsDir='.claude', got %q", p.SkillsDir)
	}
}

func TestByID_notFound(t *testing.T) {
	p := ByID("nonexistent-platform")
	if p != nil {
		t.Errorf("expected nil for unknown platform, got %+v", p)
	}
}

func TestAllPlatforms_noDuplicateIDs(t *testing.T) {
	seen := make(map[string]bool)
	for _, p := range AllPlatforms {
		if seen[p.ID] {
			t.Errorf("duplicate platform ID: %s", p.ID)
		}
		seen[p.ID] = true
	}
}

func TestAllPlatforms_requiredFields(t *testing.T) {
	for _, p := range AllPlatforms {
		if p.ID == "" {
			t.Error("platform with empty ID")
		}
		if p.Name == "" {
			t.Errorf("platform %q has empty Name", p.ID)
		}
		if p.SkillsDir == "" {
			t.Errorf("platform %q has empty SkillsDir", p.ID)
		}
		if p.OpenSpecToolID == "" {
			t.Errorf("platform %q has empty OpenSpecToolID", p.ID)
		}
	}
}
