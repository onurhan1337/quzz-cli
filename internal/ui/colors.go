package ui

import "github.com/fatih/color"

var (
	Orange1 = color.New(color.FgHiYellow, color.Bold)
	Orange2 = color.New(color.FgYellow)
	Orange3 = color.New(color.FgHiYellow)
	Orange4 = color.New(color.FgYellow)

	White = color.New(color.FgWhite)
	Gray  = color.New(color.FgHiBlack)
	Black = color.New(color.Faint)
	Dim   = color.New(color.Faint)
	Bold  = color.New(color.Bold)

	Success = color.New(color.FgGreen, color.Bold)
	Error   = color.New(color.FgHiYellow, color.Bold)
	Warning = color.New(color.FgYellow, color.Bold)
	Info    = color.New(color.FgHiYellow)

	Highlight = color.New(color.FgHiYellow)
	Value     = color.New(color.FgWhite)
	Muted     = color.New(color.FgHiBlack)

	Primary = Orange1
	Accent  = Orange2
)

const (
	IconBullet  = "•"
	IconInfo    = "ℹ"
	IconCheck   = "✓"
	IconWarning = "⚠"
	IconCross   = "✗"
	IconDash    = "-"
	IconArrow   = "→"
)
