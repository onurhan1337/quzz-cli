package trace

import "time"

type TraceEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Component   string                 `json:"component"`
	Duration    int                    `json:"duration"`
	Props       map[string]interface{} `json:"props,omitempty"`
	Level       string                 `json:"level"`
	Operation   string                 `json:"operation,omitempty"`
	Category    string                 `json:"category,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Performance *PerformanceData       `json:"performance,omitempty"`
}

type PerformanceData struct {
	WarnThreshold int  `json:"warnThreshold"`
	IsWarning     bool `json:"isWarning"`
}

type TraceFile struct {
	Traces []TraceEntry `json:"traces"`
}

type TraceStats struct {
	TotalTraces      int
	TotalDuration    int64
	AverageDuration  float64
	MinDuration      int
	MaxDuration      int
	ComponentCounts  map[string]int
	OperationCounts  map[string]int
	CategoryCounts   map[string]int
	LevelCounts      map[string]int
	ErrorCount       int
	WarningCount     int
	DateRange        DateRange
}

type DateRange struct {
	Earliest time.Time
	Latest   time.Time
}

type FilterOptions struct {
	Component    string
	Operation    string
	Category     string
	Level        string
	MinDuration  int
	MaxDuration  int
	StartDate    time.Time
	EndDate      time.Time
	HasError     bool
	HasWarning   bool
	ComponentRegex string
}
