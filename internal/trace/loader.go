package trace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func LoadTraceFile(path string) (*TraceFile, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read trace file: %w", err)
	}

	var traceFile TraceFile
	if err := json.Unmarshal(data, &traceFile); err != nil {
		return nil, fmt.Errorf("failed to parse trace file: %w", err)
	}

	return &traceFile, nil
}

func FilterTraces(traces []TraceEntry, opts FilterOptions) []TraceEntry {
	filtered := make([]TraceEntry, 0)

	var componentRegex *regexp.Regexp
	if opts.ComponentRegex != "" {
		componentRegex = regexp.MustCompile(opts.ComponentRegex)
	}

	for _, trace := range traces {
		if opts.Component != "" && trace.Component != opts.Component {
			continue
		}

		if opts.Operation != "" && trace.Operation != opts.Operation {
			continue
		}

		if opts.Category != "" && trace.Category != opts.Category {
			continue
		}

		if opts.Level != "" && trace.Level != opts.Level {
			continue
		}

		if opts.MinDuration > 0 && trace.Duration < opts.MinDuration {
			continue
		}

		if opts.MaxDuration > 0 && trace.Duration > opts.MaxDuration {
			continue
		}

		if !opts.StartDate.IsZero() && trace.Timestamp.Before(opts.StartDate) {
			continue
		}

		if !opts.EndDate.IsZero() && trace.Timestamp.After(opts.EndDate) {
			continue
		}

		if opts.HasError && trace.Error == "" {
			continue
		}

		if opts.HasWarning && (trace.Performance == nil || !trace.Performance.IsWarning) {
			continue
		}

		if componentRegex != nil && !componentRegex.MatchString(trace.Component) {
			continue
		}

		filtered = append(filtered, trace)
	}

	return filtered
}

func CalculateStats(traces []TraceEntry) TraceStats {
	stats := TraceStats{
		ComponentCounts: make(map[string]int),
		OperationCounts: make(map[string]int),
		CategoryCounts:  make(map[string]int),
		LevelCounts:     make(map[string]int),
	}

	if len(traces) == 0 {
		return stats
	}

	stats.TotalTraces = len(traces)
	stats.MinDuration = traces[0].Duration
	stats.MaxDuration = traces[0].Duration
	stats.DateRange.Earliest = traces[0].Timestamp
	stats.DateRange.Latest = traces[0].Timestamp

	var totalDuration int64
	for _, trace := range traces {
		totalDuration += int64(trace.Duration)

		if trace.Duration < stats.MinDuration {
			stats.MinDuration = trace.Duration
		}
		if trace.Duration > stats.MaxDuration {
			stats.MaxDuration = trace.Duration
		}

		stats.ComponentCounts[trace.Component]++
		if trace.Operation != "" {
			stats.OperationCounts[trace.Operation]++
		}
		if trace.Category != "" {
			stats.CategoryCounts[trace.Category]++
		}
		if trace.Level != "" {
			stats.LevelCounts[trace.Level]++
		}

		if trace.Error != "" {
			stats.ErrorCount++
		}

		if trace.Performance != nil && trace.Performance.IsWarning {
			stats.WarningCount++
		}

		if trace.Timestamp.Before(stats.DateRange.Earliest) {
			stats.DateRange.Earliest = trace.Timestamp
		}
		if trace.Timestamp.After(stats.DateRange.Latest) {
			stats.DateRange.Latest = trace.Timestamp
		}
	}

	stats.TotalDuration = totalDuration
	stats.AverageDuration = float64(totalDuration) / float64(len(traces))

	return stats
}
