package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/onurhan1337/quzz-cli/internal/trace"
)

func RenderTraceTable(traces []trace.TraceEntry, limit int) {
	headerStyle := HeaderStyle.Render

	fmt.Printf("%-12s %-25s %-15s %-10s %-15s %-15s\n",
		headerStyle("Time"),
		headerStyle("Component"),
		headerStyle("Duration"),
		headerStyle("Level"),
		headerStyle("Operation"),
		headerStyle("Status"))
	fmt.Println(strings.Repeat("─", 100))

	displayCount := len(traces)
	if limit > 0 && limit < len(traces) {
		displayCount = limit
	}

	for i := 0; i < displayCount; i++ {
		t := traces[i]
		status := "✓"
		statusStyle := SuccessStyle.Render
		if t.Error != "" {
			status = "✗ Error"
			statusStyle = ErrorStyle.Render
		} else if t.Performance != nil && t.Performance.IsWarning {
			status = "⚠ Warning"
			statusStyle = WarningStyle.Render
		}

		operation := t.Operation
		if operation == "" {
			operation = "-"
		}

		levelStyle := InfoStyle.Render
		switch strings.ToLower(t.Level) {
		case "error":
			levelStyle = ErrorStyle.Render
		case "warn":
			levelStyle = WarningStyle.Render
		case "debug":
			levelStyle = DimStyle.Render
		}

		fmt.Printf("%-12s %-25s %-15s %-10s %-15s %-15s\n",
			t.Timestamp.Format("15:04:05"),
			t.Component,
			fmt.Sprintf("%dms", t.Duration),
			levelStyle(strings.ToUpper(t.Level)),
			operation,
			statusStyle(status))
	}

	if limit > 0 && len(traces) > limit {
		fmt.Printf("\n%s\n", DimStyle.Render(fmt.Sprintf("Showing %d of %d traces. Use --limit to see more.", limit, len(traces))))
	}
}

func RenderStatsTable(stats trace.TraceStats) {
	fmt.Println(TitleStyle.Render("📊 Trace Statistics"))

	fmt.Printf("%s %s\n", LabelStyle.Render("Total Traces:"), ValueStyle.Render(fmt.Sprintf("%d", stats.TotalTraces)))
	fmt.Printf("%s %s\n", LabelStyle.Render("Total Duration:"), ValueStyle.Render(fmt.Sprintf("%dms", stats.TotalDuration)))
	fmt.Printf("%s %s\n", LabelStyle.Render("Average Duration:"), ValueStyle.Render(fmt.Sprintf("%.2fms", stats.AverageDuration)))
	fmt.Printf("%s %s\n", LabelStyle.Render("Min Duration:"), ValueStyle.Render(fmt.Sprintf("%dms", stats.MinDuration)))
	fmt.Printf("%s %s\n", LabelStyle.Render("Max Duration:"), ValueStyle.Render(fmt.Sprintf("%dms", stats.MaxDuration)))

	if stats.ErrorCount > 0 {
		fmt.Printf("%s %s\n", LabelStyle.Render("Errors:"), ErrorStyle.Render(fmt.Sprintf("%d", stats.ErrorCount)))
	}
	if stats.WarningCount > 0 {
		fmt.Printf("%s %s\n", LabelStyle.Render("Warnings:"), WarningStyle.Render(fmt.Sprintf("%d", stats.WarningCount)))
	}

	if !stats.DateRange.Earliest.IsZero() {
		fmt.Printf("%s %s - %s\n",
			LabelStyle.Render("Date Range:"),
			ValueStyle.Render(stats.DateRange.Earliest.Format(time.RFC3339)),
			ValueStyle.Render(stats.DateRange.Latest.Format(time.RFC3339)))
	}

	if len(stats.ComponentCounts) > 0 {
		fmt.Println(HeaderStyle.Render("\n🔧 Components"))
		for component, count := range stats.ComponentCounts {
			fmt.Printf("  %s: %s\n", ValueStyle.Render(component), HighlightStyle.Render(fmt.Sprintf("%d", count)))
		}
	}

	if len(stats.LevelCounts) > 0 {
		fmt.Println(HeaderStyle.Render("\n📈 Log Levels"))
		for level, count := range stats.LevelCounts {
			fmt.Printf("  %s: %s\n", ValueStyle.Render(level), HighlightStyle.Render(fmt.Sprintf("%d", count)))
		}
	}

	if len(stats.OperationCounts) > 0 {
		fmt.Println(HeaderStyle.Render("\n⚙️  Operations"))
		for op, count := range stats.OperationCounts {
			fmt.Printf("  %s: %s\n", ValueStyle.Render(op), HighlightStyle.Render(fmt.Sprintf("%d", count)))
		}
	}

	if len(stats.CategoryCounts) > 0 {
		fmt.Println(HeaderStyle.Render("\n📁 Categories"))
		for cat, count := range stats.CategoryCounts {
			fmt.Printf("  %s: %s\n", ValueStyle.Render(cat), HighlightStyle.Render(fmt.Sprintf("%d", count)))
		}
	}
}
