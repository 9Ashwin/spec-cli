package cmdutil

import (
	"bytes"
	"testing"
)

func TestIOStreamsNormalizeFillsNilFields(t *testing.T) {
	out := &bytes.Buffer{}

	streams := (&IOStreams{Out: out}).Normalize()

	if streams.In == nil {
		t.Fatal("expected In to be filled")
	}
	if streams.Out != out {
		t.Fatal("expected existing Out to be preserved")
	}
	if streams.ErrOut == nil {
		t.Fatal("expected ErrOut to be filled")
	}
}

func TestIOStreamsNormalizeNilUsesSystemIO(t *testing.T) {
	streams := (*IOStreams)(nil).Normalize()

	if streams.In == nil || streams.Out == nil || streams.ErrOut == nil {
		t.Fatal("expected nil IOStreams to normalize to complete system streams")
	}
}

func TestTestIOStreamsBuildsBufferedStreams(t *testing.T) {
	streams, out, errOut := TestIOStreams(t)

	if streams.In == nil {
		t.Fatal("expected test input stream")
	}
	if streams.Out != out {
		t.Fatal("expected stdout buffer to be wired")
	}
	if streams.ErrOut != errOut {
		t.Fatal("expected stderr buffer to be wired")
	}
}
