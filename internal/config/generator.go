package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type GeneratorOptions struct {
	Format  ConfigFormat
	Config  QuzzConfig
	Minimal bool
}

// templateData holds all data needed for config file generation
type templateData struct {
	IsTypeScript bool
	Config       QuzzConfig
	Minimal      bool
}

// configTemplate is a unified template for both JS and TS config files
const configTemplate = `{{- if .IsTypeScript -}}
// quzz.config.ts

import type { QuzzConfig } from 'quzz';

const config: QuzzConfig = {
{{- else -}}
// quzz.config.js

/** @type {import('quzz').QuzzConfig} */
module.exports = {
{{- end }}
  logLevel: "{{.Config.LogLevel}}",
  outputFormat: "{{.Config.OutputFormat}}",

  performance: {
    enabled: {{.Config.Performance.Enabled}},
    warnThreshold: {{.Config.Performance.WarnThreshold}},
{{- if not .Minimal }}
    trackMemory: {{.Config.Performance.TrackMemory}},
{{- if gt .Config.Performance.MemoryThreshold 0 }}
    memoryThreshold: {{.Config.Performance.MemoryThreshold}},
{{- end }}
    aggregate: {{.Config.Performance.Aggregate}},
    enableHeapSnapshots: {{.Config.Performance.EnableHeapSnapshots}},
{{- if ne .Config.Performance.HeapSnapshotDir "" }}
    heapSnapshotDir: "{{.Config.Performance.HeapSnapshotDir}}",
{{- end }}
{{- end }}
  },
{{- if not .Minimal }}

  props: {
    awaitProps: {{.Config.Props.AwaitProps}},
    awaitTimeout: {{.Config.Props.AwaitTimeout}},
    showPromiseTypes: {{.Config.Props.ShowPromiseTypes}},
    maxArrayItems: {{.Config.Props.MaxArrayItems}},
    maxObjectProps: {{.Config.Props.MaxObjectProps}},
    maxErrorDepth: {{.Config.Props.MaxErrorDepth}},
    serializationStrategy: "{{.Config.Props.SerializationStrategy}}",
  },

  forceEnable: {{.Config.ForceEnable}},
  maxPropDepth: {{.Config.MaxPropDepth}},
  maxStringLength: {{.Config.MaxStringLength}},
{{- if gt (len .Config.SensitiveKeys) 0 }}

  sensitiveKeys: [{{range $i, $key := .Config.SensitiveKeys}}{{if $i}}, {{end}}"{{$key}}"{{end}}],
{{- end }}

  contextTracking: {{.Config.ContextTracking}},
  includeSourceLocation: {{.Config.IncludeSourceLocation}},
{{- if ne .Config.ComponentFilter "" }}
  componentFilter: /{{.Config.ComponentFilter}}/,
{{- end }}
{{- if gt .Config.ThrottleMs 0 }}
  throttleMs: {{.Config.ThrottleMs}},
{{- end }}

  trackTotalLatency: {{.Config.TrackTotalLatency}},
{{- end }}

  visualizer: {
    enabled: {{.Config.Visualizer.Enabled}},
    output: "{{.Config.Visualizer.Output}}",
  },
{{- if not .Minimal }}

  debugContext: {{.Config.DebugContext}},
  enableSnapshots: {{.Config.EnableSnapshots}},
  verboseMode: {{.Config.VerboseMode}},
  autoLinkParent: {{.Config.AutoLinkParent}},
  suppressConfigWarnings: {{.Config.SuppressConfigWarnings}},
  enableHyperlinks: {{.Config.EnableHyperlinks}},
{{- end }}
{{- if .IsTypeScript }}
};

export default config;
{{- else }}
};
{{- end }}
`

var (
	// tmpl is the compiled template, initialized once
	tmpl = template.Must(template.New("config").Parse(configTemplate))
)

// GenerateConfigFile creates a configuration file with the specified options
func GenerateConfigFile(opts GeneratorOptions) (string, error) {
	filename := "quzz.config" + opts.Format.Extension()

	content, err := generateConfig(opts.Format, opts.Config, opts.Minimal)
	if err != nil {
		return "", fmt.Errorf("failed to generate config: %w", err)
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write config file: %w", err)
	}

	return absPath, nil
}

// generateConfig generates configuration file content using a template
func generateConfig(format ConfigFormat, cfg QuzzConfig, minimal bool) (string, error) {
	data := templateData{
		IsTypeScript: format == ConfigFormatTS,
		Config:       cfg,
		Minimal:      minimal,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// formatStringArray formats a string slice as a JavaScript/TypeScript array
// Deprecated: This is now handled by the template
func formatStringArray(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}

	parts := make([]string, len(arr))
	for i, s := range arr {
		parts[i] = fmt.Sprintf("\"%s\"", s)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// DetectProjectType checks for common project files to help determine project type
func DetectProjectType() string {
	files := []string{"package.json", "tsconfig.json", "next.config.js", "next.config.ts"}

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}

	return ""
}
