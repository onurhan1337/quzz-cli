package config

type QuzzConfig struct {
	LogLevel        string              `json:"logLevel"`
	OutputFormat    string              `json:"outputFormat"`
	Performance     PerformanceConfig   `json:"performance"`
	ComponentFilter string              `json:"componentFilter,omitempty"`
	SensitiveKeys   []string            `json:"sensitiveKeys,omitempty"`
}

type PerformanceConfig struct {
	WarnThreshold int `json:"warnThreshold"`
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
