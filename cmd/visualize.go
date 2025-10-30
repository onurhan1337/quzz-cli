package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/onurhan1337/quzz-cli/internal/trace"
	"github.com/onurhan1337/quzz-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
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
)

var visualizeCmd = &cobra.Command{
	Use:   "visualize [path]",
	Short: "Visualize and explore traces from traces.json",
	Long: `Visualize and explore traces with powerful filtering options.

Examples:
  # View all traces
  quzz visualize traces.json

  # Show only statistics
  quzz visualize traces.json --stats

  # Filter by component
  quzz visualize traces.json --component UserProfile

  # Filter by duration range
  quzz visualize traces.json --min-duration 100 --max-duration 500

  # Filter by date range
  quzz visualize traces.json --start-date 2024-01-01 --end-date 2024-01-31

  # Show only errors
  quzz visualize traces.json --errors

  # Component regex filtering
  quzz visualize traces.json --component-regex "^(Blog|Product)"

  # JSON output for programmatic processing
  quzz visualize traces.json --json`,
	Args: cobra.ExactArgs(1),
	RunE: runVisualize,
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	visualizeCmd.Flags().BoolVarP(&statsOnly, "stats", "s", false, "Show only statistics")
	visualizeCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")
	visualizeCmd.Flags().StringVarP(&component, "component", "c", "", "Filter by component name")
	visualizeCmd.Flags().StringVarP(&operation, "operation", "o", "", "Filter by operation")
	visualizeCmd.Flags().StringVar(&category, "category", "", "Filter by category")
	visualizeCmd.Flags().StringVarP(&level, "level", "l", "", "Filter by log level (info, warn, error, debug)")
	visualizeCmd.Flags().IntVar(&minDuration, "min-duration", 0, "Filter by minimum duration (ms)")
	visualizeCmd.Flags().IntVar(&maxDuration, "max-duration", 0, "Filter by maximum duration (ms)")
	visualizeCmd.Flags().StringVar(&startDate, "start-date", "", "Filter by start date (RFC3339 format)")
	visualizeCmd.Flags().StringVar(&endDate, "end-date", "", "Filter by end date (RFC3339 format)")
	visualizeCmd.Flags().BoolVar(&hasError, "errors", false, "Show only traces with errors")
	visualizeCmd.Flags().BoolVar(&hasWarning, "warnings", false, "Show only traces with warnings")
	visualizeCmd.Flags().StringVar(&componentRegex, "component-regex", "", "Filter components by regex pattern")
	visualizeCmd.Flags().IntVar(&limit, "limit", 50, "Limit number of traces displayed (0 for all)")
}

func runVisualize(cmd *cobra.Command, args []string) error {
	tracePath := args[0]

	traceFile, err := trace.LoadTraceFile(tracePath)
	if err != nil {
		return fmt.Errorf("%s %s", ui.ErrorStyle.Render("Error loading trace file:"), err)
	}

	filterOpts := trace.FilterOptions{
		Component:      component,
		Operation:      operation,
		Category:       category,
		Level:          level,
		MinDuration:    minDuration,
		MaxDuration:    maxDuration,
		HasError:       hasError,
		HasWarning:     hasWarning,
		ComponentRegex: componentRegex,
	}

	if startDate != "" {
		t, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return fmt.Errorf("%s invalid start date format (use RFC3339, e.g., 2024-01-01T00:00:00Z)", ui.ErrorStyle.Render("Error:"))
		}
		filterOpts.StartDate = t
	}

	if endDate != "" {
		t, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return fmt.Errorf("%s invalid end date format (use RFC3339, e.g., 2024-01-31T23:59:59Z)", ui.ErrorStyle.Render("Error:"))
		}
		filterOpts.EndDate = t
	}

	filteredTraces := trace.FilterTraces(traceFile.Traces, filterOpts)

	if len(filteredTraces) == 0 {
		fmt.Println(ui.WarningStyle.Render("⚠ No traces match the specified filters"))
		return nil
	}

	if jsonOutput {
		output := map[string]interface{}{
			"traces": filteredTraces,
			"stats":  trace.CalculateStats(filteredTraces),
		}
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("%s failed to generate JSON output: %w", ui.ErrorStyle.Render("Error:"), err)
		}
		fmt.Println(string(jsonData))
		return nil
	}

	stats := trace.CalculateStats(filteredTraces)

	if statsOnly {
		ui.RenderStatsTable(stats)
	} else {
		fmt.Println(ui.TitleStyle.Render("📋 Quzz Trace Visualization"))
		fmt.Println()
		ui.RenderStatsTable(stats)
		fmt.Println(ui.HeaderStyle.Render("\n📊 Trace Details"))
		fmt.Println()
		ui.RenderTraceTable(filteredTraces, limit)
	}

	return nil
}
