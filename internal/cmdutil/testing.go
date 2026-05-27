package cmdutil

import (
	"bytes"
	"strings"
	"testing"
)

// TestIOStreams creates buffered streams for command tests.
func TestIOStreams(t testing.TB) (*IOStreams, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	return NewIOStreams(strings.NewReader(""), out, errOut), out, errOut
}
