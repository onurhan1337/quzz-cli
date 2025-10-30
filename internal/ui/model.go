package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/onurhan1337/quzz-cli/internal/trace"
	"github.com/onurhan1337/quzz-cli/internal/ui/format"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Help     key.Binding
	Quit     key.Binding
	Toggle   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "b"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "f"),
	),
	Home: key.NewBinding(
		key.WithKeys("home", "g"),
	),
	End: key.NewBinding(
		key.WithKeys("end", "G"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("tab"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
	),
}

type model struct {
	traces      []trace.TraceEntry
	stats       trace.TraceStats
	viewport    viewport.Model
	spinner     spinner.Model
	ready       bool
	loading     bool
	showStats   bool
	showHelp    bool
	width       int
	height      int
}

type loadCompleteMsg struct{}

func NewModel(traces []trace.TraceEntry, stats trace.TraceStats) model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return model{
		traces:    traces,
		stats:     stats,
		showStats: false,
		showHelp:  false,
		spinner:   sp,
		loading:   true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadData(),
	)
}

func loadData() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return loadCompleteMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case loadCompleteMsg:
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		if m.showHelp {
			m.showHelp = false
			return m, nil
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Toggle):
			m.showStats = !m.showStats
			m.viewport.SetContent(m.renderContent())
			return m, nil
		case key.Matches(msg, keys.Help):
			m.showHelp = true
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 4
		footerHeight := 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.renderContent())
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.viewport.SetContent(m.renderContent())
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if !m.ready {
		return MutedText("\n  Initializing...")
	}

	if m.showHelp {
		return m.renderHelpView()
	}

	return fmt.Sprintf("%s\n%s\n%s", m.renderHeader(), m.viewport.View(), m.renderFooter())
}

func (m model) renderLoading() string {
	return fmt.Sprintf("\n\n  %s %s\n\n  %s\n\n",
		m.spinner.View(),
		Orange1.Sprint("Loading Traces"),
		MutedText("Analyzing your trace data..."))
}

func (m model) renderHeader() string {
	var b strings.Builder

	b.WriteString("\n ")
	b.WriteString(Title("QUZZ TRACE VISUALIZATION"))
	b.WriteString("\n ")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n ")

	b.WriteString(Gray.Sprintf("Traces: "))
	b.WriteString(Orange2.Sprintf("%d", len(m.traces)))
	b.WriteString(Gray.Sprint(" | Duration: "))
	b.WriteString(format.Milliseconds(m.stats.TotalDuration))

	if m.showStats {
		b.WriteString(Success.Sprint(" | [Stats Mode]"))
	}
	if m.stats.ErrorCount > 0 {
		b.WriteString(Gray.Sprint(" | Errors: "))
		b.WriteString(Error.Sprintf("%d", m.stats.ErrorCount))
	}
	if m.stats.WarningCount > 0 {
		b.WriteString(Gray.Sprint(" | Warnings: "))
		b.WriteString(Warning.Sprintf("%d", m.stats.WarningCount))
	}

	b.WriteString("\n")

	return b.String()
}

func (m model) renderFooter() string {
	var b strings.Builder

	percent := m.viewport.ScrollPercent()
	progressBar := RenderProgressBar(percent, (m.width - 4) / 8)

	b.WriteString("\n ")
	b.WriteString(progressBar)
	b.WriteString("\n ")

	b.WriteString(Orange2.Sprint("↑/↓ j/k"))
	b.WriteString(Gray.Sprint(" scroll "))
	b.WriteString(Orange2.Sprint("pgup/pgdn"))
	b.WriteString(Gray.Sprint(" page "))
	b.WriteString(Orange2.Sprint("tab"))
	b.WriteString(Gray.Sprint(" stats "))
	b.WriteString(Orange2.Sprint("?"))
	b.WriteString(Gray.Sprint(" help "))
	b.WriteString(Orange2.Sprint("q"))
	b.WriteString(Gray.Sprint(" quit"))

	b.WriteString("\n")

	return b.String()
}

func (m model) renderContent() string {
	var content strings.Builder

	if m.showStats {
		content.WriteString(m.renderStats())
		content.WriteString("\n\n")
	}

	content.WriteString(m.renderTraces())

	return content.String()
}

func (m model) renderStats() string {
	var b strings.Builder

	b.WriteString("\n ")
	b.WriteString(Orange1.Sprint("STATISTICS"))
	b.WriteString("\n ")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n\n ")

	b.WriteString(Orange2.Sprint("Overview:"))
	b.WriteString("\n   ")
	b.WriteString(Label("Total Traces:       "))
	b.WriteString(ValueText(fmt.Sprintf("%d", m.stats.TotalTraces)))
	b.WriteString("\n   ")
	b.WriteString(Label("Total Duration:     "))
	b.WriteString(format.Milliseconds(m.stats.TotalDuration))
	b.WriteString("\n   ")
	b.WriteString(Label("Average Duration:   "))
	b.WriteString(format.Milliseconds(int64(m.stats.AverageDuration)))
	b.WriteString("\n   ")
	b.WriteString(Label("Min Duration:       "))
	b.WriteString(Info.Sprintf("%dms", m.stats.MinDuration))
	b.WriteString("\n   ")
	b.WriteString(Label("Max Duration:       "))
	b.WriteString(Warning.Sprintf("%dms", m.stats.MaxDuration))

	if m.stats.ErrorCount > 0 || m.stats.WarningCount > 0 {
		b.WriteString("\n\n ")
		b.WriteString(Orange2.Sprint("Issues:"))
		b.WriteString("\n")
		if m.stats.ErrorCount > 0 {
			b.WriteString("   ")
			b.WriteString(Label("Errors:   "))
			b.WriteString(Error.Sprintf("%d", m.stats.ErrorCount))
			b.WriteString("\n")
		}
		if m.stats.WarningCount > 0 {
			b.WriteString("   ")
			b.WriteString(Label("Warnings: "))
			b.WriteString(Warning.Sprintf("%d", m.stats.WarningCount))
			b.WriteString("\n")
		}
	}

	if len(m.stats.ComponentCounts) > 0 {
		b.WriteString("\n ")
		b.WriteString(Orange2.Sprint("Top Components:"))
		b.WriteString("\n")
		count := 0
		for name, cnt := range m.stats.ComponentCounts {
			if count >= 5 {
				break
			}
			b.WriteString(fmt.Sprintf("   %-30s %s\n",
				ValueText(name),
				Gray.Sprintf("(%d traces)", cnt)))
			count++
		}
	}

	if len(m.stats.LevelCounts) > 0 {
		b.WriteString("\n ")
		b.WriteString(Orange2.Sprint("Log Levels:"))
		b.WriteString("\n")
		levels := []string{"error", "warn", "info", "debug"}
		for _, level := range levels {
			if count, ok := m.stats.LevelCounts[level]; ok {
				var levelColor *color.Color
				switch level {
				case "error":
					levelColor = Error
				case "warn":
					levelColor = Warning
				case "info":
					levelColor = Info
				default:
					levelColor = Gray
				}
				b.WriteString(fmt.Sprintf("   %-10s %s\n",
					levelColor.Sprint(level),
					Gray.Sprintf("(%d)", count)))
			}
		}
	}

	b.WriteString("\n")
	return b.String()
}

func (m model) renderTraces() string {
	var b strings.Builder

	b.WriteString("\n ")
	b.WriteString(Orange1.Sprint("TRACES"))
	b.WriteString("\n ")
	b.WriteString(strings.Repeat("─", 80))
	b.WriteString("\n\n")

	if len(m.traces) == 0 {
		return b.String() + MutedText("  No traces to display\n")
	}

	columns := []table.Column{
		{Title: "TIME", Width: 10},
		{Title: "COMPONENT", Width: 25},
		{Title: "DURATION", Width: 10},
		{Title: "LEVEL", Width: 8},
		{Title: "OPERATION", Width: 15},
		{Title: "STATUS", Width: 10},
	}

	rows := []table.Row{}
	for _, t := range m.traces {
		timeStr := Gray.Sprint(t.Timestamp.Format("15:04:05"))
		componentStr := ValueText(format.Component(t.Component, 25))

		var durationStr string
		if t.Duration > 1000 {
			durationStr = Error.Sprintf("%dms", t.Duration)
		} else if t.Duration > 500 {
			durationStr = Warning.Sprintf("%dms", t.Duration)
		} else {
			durationStr = Success.Sprintf("%dms", t.Duration)
		}

		var levelStr string
		switch strings.ToLower(t.Level) {
		case "error":
			levelStr = Error.Sprint("ERROR")
		case "warn":
			levelStr = Warning.Sprint("WARN")
		case "info":
			levelStr = Info.Sprint("INFO")
		case "debug":
			levelStr = Gray.Sprint("DEBUG")
		default:
			levelStr = t.Level
		}

		operationStr := t.Operation
		if operationStr == "" {
			operationStr = Gray.Sprint("-")
		} else {
			operationStr = ValueText(operationStr)
		}

		var statusStr string
		if t.Error != "" {
			statusStr = Error.Sprint("ERROR")
		} else if t.Performance != nil && t.Performance.IsWarning {
			statusStr = Warning.Sprint("SLOW")
		} else {
			statusStr = Success.Sprint("OK")
		}

		rows = append(rows, table.Row{
			timeStr, componentStr, durationStr, levelStr, operationStr, statusStr,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows)),
	)

	b.WriteString(t.View())
	b.WriteString("\n")

	return b.String()
}

func (m model) renderHelpView() string {
	var b strings.Builder

	b.WriteString("\n\n ")
	b.WriteString(Title("KEYBOARD SHORTCUTS"))
	b.WriteString("\n ")
	b.WriteString(strings.Repeat("─", 40))
	b.WriteString("\n\n")

	helpItems := []struct {
		key  string
		desc string
	}{
		{"↑ / k", "Scroll up"},
		{"↓ / j", "Scroll down"},
		{"PgUp / b", "Page up"},
		{"PgDn / f", "Page down"},
		{"Home / g", "Jump to top"},
		{"End / G", "Jump to bottom"},
		{"Tab", "Toggle statistics"},
		{"?", "Toggle help"},
		{"q / Esc", "Quit"},
	}

	for _, item := range helpItems {
		b.WriteString(fmt.Sprintf("  %-15s %s\n",
			Orange2.Sprint(item.key),
			Gray.Sprint(item.desc)))
	}

	b.WriteString("\n ")
	b.WriteString(MutedText("Press any key to return"))
	b.WriteString("\n\n")

	return b.String()
}
