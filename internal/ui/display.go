package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

const (
	boxWidth = 80
)

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

// Box drawing characters
const (
	topLeft     = "╭"
	topRight    = "╮"
	bottomLeft  = "╰"
	bottomRight = "╯"
	horizontal  = "─"
	vertical    = "│"
	vertBar     = "┃"
	teeRight    = "├"
	teeLeft     = "┤"
	cross       = "┼"
)

func PrintTitle(title string) {
	width := runeLen(title) + 2
	fmt.Println()
	Primary.Print(topLeft + strings.Repeat(horizontal, width) + topRight)
	fmt.Println()
	Primary.Print(vertical + " ")
	Bold.Print(title)
	Primary.Print(" " + vertical)
	fmt.Println()
	Primary.Print(bottomLeft + strings.Repeat(horizontal, width) + bottomRight)
	fmt.Println()
}

func PrintBlockStart(title string) {
	fmt.Println()
	Primary.Print(topLeft + strings.Repeat(horizontal, boxWidth-2) + topRight)
	fmt.Println()
	if title != "" {
		Primary.Print(vertical + " ")
		Bold.Print(title)
		titleLen := runeLen(title)
		padding := boxWidth - titleLen - 4
		Primary.Print(strings.Repeat(" ", padding) + " " + vertical)
		fmt.Println()
		Primary.Print(teeRight + strings.Repeat(horizontal, boxWidth-2) + teeLeft)
		fmt.Println()
	}
}

func PrintBlockEnd() {
	Primary.Print(bottomLeft + strings.Repeat(horizontal, boxWidth-2) + bottomRight)
	fmt.Println()
}

func PrintBlockLine(content string) {
	stripped := stripAnsi(content)
	contentLen := len(stripped)
	padding := boxWidth - contentLen - 4
	if padding < 0 {
		padding = 0
	}
	Primary.Print(vertical + " ")
	fmt.Print(content)
	Primary.Print(strings.Repeat(" ", padding) + " " + vertical)
	fmt.Println()
}

func PrintBlockSection(title string) {
	fmt.Println()
	sectionTitle := IconBullet + " " + title
	titleLen := runeLen(sectionTitle)
	padding := boxWidth - titleLen - 4
	Primary.Print(vertical + " ")
	Accent.Print(sectionTitle)
	Primary.Print(strings.Repeat(" ", padding) + " " + vertical)
	fmt.Println()
	Primary.Print(vertical + " ")
	Dim.Print(strings.Repeat(horizontal, boxWidth-4))
	Primary.Print(" " + vertical)
	fmt.Println()
}

func PrintBlockItem(content string) {
	stripped := stripAnsi(content)
	contentLen := len(stripped)
	padding := boxWidth - contentLen - 6
	if padding < 0 {
		padding = 0
	}
	Primary.Print(vertical + "   ")
	fmt.Print(content)
	Primary.Print(strings.Repeat(" ", padding) + " " + vertical)
	fmt.Println()
}

func PrintSection(title string) {
	fmt.Println()
	Primary.Print(vertBar + " ")
	Bold.Print(title)
	fmt.Println()
	Primary.Print(vertBar)
	Dim.Print(strings.Repeat(horizontal, boxWidth-2))
	fmt.Println()
}

func PrintSectionItem(label, value string) {
	fmt.Printf("  %-30s %s\n", Dim.Sprint(label+":"), value)
}

func PrintListItem(bullet, content string) {
	fmt.Printf("  %s %s\n", Dim.Sprint(bullet+"."), content)
}

// PanelStyle defines the visual style of a panel
type PanelStyle struct {
	Color *color.Color
	Icon  string
}

var (
	// InfoPanelStyle for informational panels
	InfoPanelStyle = PanelStyle{Color: Info, Icon: IconInfo}
	// SuccessPanelStyle for success panels
	SuccessPanelStyle = PanelStyle{Color: Success, Icon: IconCheck}
	// WarningPanelStyle for warning panels
	WarningPanelStyle = PanelStyle{Color: Warning, Icon: IconWarning}
)

// PrintPanel prints a styled panel with title and items
func PrintPanel(style PanelStyle, title string, items []string) {
	fmt.Println()
	style.Color.Print(topLeft + strings.Repeat(horizontal, boxWidth-2) + topRight)
	fmt.Println()

	panelTitle := style.Icon + " " + title
	titleLen := runeLen(panelTitle)
	titlePadding := boxWidth - titleLen - 4
	style.Color.Print(vertical + " ")
	Bold.Print(panelTitle)
	style.Color.Print(strings.Repeat(" ", titlePadding) + " " + vertical)
	fmt.Println()

	for _, item := range items {
		lines := wrapText(item, boxWidth-4)
		for _, line := range lines {
			stripped := stripAnsi(line)
			contentLen := len(stripped)
			padding := boxWidth - contentLen - 4
			if padding < 0 {
				padding = 0
			}
			style.Color.Print(vertical + " ")
			fmt.Print(line)
			style.Color.Print(strings.Repeat(" ", padding) + " " + vertical)
			fmt.Println()
		}
	}

	style.Color.Print(bottomLeft + strings.Repeat(horizontal, boxWidth-2) + bottomRight)
	fmt.Println()
}

// PrintInfoPanel prints an informational panel (convenience wrapper)
func PrintInfoPanel(title string, items []string) {
	PrintPanel(InfoPanelStyle, title, items)
}

// PrintSuccessPanel prints a success panel (convenience wrapper)
func PrintSuccessPanel(title string, items []string) {
	PrintPanel(SuccessPanelStyle, title, items)
}

// PrintWarningPanel prints a warning panel (convenience wrapper)
func PrintWarningPanel(title string, items []string) {
	PrintPanel(WarningPanelStyle, title, items)
}

func PrintDivider() {
	Dim.Println(strings.Repeat(horizontal, boxWidth))
}

func stripAnsi(str string) string {
	result := ""
	inEscape := false
	for _, r := range str {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(r)
	}
	return result
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	stripped := stripAnsi(text)
	if len(stripped) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	currentLine := ""
	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(stripAnsi(testLine)) <= width {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
