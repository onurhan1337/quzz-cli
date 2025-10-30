package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/onurhan1337/quzz-cli/internal/trace"
	"github.com/onurhan1337/quzz-cli/internal/ui"
	"github.com/spf13/cobra"
)

type visualizeFlags struct {
	statsOnly      bool
	jsonOutput     bool
	component      string
	operation      string
	category       string
	level          string
	minDuration    int
	maxDuration    int
	startDate      string
	endDate        string
	hasError       bool
	hasWarning     bool
	componentRegex string
	limit          int
}

var vFlags visualizeFlags

var visualizeCmd = &cobra.Command{
	Use:   "visualize [path]",
	Short: "Visualize and explore traces from traces.json",
	Long: `Visualize and explore traces with powerful filtering options.

Examples:
  quzz visualize traces.json
  quzz visualize traces.json --stats
  quzz visualize traces.json --component UserProfile
  quzz visualize traces.json --min-duration 100 --max-duration 500
  quzz visualize traces.json --errors
  quzz visualize traces.json --component-regex "^(Blog|Product)"
  quzz visualize traces.json --json`,
	Args: cobra.ExactArgs(1),
	RunE: runVisualize,
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	flags := visualizeCmd.Flags()
	flags.BoolVarP(&vFlags.statsOnly, "stats", "s", false, "Show only statistics")
	flags.BoolVarP(&vFlags.jsonOutput, "json", "j", false, "Output as JSON")
	flags.StringVarP(&vFlags.component, "component", "c", "", "Filter by component name")
	flags.StringVarP(&vFlags.operation, "operation", "o", "", "Filter by operation")
	flags.StringVar(&vFlags.category, "category", "", "Filter by category")
	flags.StringVarP(&vFlags.level, "level", "l", "", "Filter by log level")
	flags.IntVar(&vFlags.minDuration, "min-duration", 0, "Filter by minimum duration (ms)")
	flags.IntVar(&vFlags.maxDuration, "max-duration", 0, "Filter by maximum duration (ms)")
	flags.StringVar(&vFlags.startDate, "start-date", "", "Filter by start date (RFC3339)")
	flags.StringVar(&vFlags.endDate, "end-date", "", "Filter by end date (RFC3339)")
	flags.BoolVar(&vFlags.hasError, "errors", false, "Show only traces with errors")
	flags.BoolVar(&vFlags.hasWarning, "warnings", false, "Show only traces with warnings")
	flags.StringVar(&vFlags.componentRegex, "component-regex", "", "Filter by regex pattern")
	flags.IntVar(&vFlags.limit, "limit", 50, "Limit traces displayed (0 for all)")
}

func runVisualize(cmd *cobra.Command, args []string) error {
	tracePath := args[0]

	traceFile, err := loadTraceFile(tracePath)
	if err != nil {
		return err
	}

	filterOpts, err := buildFilterOptions()
	if err != nil {
		return err
	}

	filteredTraces := trace.FilterTraces(traceFile.Traces, filterOpts)

	if len(filteredTraces) == 0 {
		showNoResults()
		return nil
	}

	if vFlags.jsonOutput {
		return outputJSON(filteredTraces)
	}

	displayTraces(filteredTraces)
	return nil
}

func showNoResults() {
	ui.PrintInfoPanel("NO RESULTS", []string{
		"No traces match the specified filters",
		"",
		"Try adjusting your filter criteria:",
		"  - Remove or loosen duration constraints",
		"  - Check component/operation names for typos",
		"  - Verify date range includes trace data",
	})
}

func loadTraceFile(path string) (*trace.TraceFile, error) {
	traceFile, err := trace.LoadTraceFile(path)
	if err != nil {
		ui.PrintWarningPanel("ERROR", []string{
			fmt.Sprintf("Failed to load trace file: %v", err),
		})
		return nil, err
	}
	return traceFile, nil
}

func buildFilterOptions() (trace.FilterOptions, error) {
	opts := trace.FilterOptions{
		Component:      vFlags.component,
		Operation:      vFlags.operation,
		Category:       vFlags.category,
		Level:          vFlags.level,
		MinDuration:    vFlags.minDuration,
		MaxDuration:    vFlags.maxDuration,
		HasError:       vFlags.hasError,
		HasWarning:     vFlags.hasWarning,
		ComponentRegex: vFlags.componentRegex,
	}

	if vFlags.startDate != "" {
		t, err := time.Parse(time.RFC3339, vFlags.startDate)
		if err != nil {
			ui.PrintWarningPanel("INVALID DATE FORMAT", []string{
				fmt.Sprintf("Start date must be in RFC3339 format: %v", err),
			})
			return opts, err
		}
		opts.StartDate = t
	}

	if vFlags.endDate != "" {
		t, err := time.Parse(time.RFC3339, vFlags.endDate)
		if err != nil {
			ui.PrintWarningPanel("INVALID DATE FORMAT", []string{
				fmt.Sprintf("End date must be in RFC3339 format: %v", err),
			})
			return opts, err
		}
		opts.EndDate = t
	}

	return opts, nil
}

func outputJSON(traces []trace.TraceEntry) error {
	output := map[string]interface{}{
		"traces": traces,
		"stats":  trace.CalculateStats(traces),
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		ui.PrintWarningPanel("JSON ERROR", []string{
			fmt.Sprintf("Failed to generate JSON output: %v", err),
		})
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

func displayTraces(traces []trace.TraceEntry) {
	stats := trace.CalculateStats(traces)

	if vFlags.statsOnly {
		displayActiveFilters()
		displayStats(stats)
	} else {
		m := ui.NewModel(traces, stats)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			ui.PrintWarningPanel("VISUALIZATION ERROR", []string{
				fmt.Sprintf("Error running visualization: %v", err),
			})
			os.Exit(1)
		}
	}
}

func displayActiveFilters() {
	filters := []string{}

	if vFlags.component != "" {
		filters = append(filters, fmt.Sprintf("Component: %s", ui.Highlight.Sprint(vFlags.component)))
	}
	if vFlags.operation != "" {
		filters = append(filters, fmt.Sprintf("Operation: %s", ui.Highlight.Sprint(vFlags.operation)))
	}
	if vFlags.category != "" {
		filters = append(filters, fmt.Sprintf("Category: %s", ui.Highlight.Sprint(vFlags.category)))
	}
	if vFlags.level != "" {
		filters = append(filters, fmt.Sprintf("Level: %s", ui.Highlight.Sprint(vFlags.level)))
	}
	if vFlags.minDuration > 0 {
		filters = append(filters, fmt.Sprintf("Min Duration: %s", ui.Highlight.Sprintf("%dms", vFlags.minDuration)))
	}
	if vFlags.maxDuration > 0 {
		filters = append(filters, fmt.Sprintf("Max Duration: %s", ui.Highlight.Sprintf("%dms", vFlags.maxDuration)))
	}
	if vFlags.hasError {
		filters = append(filters, ui.Error.Sprint("Errors only"))
	}
	if vFlags.hasWarning {
		filters = append(filters, ui.Warning.Sprint("Warnings only"))
	}
	if vFlags.componentRegex != "" {
		filters = append(filters, fmt.Sprintf("Component Pattern: %s", ui.Highlight.Sprint(vFlags.componentRegex)))
	}
	if vFlags.startDate != "" {
		filters = append(filters, fmt.Sprintf("Start Date: %s", ui.Value.Sprint(vFlags.startDate)))
	}
	if vFlags.endDate != "" {
		filters = append(filters, fmt.Sprintf("End Date: %s", ui.Value.Sprint(vFlags.endDate)))
	}
	if vFlags.limit > 0 && vFlags.limit != 50 {
		filters = append(filters, fmt.Sprintf("Limit: %s", ui.Highlight.Sprintf("%d", vFlags.limit)))
	}

	if len(filters) > 0 {
		ui.PrintInfoPanel("ACTIVE FILTERS", filters)
		fmt.Println()
	}
}

func displayStats(stats trace.TraceStats) {
	ui.PrintTitle("QUZZ STATISTICS")
	fmt.Println()

	ui.PrintBlockStart("OVERVIEW")
	ui.PrintBlockItem(fmt.Sprintf("Total Traces:       %s", ui.Highlight.Sprint(stats.TotalTraces)))
	ui.PrintBlockItem(fmt.Sprintf("Total Duration:     %s", ui.Orange3.Sprintf("%.2fs", float64(stats.TotalDuration)/1000.0)))
	ui.PrintBlockItem(fmt.Sprintf("Average Duration:   %s", ui.Orange3.Sprintf("%.2fms", stats.AverageDuration)))
	ui.PrintBlockItem(fmt.Sprintf("Min Duration:       %s", ui.Value.Sprintf("%dms", stats.MinDuration)))
	ui.PrintBlockItem(fmt.Sprintf("Max Duration:       %s", ui.Value.Sprintf("%dms", stats.MaxDuration)))
	ui.PrintBlockItem(fmt.Sprintf("Errors:             %s", ui.Error.Sprint(stats.ErrorCount)))
	ui.PrintBlockItem(fmt.Sprintf("Warnings:           %s", ui.Warning.Sprint(stats.WarningCount)))
	ui.PrintBlockEnd()

	if len(stats.ComponentCounts) > 0 {
		fmt.Println()
		ui.PrintBlockStart("COMPONENTS")
		for component, count := range stats.ComponentCounts {
			ui.PrintBlockItem(fmt.Sprintf("%-30s %s", component, ui.Highlight.Sprintf("%d", count)))
		}
		ui.PrintBlockEnd()
	}

	if len(stats.LevelCounts) > 0 {
		fmt.Println()
		ui.PrintBlockStart("LOG LEVELS")
		for level, count := range stats.LevelCounts {
			ui.PrintBlockItem(fmt.Sprintf("%-30s %s", level, ui.Highlight.Sprintf("%d", count)))
		}
		ui.PrintBlockEnd()
	}

	fmt.Println()
}
