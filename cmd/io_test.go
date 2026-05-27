package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/9Ashwin/spec-cli/internal/cmdutil"
	"github.com/9Ashwin/spec-cli/internal/openspec"
)

func TestSetIONormalizesPartialStreams(t *testing.T) {
	out := &bytes.Buffer{}
	restore := SetIO(&cmdutil.IOStreams{Out: out})
	defer restore()

	if defaultIO.In == nil {
		t.Fatal("expected In to be filled")
	}
	if defaultIO.Out != out {
		t.Fatal("expected Out to be preserved")
	}
	if defaultIO.ErrOut == nil {
		t.Fatal("expected ErrOut to be filled")
	}
}

func TestPrintJSONWritesToConfiguredOut(t *testing.T) {
	streams, out, _ := cmdutil.TestIOStreams(t)
	restore := SetIO(streams)
	defer restore()

	printJSON(map[string]string{"status": "ok"})

	if got := out.String(); !strings.Contains(got, `"status": "ok"`) {
		t.Fatalf("expected JSON on configured stdout, got %q", got)
	}
}

func TestCompletionWritesToConfiguredOut(t *testing.T) {
	streams, out, _ := cmdutil.TestIOStreams(t)
	restore := SetIO(streams)
	defer restore()

	if err := completionCmd.RunE(completionCmd, []string{"bash"}); err != nil {
		t.Fatalf("completion command failed: %v", err)
	}

	if got := out.String(); !strings.Contains(got, "spec-cli") {
		t.Fatalf("expected completion script on configured stdout, got %q", got)
	}
}

func TestSetIOConfiguresOpenSpecLogWriter(t *testing.T) {
	streams, _, errOut := cmdutil.TestIOStreams(t)
	restore := SetIO(streams)
	defer restore()

	if _, err := openspec.LogWriter.Write([]byte("installing\n")); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	if got := errOut.String(); got != "installing\n" {
		t.Fatalf("expected openspec logs on configured stderr, got %q", got)
	}
}
