package schema

import "testing"

func TestListSchemas(t *testing.T) {
	schemas, err := ListSchemas()
	if err != nil {
		t.Fatalf("ListSchemas() error: %v", err)
	}
	if len(schemas) == 0 {
		t.Fatal("expected at least one schema bundle")
	}

	found := false
	for _, s := range schemas {
		if s.Name == "superpowers-bridge" {
			found = true
			if s.Version == "" || s.Version == "unknown" {
				t.Errorf("superpowers-bridge has invalid version: %q", s.Version)
			}
		}
	}
	if !found {
		t.Error("expected superpowers-bridge schema to be listed")
	}
}

func TestRelEmbedPath(t *testing.T) {
	tests := []struct {
		base, target string
		want         string
		wantErr      bool
	}{
		{"assets/schemas/foo", "assets/schemas/foo/bar.yaml", "bar.yaml", false},
		{"assets/schemas/foo", "assets/schemas/foo/sub/file.md", "sub/file.md", false},
		{"assets/schemas/foo", "assets/schemas/foo", ".", false},
		{"assets/schemas/foo", "assets/schemas/other/bar", "", true},
	}

	for _, tt := range tests {
		got, err := relEmbedPath(tt.base, tt.target)
		if tt.wantErr {
			if err == nil {
				t.Errorf("relEmbedPath(%q, %q) expected error", tt.base, tt.target)
			}
			continue
		}
		if err != nil {
			t.Errorf("relEmbedPath(%q, %q) error: %v", tt.base, tt.target, err)
			continue
		}
		if got != tt.want {
			t.Errorf("relEmbedPath(%q, %q) = %q, want %q", tt.base, tt.target, got, tt.want)
		}
	}
}
