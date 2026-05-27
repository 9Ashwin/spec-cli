package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/9Ashwin/spec-cli/internal/cmdutil"
	"github.com/9Ashwin/spec-cli/internal/openspec"
)

const (
	ScopeProject = "project"
	ScopeGlobal  = "global"

	LangEN = "en"
	LangZH = "zh"
)

// defaultIO is the package-level IOStreams used by all commands.
// Tests may replace it via SetIO() to capture output.
var defaultIO = cmdutil.SystemIO()

// SetIO replaces the package IOStreams. Returns a restore function.
// Use in tests: defer SetIO(testStreams)()
func SetIO(io *cmdutil.IOStreams) func() {
	prev := defaultIO
	prevLogWriter := openspec.LogWriter
	defaultIO = io.Normalize()
	applyCommandIO()
	return func() {
		defaultIO = prev
		openspec.LogWriter = prevLogWriter
		applyCommandIO()
	}
}

func applyCommandIO() {
	if rootCmd == nil || defaultIO == nil {
		return
	}
	rootCmd.SetIn(defaultIO.In)
	rootCmd.SetOut(defaultIO.Out)
	rootCmd.SetErr(defaultIO.ErrOut)
	openspec.LogWriter = defaultIO.ErrOut
}

// resolveProjectPath returns an absolute path from the optional args slice.
func resolveProjectPath(args []string) (string, error) {
	p := "."
	if len(args) > 0 {
		p = args[0]
	}
	return filepath.Abs(p)
}

// printJSON marshals v to indented JSON and writes to defaultIO.Out.
func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}
	fmt.Fprintln(defaultIO.Out, string(data))
}

// printer is a conditional stderr writer that suppresses output in JSON mode.
type printer struct {
	io    *cmdutil.IOStreams
	quiet bool
}

func newPrinter(jsonMode bool) *printer {
	return &printer{io: defaultIO, quiet: jsonMode}
}

func (p *printer) printf(format string, a ...interface{}) {
	if !p.quiet {
		fmt.Fprintf(p.io.ErrOut, format, a...)
	}
}
