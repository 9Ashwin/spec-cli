package openspec

import (
	"encoding/json"
	"testing"
)

// mockRunner implements Runner for testing.
type mockRunner struct {
	output    []byte
	err       error
	available bool
	calls     []string
}

func (m *mockRunner) Run(dir, name string, args ...string) ([]byte, error) {
	m.calls = append(m.calls, name)
	return m.output, m.err
}

func (m *mockRunner) Available(name string) bool {
	return m.available
}

func TestListChanges_unavailable(t *testing.T) {
	orig := DefaultRunner
	defer func() { DefaultRunner = orig }()

	DefaultRunner = &mockRunner{available: false}

	changes, err := ListChanges("/tmp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changes != nil {
		t.Errorf("expected nil changes, got %v", changes)
	}
}

func TestListChanges_envelopeFormat(t *testing.T) {
	orig := DefaultRunner
	defer func() { DefaultRunner = orig }()

	data, _ := json.Marshal(struct {
		Changes []ChangeInfo `json:"changes"`
	}{
		Changes: []ChangeInfo{
			{Name: "feature-x", Schema: "superpowers-bridge"},
		},
	})

	DefaultRunner = &mockRunner{available: true, output: data}

	changes, err := ListChanges("/tmp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Name != "feature-x" {
		t.Errorf("expected name 'feature-x', got %q", changes[0].Name)
	}
}

func TestListChanges_flatArrayFormat(t *testing.T) {
	orig := DefaultRunner
	defer func() { DefaultRunner = orig }()

	data, _ := json.Marshal([]ChangeInfo{
		{Name: "fix-y"},
	})

	DefaultRunner = &mockRunner{available: true, output: data}

	changes, err := ListChanges("/tmp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 1 || changes[0].Name != "fix-y" {
		t.Errorf("unexpected changes: %+v", changes)
	}
}

func TestVersion_unavailable(t *testing.T) {
	orig := DefaultRunner
	defer func() { DefaultRunner = orig }()

	DefaultRunner = &mockRunner{available: false}

	_, err := Version()
	if err == nil {
		t.Error("expected error when openspec not available")
	}
}

func TestVersion_available(t *testing.T) {
	orig := DefaultRunner
	defer func() { DefaultRunner = orig }()

	DefaultRunner = &mockRunner{available: true, output: []byte("1.2.3\n")}

	v, err := Version()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "1.2.3" {
		t.Errorf("expected '1.2.3', got %q", v)
	}
}
