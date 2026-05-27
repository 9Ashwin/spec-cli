package build

import "testing"

func TestGetInfo_nonEmpty(t *testing.T) {
	info := GetInfo()
	if info == "" {
		t.Error("GetInfo() returned empty string")
	}
}

func TestShortCommitHash_length(t *testing.T) {
	h := ShortCommitHash()
	if len(h) > 7 {
		t.Errorf("ShortCommitHash() returned %d chars, expected <= 7", len(h))
	}
}

func TestEnsureBuildInfo_setsVersion(t *testing.T) {
	EnsureBuildInfo()
	if Version == "" {
		t.Error("Version is empty after EnsureBuildInfo()")
	}
}
