// Package log is a convenience wrapper for logging messages of various levels (associated colors to come)
// to the terminal. Much of this code has been shamelessly stolen from https://github.com/helm/helm/blob/master/log/log.go
package log

import (
	"fmt"
	"io"
	"os"

	"github.com/deis/pkg/prettyprint"
)

// Stdout is the logging destination for normal messages.
var Stdout io.Writer = os.Stdout

// Stderr is the logging destination for error messages.
var Stderr io.Writer = os.Stderr

// IsDebugging toggles whether or not to enable debug output and behavior.
var IsDebugging = false

// Color is the representation of a color, to be used in Colorize
type Color string

// String is the fmt.Stringer interface implementation
func (c Color) String() string {
	return string(c)
}

var (
	// Default resets the console color
	Default = Color(prettyprint.Colors["Default"])
	// Red sets the console color to red
	Red = Color(prettyprint.Colors["Red"])
	// Cyan sets the console color to cyan
	Cyan = Color(prettyprint.Colors["Cyan"])
	// Yellow sets the console color to yellow
	Yellow = Color(prettyprint.Colors["Yellow"])
	// Green sets the console color to green
	Green = Color(prettyprint.Colors["Green"])
)

type Logger struct {
	stdout io.Writer
	stderr io.Writer
}

func NewLogger(stdout, stderr io.Writer) *Logger {
	return &Logger{stdout: stdout, stderr: stderr}
}

var defaultLogger = &Logger{stdout: os.Stdout, stderr: os.Stderr}

// Msg passes through the formatter, but otherwise prints exactly as-is.
//
// No prettification.
func (l *Logger) Msg(format string, v ...interface{}) {
	fmt.Fprintf(l.stdout, appendNewLine(format), v...)
}

// Die prints an error and then call os.Exit(1).
func (l *Logger) Die(format string, v ...interface{}) {
	l.Err(format, v...)
	if IsDebugging {
		panic(fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

// CleanExit prints a message and then exits with 0.
func (l *Logger) CleanExit(format string, v ...interface{}) {
	l.Info(format, v...)
	os.Exit(0)
}

// Err prints an error message. It does not cause an exit.
func (l *Logger) Err(format string, v ...interface{}) {
	fmt.Fprint(Stderr, addColor("[ERROR] ", Red))
	fmt.Fprintf(Stderr, appendNewLine(format), v...)
}

// Info prints a green-tinted message.
func (l *Logger) Info(format string, v ...interface{}) {
	fmt.Fprint(l.stderr, addColor("---> ", Green))
	fmt.Fprintf(l.stderr, appendNewLine(format), v...)
}

// Debug prints a cyan-tinted message if IsDebugging is true.
func (l *Logger) Debug(msg string, v ...interface{}) {
	if IsDebugging {
		fmt.Fprint(l.stderr, addColor("[DEBUG] ", Cyan))
		l.Msg(msg, v...)
	}
}

// Warn prints a yellow-tinted warning message.
func (l *Logger) Warn(format string, v ...interface{}) {
	fmt.Fprint(l.stderr, addColor("[WARN] ", Yellow))
	l.Msg(format, v...)
}

func appendNewLine(format string) string {
	return format + "\n"
}

func addColor(str string, color Color) string {
	return prettyprint.Colorize(fmt.Sprintf("{{.%s}}%s{{.%s}}", color.String(), str, Default.String()))
}
