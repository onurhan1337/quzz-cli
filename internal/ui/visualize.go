package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/onurhan1337/quzz-cli/internal/trace"
	"github.com/onurhan1337/quzz-cli/internal/ui/format"
)

func RenderTraceVisualization(traces []trace.TraceEntry, stats trace.TraceStats, limit int) {
	PrintTitle("QUZZ TRACE VISUALIZATION")

	renderStatsBlock(stats)
	fmt.Println()
	renderTracesBlock(traces, limit)
}

func renderStatsBlock(stats trace.TraceStats) {
	PrintBlockStart("STATISTICS")

	PrintBlockSection("Overview")
	PrintBlockLine(fmt.Sprintf("Total Traces      %s", Highlight.Sprintf("%d", stats.TotalTraces)))
	PrintBlockLine(fmt.Sprintf("Total Duration    %s", Highlight.Sprint(format.Milliseconds(stats.TotalDuration))))
	PrintBlockLine(fmt.Sprintf("Average Duration  %s", format.Milliseconds(int64(stats.AverageDuration))))
	PrintBlockLine(fmt.Sprintf("Min Duration      %dms", stats.MinDuration))
	PrintBlockLine(fmt.Sprintf("Max Duration      %dms", stats.MaxDuration))

	if stats.ErrorCount > 0 || stats.WarningCount > 0 {
		PrintBlockSection("Issues")
		if stats.ErrorCount > 0 {
			PrintBlockLine(fmt.Sprintf("Errors    %s", Error.Sprintf("%d", stats.ErrorCount)))
		}
		if stats.WarningCount > 0 {
			PrintBlockLine(fmt.Sprintf("Warnings  %s", Warning.Sprintf("%d", stats.WarningCount)))
		}
	}

	if !stats.DateRange.Earliest.IsZero() {
		PrintBlockSection("Time Range")
		PrintBlockLine(fmt.Sprintf("%s %s %s",
			Dim.Sprint(stats.DateRange.Earliest.Format(time.RFC3339)),
			IconArrow,
			Dim.Sprint(stats.DateRange.Latest.Format(time.RFC3339))))
	}

	if len(stats.ComponentCounts) > 0 {
		PrintBlockSection("Components")
		renderTopItems(stats.ComponentCounts, 10)
	}

	if len(stats.LevelCounts) > 0 {
		PrintBlockSection("Log Levels")
		renderLevelItems(stats.LevelCounts)
	}

	if len(stats.OperationCounts) > 0 {
		PrintBlockSection("Operations")
		renderTopItems(stats.OperationCounts, 10)
	}

	PrintBlockEnd()
}

func renderTracesBlock(traces []trace.TraceEntry, limit int) {
	displayCount := calculateDisplayCount(traces, limit)

	PrintBlockStart("TRACES")

	Primary.Printf("%s ", vertical)
	fmt.Printf("%-10s %-20s %-12s %-8s %-12s %-12s",
		Bold.Sprint("TIME"),
		Bold.Sprint("COMPONENT"),
		Bold.Sprint("DURATION"),
		Bold.Sprint("LEVEL"),
		Bold.Sprint("OPERATION"),
		Bold.Sprint("STATUS"))
	Primary.Printf(" %s\n", vertical)

	Primary.Printf("%s ", vertical)
	Dim.Print(strings.Repeat(horizontal, boxWidth-4))
	Primary.Printf(" %s\n", vertical)

	for i := 0; i < displayCount; i++ {
		renderTraceRow(traces[i])
	}

	PrintBlockEnd()

	if shouldShowLimitMessage(traces, limit) {
		fmt.Println()
		Dim.Printf("Showing %d of %d traces. Use --limit 0 to show all traces.\n",
			displayCount, len(traces))
	}
}

func renderTraceRow(t trace.TraceEntry) {
	timeStr := t.Timestamp.Format("15:04:05")
	componentStr := format.Component(t.Component, 20)
	durationStr := format.Duration(t.Duration)
	levelStr := format.Level(t.Level)
	operationStr := format.Operation(t.Operation)
	statusStr := format.Status(t)

	Primary.Printf("%s ", vertical)
	fmt.Printf("%-10s %-20s %-12s %-8s %-12s %-12s",
		timeStr,
		componentStr,
		durationStr,
		levelStr,
		operationStr,
		statusStr)
	Primary.Printf(" %s\n", vertical)
}

func renderTopItems(counts map[string]int, maxItems int) {
	type countPair struct {
		name  string
		count int
	}

	pairs := make([]countPair, 0, len(counts))
	for name, count := range counts {
		pairs = append(pairs, countPair{name, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	displayCount := len(pairs)
	if maxItems > 0 && maxItems < len(pairs) {
		displayCount = maxItems
	}

	for i := 0; i < displayCount; i++ {
		PrintBlockLine(fmt.Sprintf("%s %-25s %s",
			IconBullet,
			pairs[i].name,
			Dim.Sprintf("(%d)", pairs[i].count)))
	}

	if len(pairs) > displayCount {
		PrintBlockLine(Dim.Sprintf("... and %d more", len(pairs)-displayCount))
	}
}

func renderLevelItems(counts map[string]int) {
	levels := []string{"error", "warn", "info", "debug"}

	for _, level := range levels {
		if count, ok := counts[level]; ok {
			levelStr := format.LevelColored(level)
			PrintBlockLine(fmt.Sprintf("%s %s", levelStr, Dim.Sprintf("(%d)", count)))
		}
	}
}

func calculateDisplayCount(traces []trace.TraceEntry, limit int) int {
	if limit <= 0 || limit > len(traces) {
		return len(traces)
	}
	return limit
}

func shouldShowLimitMessage(traces []trace.TraceEntry, limit int) bool {
	return limit > 0 && len(traces) > limit
}
