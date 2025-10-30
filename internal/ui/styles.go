package ui

import "fmt"

// Title formats text as a title
func Title(text string) string {
	return Orange1.Sprint(text)
}

// HighlightText formats text as highlighted
func HighlightText(text string) string {
	return Orange2.Sprint(text)
}

// Label formats text as a label
func Label(text string) string {
	return Gray.Sprint(text)
}

// ValueText formats text as a value
func ValueText(text string) string {
	return White.Sprint(text)
}

// MutedText formats text as muted/dim
func MutedText(text string) string {
	return Dim.Sprint(text)
}

// RenderProgressBar renders a progress bar with the given percentage and width
func RenderProgressBar(percent float64, width int) string {
	if width < 10 {
		width = 10
	}
	filled := int(float64(width) * percent)
	if filled > width {
		filled = width
	}
	empty := width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += Black.Sprint("█")
	}
	for i := 0; i < empty; i++ {
		bar += Gray.Sprint("░")
	}
	return fmt.Sprintf("%s %s", bar, Dim.Sprintf("%0.0f%%", percent*100))
}