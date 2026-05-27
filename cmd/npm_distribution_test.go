package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestNPMDistributionBundlesReleaseArchives(t *testing.T) {
	packageJSON, err := os.ReadFile("../package.json")
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}
	if !strings.Contains(string(packageJSON), `"dist/*"`) {
		t.Fatal("package.json must include dist/* so npm installs can use packaged release archives")
	}

	workflow, err := os.ReadFile("../.github/workflows/release.yml")
	if err != nil {
		t.Fatalf("read release workflow: %v", err)
	}
	workflowContent := string(workflow)
	for _, want := range []string{
		"mkdir -p dist",
		`gh release download "${TAG}" --pattern 'spec-cli-*' --dir dist`,
	} {
		if !strings.Contains(workflowContent, want) {
			t.Fatalf("release workflow missing %q", want)
		}
	}

	installer, err := os.ReadFile("../scripts/install.js")
	if err != nil {
		t.Fatalf("read npm installer: %v", err)
	}
	installerContent := string(installer)
	for _, want := range []string{
		"packagedArchive",
		"path.join(__dirname, \"..\", \"dist\", archiveName)",
		"fs.copyFileSync(packagedArchive, archivePath)",
	} {
		if !strings.Contains(installerContent, want) {
			t.Fatalf("npm installer missing packaged archive fallback %q", want)
		}
	}
}
