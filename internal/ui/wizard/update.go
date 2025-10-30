package wizard

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/onurhan1337/quzz-cli/internal/config"
)

// Update handles messages and updates the wizard state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.step == StepSuccess {
			return m, tea.Quit
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.step > StepFormat {
				m.step--
				m.updateListForStep()
			} else {
				return m, tea.Quit
			}
			return m, nil
		case "enter":
			return m.handleEnter()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := 60, 10
		m.list.SetSize(h, v)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case creationCompleteMsg:
		m.result = msg.result
		m.step = StepSuccess
		return m, nil
	}

	var cmd tea.Cmd
	if m.step == StepThreshold || m.step == StepComponentFilter {
		m.textInput, cmd = m.textInput.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

// updateListForStep updates the list items based on the current step
func (m *Model) updateListForStep() {
	var items []item
	var title string

	switch m.step {
	case StepFormat:
		title = "Select config file format"
		items = []item{
			{title: "JavaScript (.js)", desc: "Standard Node.js configuration"},
			{title: "TypeScript (.ts)", desc: "Type-safe configuration"},
		}

	case StepMode:
		title = "Select configuration mode"
		items = []item{
			{title: "Minimal", desc: "Quick setup with defaults (recommended)"},
			{title: "Full", desc: "Complete control over all options"},
		}

	case StepLogLevel:
		title = "Choose log level"
		items = []item{
			{title: "error", desc: "Only errors"},
			{title: "warn", desc: "Errors and warnings"},
			{title: "info", desc: "Standard logging (recommended)"},
			{title: "debug", desc: "Verbose debug output"},
		}

	case StepOutputFormat:
		title = "Choose output format"
		items = []item{
			{title: "pretty", desc: "Rich formatting with colors"},
			{title: "compact", desc: "Dense terminal output"},
			{title: "json", desc: "Machine-readable format"},
		}

	case StepPerformance:
		title = "Enable performance monitoring"
		items = []item{
			{title: "No", desc: "Disabled (production)"},
			{title: "Yes", desc: "Enabled (development)"},
		}

	case StepConfirm:
		title = "Confirm configuration"
		items = []item{
			{title: "Create", desc: "Generate configuration file"},
			{title: "Cancel", desc: "Abort setup"},
		}
	}

	m.list.Title = title
	// Convert items to []list.Item
	listItems := make([]list.Item, len(items))
	for i, it := range items {
		listItems[i] = it
	}
	m.list.SetItems(listItems)
}

// handleEnter processes the enter key based on the current step
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepFormat:
		selected := m.list.SelectedItem().(item)
		if strings.Contains(selected.title, "TypeScript") {
			m.format = config.ConfigFormatTS
		} else {
			m.format = config.ConfigFormatJS
		}
		m.step = StepMode
		m.updateListForStep()

	case StepMode:
		selected := m.list.SelectedItem().(item)
		m.minimal = selected.title == "Minimal"
		m.step = StepLogLevel
		m.updateListForStep()

	case StepLogLevel:
		selected := m.list.SelectedItem().(item)
		m.logLevel = selected.title
		m.step = StepOutputFormat
		m.updateListForStep()

	case StepOutputFormat:
		selected := m.list.SelectedItem().(item)
		m.outputFmt = selected.title

		if m.minimal {
			m.step = StepConfirm
			m.updateListForStep()
		} else {
			m.step = StepThreshold
			m.textInput.SetValue(fmt.Sprintf("%d", config.DefaultWarnThreshold))
			m.textInput.Focus()
		}

	case StepThreshold:
		var threshold int
		fmt.Sscanf(m.textInput.Value(), "%d", &threshold)
		if threshold > 0 {
			m.threshold = threshold
		}
		m.step = StepPerformance
		m.textInput.Blur()
		m.updateListForStep()

	case StepPerformance:
		selected := m.list.SelectedItem().(item)
		m.enablePerf = selected.title == "Yes"
		m.step = StepComponentFilter
		m.textInput.SetValue("")
		m.textInput.Focus()

	case StepComponentFilter:
		m.compFilter = m.textInput.Value()
		m.step = StepConfirm
		m.textInput.Blur()
		m.updateListForStep()

	case StepConfirm:
		selected := m.list.SelectedItem().(item)
		if selected.title == "Create" {
			m.step = StepCreating
			return m, m.createConfig()
		}
		return m, tea.Quit
	}

	return m, nil
}

// createConfig creates the configuration file
func (m Model) createConfig() tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		cfg := m.buildConfig()
		opts := config.GeneratorOptions{
			Format:  m.format,
			Config:  cfg,
			Minimal: m.minimal,
		}
		filePath, err := config.GenerateConfigFile(opts)
		if err != nil {
			return err
		}

		duration := time.Since(start)

		return creationCompleteMsg{
			result: &Result{
				Config:   cfg,
				Format:   m.format,
				Minimal:  m.minimal,
				FilePath: filePath,
			},
			duration: duration,
		}
	}
}

// buildConfig builds the configuration from wizard selections
func (m Model) buildConfig() config.QuzzConfig {
	if m.minimal {
		cfg := config.MinimalConfig()
		cfg.LogLevel = m.logLevel
		cfg.OutputFormat = m.outputFmt
		return cfg
	}

	cfg := config.DefaultConfig()
	cfg.LogLevel = m.logLevel
	cfg.OutputFormat = m.outputFmt
	cfg.Performance.Enabled = m.enablePerf
	cfg.Performance.WarnThreshold = m.threshold
	cfg.ComponentFilter = m.compFilter
	return cfg
}
