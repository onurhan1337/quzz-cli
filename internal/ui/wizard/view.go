package wizard

import (
	"strings"

	"github.com/onurhan1337/quzz-cli/internal/ui"
)

// View renders the wizard UI
func (m Model) View() string {
	if m.step == StepCreating {
		return m.renderCreating()
	}

	if m.step == StepSuccess {
		return m.renderSuccess()
	}

	var b strings.Builder

	b.WriteString("\n ")
	b.WriteString(ui.Title("QUZZ CONFIGURATION SETUP"))
	b.WriteString("\n\n")

	switch m.step {
	case StepThreshold:
		b.WriteString(" ")
		b.WriteString(ui.Orange2.Sprint("Performance Warning Threshold (ms)"))
		b.WriteString("\n ")
		b.WriteString(ui.MutedText("Components slower than this will be flagged"))
		b.WriteString("\n\n ")
		b.WriteString(m.textInput.View())
		b.WriteString("\n")

	case StepComponentFilter:
		b.WriteString(" ")
		b.WriteString(ui.Orange2.Sprint("Component Filter (regex pattern)"))
		b.WriteString("\n ")
		b.WriteString(ui.MutedText("Leave empty to trace all components"))
		b.WriteString("\n ")
		b.WriteString(ui.MutedText("Example: ^(User|Product) for User* and Product* components"))
		b.WriteString("\n\n ")
		b.WriteString(m.textInput.View())
		b.WriteString("\n")

	case StepConfirm:
		b.WriteString(m.renderSummary())
		b.WriteString("\n\n")
		b.WriteString(m.list.View())

	default:
		b.WriteString(m.list.View())
	}

	b.WriteString("\n\n ")

	helpItems := []string{
		ui.Orange2.Sprint("enter") + ui.Gray.Sprint(" select"),
		ui.Orange2.Sprint("esc") + ui.Gray.Sprint(" back"),
		ui.Orange2.Sprint("ctrl+c") + ui.Gray.Sprint(" quit"),
	}
	if m.step == StepFormat {
		helpItems = []string{
			ui.Orange2.Sprint("enter") + ui.Gray.Sprint(" select"),
			ui.Orange2.Sprint("esc") + ui.Gray.Sprint(" quit"),
		}
	}
	b.WriteString(strings.Join(helpItems, "  "))
	b.WriteString("\n\n")

	return b.String()
}

// renderCreating renders the creating configuration view
func (m Model) renderCreating() string {
	var b strings.Builder

	b.WriteString("\n\n  ")
	b.WriteString(m.spinner.View())
	b.WriteString(" ")
	b.WriteString(ui.Orange1.Sprint("Creating Configuration"))
	b.WriteString("\n\n  ")
	b.WriteString(ui.MutedText("Generating your Quzz configuration file..."))
	b.WriteString("\n\n")

	return b.String()
}

// renderSuccess renders the success view
func (m Model) renderSuccess() string {
	var b strings.Builder

	b.WriteString("\n\n ")
	b.WriteString(ui.Success.Sprint("Configuration Created Successfully!"))
	b.WriteString("\n\n")

	if m.result != nil {
		b.WriteString(" ")
		b.WriteString(ui.Orange1.Sprint("CONFIGURATION DETAILS"))
		b.WriteString("\n ")
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n\n")

		b.WriteString("  ")
		b.WriteString(ui.Label("File:       "))
		b.WriteString(ui.ValueText(m.result.FilePath))
		b.WriteString("\n  ")
		b.WriteString(ui.Label("Format:     "))
		b.WriteString(ui.Info.Sprint(m.result.Format.String()))
		b.WriteString("\n  ")
		b.WriteString(ui.Label("Mode:       "))
		b.WriteString(ui.Orange2.Sprint(map[bool]string{true: "Minimal", false: "Full"}[m.result.Minimal]))
		b.WriteString("\n  ")
		b.WriteString(ui.Label("Log Level:  "))
		b.WriteString(ui.Info.Sprint(m.result.Config.LogLevel))
		b.WriteString("\n  ")
		b.WriteString(ui.Label("Output:     "))
		b.WriteString(ui.Info.Sprint(m.result.Config.OutputFormat))
		b.WriteString("\n")

		if !m.result.Minimal {
			b.WriteString("  ")
			b.WriteString(ui.Label("Threshold:  "))
			b.WriteString(ui.Warning.Sprintf("%dms", m.result.Config.Performance.WarnThreshold))
			b.WriteString("\n  ")
			b.WriteString(ui.Label("Performance:"))
			if m.result.Config.Performance.Enabled {
				b.WriteString(ui.Success.Sprint("Enabled"))
			} else {
				b.WriteString(ui.MutedText("Disabled"))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n ")
		b.WriteString(ui.Orange1.Sprint("NEXT STEPS"))
		b.WriteString("\n ")
		b.WriteString(strings.Repeat("─", 60))
		b.WriteString("\n\n")
		b.WriteString("  ")
		b.WriteString(ui.MutedText("1. Review your configuration file"))
		b.WriteString("\n  ")
		b.WriteString(ui.MutedText("2. Import quzz in your React components"))
		b.WriteString("\n  ")
		b.WriteString(ui.MutedText("3. Wrap components with quzz() HOC"))
		b.WriteString("\n  ")
		b.WriteString(ui.MutedText("4. Run your app and generate traces"))
		b.WriteString("\n  ")
		b.WriteString(ui.Orange2.Sprint("5. Visualize: "))
		b.WriteString(ui.Orange3.Sprint("quzz visualize traces.json"))
		b.WriteString("\n")
	}

	b.WriteString("\n ")
	b.WriteString(ui.MutedText("Press any key to exit"))
	b.WriteString("\n\n")

	return b.String()
}

// renderSummary renders the configuration summary
func (m Model) renderSummary() string {
	var b strings.Builder

	b.WriteString(" ")
	b.WriteString(ui.Orange2.Sprint("Configuration Summary"))
	b.WriteString("\n ")
	b.WriteString(strings.Repeat("─", 50))
	b.WriteString("\n\n  ")
	b.WriteString(ui.Label("Format:     "))
	b.WriteString(ui.Info.Sprint(m.format.String()))
	b.WriteString("\n  ")
	b.WriteString(ui.Label("Mode:       "))
	b.WriteString(ui.Orange2.Sprint(map[bool]string{true: "Minimal", false: "Full"}[m.minimal]))
	b.WriteString("\n  ")
	b.WriteString(ui.Label("Log Level:  "))
	b.WriteString(ui.Info.Sprint(m.logLevel))
	b.WriteString("\n  ")
	b.WriteString(ui.Label("Output:     "))
	b.WriteString(ui.Info.Sprint(m.outputFmt))
	b.WriteString("\n")

	if !m.minimal {
		b.WriteString("  ")
		b.WriteString(ui.Label("Threshold:  "))
		b.WriteString(ui.Warning.Sprintf("%dms", m.threshold))
		b.WriteString("\n  ")
		b.WriteString(ui.Label("Performance:"))
		if m.enablePerf {
			b.WriteString(ui.Success.Sprint("Enabled"))
		} else {
			b.WriteString(ui.MutedText("Disabled"))
		}
		b.WriteString("\n")

		if m.compFilter != "" {
			b.WriteString("  ")
			b.WriteString(ui.Label("Filter:     "))
			b.WriteString(ui.ValueText(m.compFilter))
			b.WriteString("\n")
		}
	}

	b.WriteString(" ")
	b.WriteString(strings.Repeat("─", 50))
	b.WriteString("\n")

	return b.String()
}
