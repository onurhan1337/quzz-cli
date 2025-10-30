package config

const (
    DefaultLogLevel            = "error"
    DefaultOutputFormat        = "pretty"
    DefaultWarnThreshold       = 1000
    DefaultTrackMemory         = false
    DefaultMemoryThreshold     = 50 * 1024 * 1024 // 50MB
    DefaultAggregate           = false
    DefaultEnableHeapSnapshots = false
    DefaultHeapSnapshotDir     = "./heap-snapshots"
    DefaultAwaitProps          = false
    DefaultAwaitTimeout        = 5000
    DefaultShowPromiseTypes    = true
    DefaultMaxArrayItems       = 10
    DefaultMaxObjectProps      = 20
    DefaultMaxErrorDepth       = 3
    DefaultSerializationStrategy = "standard"
    DefaultLogProps            = false
    DefaultForceEnable         = false
    DefaultMaxPropDepth        = 3
    DefaultMaxStringLength     = 200
    DefaultContextTracking     = true
    DefaultIncludeSourceLocation = false
    DefaultThrottleMs          = 0
    DefaultTrackTotalLatency   = false
    DefaultAutoLinkParent      = true
    DefaultVisualizerEnabled   = false
    DefaultVisualizerOutput    = "./traces.json"
    DefaultDebugContext        = false
    DefaultEnableSnapshots     = false
    DefaultVerboseMode         = false
    DefaultSuppressConfigWarnings = false
    DefaultEnableHyperlinks    = true
)

type QuzzConfig struct {
    LogLevel               string           `json:"logLevel"`
    OutputFormat           string           `json:"outputFormat"`
    Performance            PerformanceConfig  `json:"performance"`
    Props                  PropsConfig      `json:"props"`
    LogProps               bool             `json:"logProps,omitempty"`
    ForceEnable            bool             `json:"forceEnable,omitempty"`
    MaxPropDepth           int              `json:"maxPropDepth,omitempty"`
    MaxStringLength        int              `json:"maxStringLength,omitempty"`
    SensitiveKeys          []string         `json:"sensitiveKeys,omitempty"`
    ContextTracking        bool             `json:"contextTracking,omitempty"`
    IncludeSourceLocation  bool             `json:"includeSourceLocation,omitempty"`
    ComponentFilter        string           `json:"componentFilter,omitempty"`
    ThrottleMs             int              `json:"throttleMs,omitempty"`
    TrackTotalLatency      bool             `json:"trackTotalLatency,omitempty"`
    Visualizer             VisualizerConfig `json:"visualizer"`
    DebugContext           bool             `json:"debugContext,omitempty"`
    EnableSnapshots        bool             `json:"enableSnapshots,omitempty"`
    VerboseMode            bool             `json:"verboseMode,omitempty"`
    AutoLinkParent         bool             `json:"autoLinkParent,omitempty"`
    SuppressConfigWarnings bool             `json:"suppressConfigWarnings,omitempty"`
    EnableHyperlinks       bool             `json:"enableHyperlinks,omitempty"`
}

type PerformanceConfig struct {
    Enabled              bool   `json:"enabled"`
    WarnThreshold        int    `json:"warnThreshold"`
    TrackMemory          bool   `json:"trackMemory,omitempty"`
    MemoryThreshold      int    `json:"memoryThreshold,omitempty"`
    Aggregate            bool   `json:"aggregate,omitempty"`
    EnableHeapSnapshots  bool   `json:"enableHeapSnapshots,omitempty"`
    HeapSnapshotDir      string `json:"heapSnapshotDir,omitempty"`
}

type PropsConfig struct {
    AwaitProps            bool   `json:"awaitProps"`
    AwaitTimeout          int    `json:"awaitTimeout,omitempty"`
    ShowPromiseTypes      bool   `json:"showPromiseTypes,omitempty"`
    MaxArrayItems         int    `json:"maxArrayItems,omitempty"`
    MaxObjectProps        int    `json:"maxObjectProps,omitempty"`
    MaxErrorDepth         int    `json:"maxErrorDepth,omitempty"`
    SerializationStrategy string `json:"serializationStrategy,omitempty"`
}

type VisualizerConfig struct {
    Enabled bool   `json:"enabled"`
    Output  string `json:"output"`
}

type ConfigFormat int

const (
	ConfigFormatJS ConfigFormat = iota
	ConfigFormatTS
)

func (f ConfigFormat) String() string {
	switch f {
	case ConfigFormatJS:
		return "JavaScript"
	case ConfigFormatTS:
		return "TypeScript"
	default:
		return "JavaScript"
	}
}

func (f ConfigFormat) Extension() string {
	switch f {
	case ConfigFormatJS:
		return ".js"
	case ConfigFormatTS:
		return ".ts"
	default:
		return ".js"
	}
}

func DefaultConfig() QuzzConfig {
    return QuzzConfig{
        LogLevel:     DefaultLogLevel,
        OutputFormat: DefaultOutputFormat,
        Performance: PerformanceConfig{
            Enabled:             false,
            WarnThreshold:       DefaultWarnThreshold,
            TrackMemory:         DefaultTrackMemory,
            MemoryThreshold:     DefaultMemoryThreshold,
            Aggregate:           DefaultAggregate,
            EnableHeapSnapshots: DefaultEnableHeapSnapshots,
            HeapSnapshotDir:     DefaultHeapSnapshotDir,
        },
        Props: PropsConfig{
            AwaitProps:            DefaultAwaitProps,
            AwaitTimeout:          DefaultAwaitTimeout,
            ShowPromiseTypes:      DefaultShowPromiseTypes,
            MaxArrayItems:         DefaultMaxArrayItems,
            MaxObjectProps:        DefaultMaxObjectProps,
            MaxErrorDepth:         DefaultMaxErrorDepth,
            SerializationStrategy: DefaultSerializationStrategy,
        },
        LogProps:               DefaultLogProps,
        ForceEnable:            DefaultForceEnable,
        MaxPropDepth:           DefaultMaxPropDepth,
        MaxStringLength:        DefaultMaxStringLength,
        SensitiveKeys:          []string{}, // empty by default, user-extendable
        ContextTracking:        DefaultContextTracking,
        IncludeSourceLocation:  DefaultIncludeSourceLocation,
        ComponentFilter:        "",
        ThrottleMs:             DefaultThrottleMs,
        TrackTotalLatency:      DefaultTrackTotalLatency,
        Visualizer: VisualizerConfig{
            Enabled: DefaultVisualizerEnabled,
            Output:  DefaultVisualizerOutput,
        },
        DebugContext:           DefaultDebugContext,
        EnableSnapshots:        DefaultEnableSnapshots,
        VerboseMode:            DefaultVerboseMode,
        AutoLinkParent:         DefaultAutoLinkParent,
        SuppressConfigWarnings: DefaultSuppressConfigWarnings,
        EnableHyperlinks:       DefaultEnableHyperlinks,
    }
}

func MinimalConfig() QuzzConfig {
	return QuzzConfig{
		LogLevel:     DefaultLogLevel,
		OutputFormat: DefaultOutputFormat,
		Performance: PerformanceConfig{
			Enabled:       true,
			WarnThreshold: DefaultWarnThreshold,
		},
		Props: PropsConfig{
			AwaitProps: true,
		},
		Visualizer: VisualizerConfig{
			Enabled: true,
			Output:  DefaultVisualizerOutput,
		},
	}
}
