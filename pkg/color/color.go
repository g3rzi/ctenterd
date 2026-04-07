package color

import (
	"fmt"
	"os"
)

var enabled = detectColor()

func detectColor() bool {
	// honor NO_COLOR (https://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// crude TTY check; you could use "golang.org/x/term" IsTerminal
	fi, _ := os.Stdout.Stat()
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func Enable()  { enabled = true }
func Disable() { enabled = false }
func Enabled() bool { return enabled }

// raw ANSI sequences
const (
	reset = "\x1b[0m"

	// basic colors
	Red    = "\x1b[31m"
	Green  = "\x1b[32m"
	Yellow = "\x1b[33m"
	Blue   = "\x1b[34m"
	Magenta= "\x1b[35m"
	Cyan   = "\x1b[36m"
	Gray   = "\x1b[90m"

	// styles
	Bold   = "\x1b[1m"
)

func wrap(code, s string) string {
	if !enabled {
		return s
	}
	return code + s + reset
}

// Convenience helpers
func Prompt(s string) string   { return wrap(Bold+Cyan, s) }
func Info(s string) string     { return wrap(Cyan, s) }
func Warn(s string) string     { return wrap(Yellow, s) }
func Error(s string) string    { return wrap(Red, s) }
func Ok(s string) string       { return wrap(Green, s) }
func Dim(s string) string      { return wrap(Gray, s) }

// Generic formatters
func F(code, format string, a ...any) string { return wrap(code, fmt.Sprintf(format, a...)) }
