package wizard

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/onurhan1337/quzz-cli/internal/config"
)

// Step represents a step in the initialization wizard
type Step int

const (
	StepFormat Step = iota
	StepMode
	StepLogLevel
	StepOutputFormat
	StepThreshold
	StepPerformance
	StepComponentFilter
	StepConfirm
	StepCreating
	StepSuccess
)

// Result contains the result of the initialization wizard
type Result struct {
	Config   config.QuzzConfig
	Format   config.ConfigFormat
	Minimal  bool
	FilePath string
}

// Model is the Bubble Tea model for the initialization wizard
type Model struct {
	step        Step
	format      config.ConfigFormat
	minimal     bool
	logLevel    string
	outputFmt   string
	threshold   int
	enablePerf  bool
	compFilter  string
	list        list.Model
	textInput   textinput.Model
	spinner     spinner.Model
	result      *Result
	width       int
	height      int
	err         error
}

// item is a list item for the wizard selections
type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// creationCompleteMsg is sent when configuration creation is complete
type creationCompleteMsg struct {
	result   *Result
	duration time.Duration
}

// ExtractResult extracts the Result from a tea.Model
func ExtractResult(m tea.Model) *Result {
	if model, ok := m.(Model); ok {
		return model.result
	}
	return nil
}

// New creates a new initialization wizard model
func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.CharLimit = 100

	items := []list.Item{
		item{title: "JavaScript (.js)", desc: "Standard Node.js configuration"},
		item{title: "TypeScript (.ts)", desc: "Type-safe configuration"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select config file format"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return Model{
		step:      StepFormat,
		list:      l,
		textInput: ti,
		spinner:   sp,
		threshold: config.DefaultWarnThreshold,
	}
}

// Init initializes the wizard model
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}
