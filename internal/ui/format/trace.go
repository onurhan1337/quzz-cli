package format

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/onurhan1337/quzz-cli/internal/trace"
)

// Common color definitions
var (
	Error   = color.New(color.FgRed, color.Bold)
	Warning = color.New(color.FgYellow)
	Info    = color.New(color.FgCyan)
	Success = color.New(color.FgGreen)
	Dim     = color.New(color.Faint)
)

// Icons for status display
const (
	IconCheck   = "✓"
	IconCross   = "✗"
	IconWarning = "⚠"
	IconInfo    = "ℹ"
	IconDash    = "−"
)

// Duration formats a duration in milliseconds with color coding
// Red for slow (>1000ms), yellow for moderate (>500ms), default otherwise
func Duration(duration int) string {
	durStr := fmt.Sprintf("%dms", duration)
	if duration > 1000 {
		return Error.Sprint(durStr)
	} else if duration > 500 {
		return Warning.Sprint(durStr)
	}
	return durStr
}

// Milliseconds formats milliseconds with appropriate color and units
// Converts to seconds if > 1000ms
func Milliseconds(ms int64) string {
	if ms > 1000 {
		// Use a custom orange color for large durations in stats
		orange := color.New(color.FgHiYellow)
		return orange.Sprintf("%.2fs", float64(ms)/1000.0)
	}
	return fmt.Sprintf("%dms", ms)
}

// Level formats a log level string with appropriate color
func Level(level string) string {
	levelUpper := strings.ToUpper(level)
	switch strings.ToLower(level) {
	case "error":
		return Error.Sprint(levelUpper)
	case "warn":
		return Warning.Sprint(levelUpper)
	case "debug":
		return Dim.Sprint(levelUpper)
	case "info":
		return Info.Sprint(levelUpper)
	default:
		return levelUpper
	}
}

// LevelColored formats a log level with an icon and color
func LevelColored(level string) string {
	switch level {
	case "error":
		return Error.Sprintf("%s ERROR", IconCross)
	case "warn":
		return Warning.Sprintf("%s WARN ", IconWarning)
	case "info":
		return Info.Sprintf("%s INFO ", IconInfo)
	case "debug":
		return Dim.Sprintf("%s DEBUG", IconDash)
	default:
		return fmt.Sprintf("%s %s", IconDash, level)
	}
}

// Operation formats an operation string, showing a dash if empty
func Operation(operation string) string {
	if operation == "" {
		return Dim.Sprint(IconDash)
	}
	return operation
}

// Status formats the status of a trace entry based on errors and warnings
func Status(t trace.TraceEntry) string {
	if t.Error != "" {
		return Error.Sprintf("%s ERROR", IconCross)
	}
	if t.Performance != nil && t.Performance.IsWarning {
		return Warning.Sprintf("%s SLOW", IconWarning)
	}
	return Success.Sprintf("%s OK", IconCheck)
}

// Component truncates a component name to fit the specified width
func Component(component string, maxWidth int) string {
	if len(component) <= maxWidth {
		return component
	}
	if maxWidth <= 3 {
		return component[:maxWidth]
	}
	return component[:maxWidth-3] + "..."
}
