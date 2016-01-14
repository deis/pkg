package log

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/arschles/assert"
)

var (
	world = "world"
)

func getWriters() (io.Writer, *bytes.Buffer, io.Writer, *bytes.Buffer) {
	var out, err bytes.Buffer
	stdout, stderr := io.MultiWriter(os.Stdout, &out), io.MultiWriter(os.Stderr, &err)
	return stdout, &out, stderr, &err
}

func TestMsg(t *testing.T) {
	stdout, _, stderr, _ := getWriters()
	lg := NewLogger(stdout, stderr, true)
	lg.Msg("hello %s", world)
}

func TestErr(t *testing.T) {
	stdout, _, stderr, _ := getWriters()
	lg := NewLogger(stdout, stderr, true)
	lg.Err("hello %s", world)
}

func TestInfo(t *testing.T) {
	stdout, _, stderr, err := getWriters()
	lg := NewLogger(stdout, stderr, true)
	lg.Info("hello %s", world)
	assert.Equal(t, err.Len(), 0, "stderr buffer length")
}

func TestDebug(t *testing.T) {
	stdout, out, stderr, err := getWriters()
	lgOn := NewLogger(stdout, stderr, true)
	lgOff := NewLogger(stdout, stderr, false)
	lgOff.Debug("hello %s", world)
	assert.Equal(t, out.Len(), 0, "stdout buffer length")
	assert.Equal(t, err.Len(), 0, "stderr buffer length")
	lgOn.Debug("hello %s", world)
	assert.Equal(t, err.Len(), 0, "stderr buffer length")
	assert.True(t, out.Len() > 0, "stdout buffer was empty, expected it to have debug output")
}
