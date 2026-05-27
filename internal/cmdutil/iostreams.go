// Package cmdutil provides shared utilities for CLI commands.
package cmdutil

import (
	"io"
	"os"

	"golang.org/x/term"
)

// IOStreams provides the standard input/output/error streams.
// Commands should use these instead of os.Stdin/Stdout/Stderr
// to enable testing and output capture.
type IOStreams struct {
	In         io.Reader
	Out        io.Writer
	ErrOut     io.Writer
	IsTerminal bool
}

// NewIOStreams builds an IOStreams from arbitrary readers/writers.
// IsTerminal is derived from in's underlying *os.File, if any; non-file
// readers (bytes.Buffer, strings.Reader, ...) yield IsTerminal=false.
func NewIOStreams(in io.Reader, out, errOut io.Writer) *IOStreams {
	isTerminal := false
	if f, ok := in.(*os.File); ok {
		isTerminal = term.IsTerminal(int(f.Fd()))
	}
	return &IOStreams{In: in, Out: out, ErrOut: errOut, IsTerminal: isTerminal}
}

// SystemIO creates an IOStreams wired to the process's standard file descriptors.
func SystemIO() *IOStreams {
	return NewIOStreams(os.Stdin, os.Stdout, os.Stderr)
}

// Normalize returns a fresh IOStreams with nil streams filled from SystemIO.
// It lets tests provide partial streams, such as &IOStreams{Out: buf}, without
// leaking nil readers or writers into command code.
func (s *IOStreams) Normalize() *IOStreams {
	if s == nil {
		return SystemIO()
	}

	out := *s
	if out.In == nil || out.Out == nil || out.ErrOut == nil {
		sys := SystemIO()
		if out.In == nil {
			out.In = sys.In
		}
		if out.Out == nil {
			out.Out = sys.Out
		}
		if out.ErrOut == nil {
			out.ErrOut = sys.ErrOut
		}
	}
	return &out
}
