package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateConfigFile(format ConfigFormat, config QuzzConfig) (string, error) {
	var content string
	filename := "quzz.config" + format.Extension()

	if format == ConfigFormatTS {
		content = generateTSConfig(config)
	} else {
		content = generateJSConfig(config)
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

func generateJSConfig(config QuzzConfig) string {
	sensitiveKeys := ""
	if len(config.SensitiveKeys) > 0 {
		sensitiveKeys = fmt.Sprintf(`  sensitiveKeys: %v,`, formatStringArray(config.SensitiveKeys))
	}

	componentFilter := ""
	if config.ComponentFilter != "" {
		componentFilter = fmt.Sprintf("  componentFilter: /%s/,\n", config.ComponentFilter)
	}

	return fmt.Sprintf(`/**
 * Quzz Configuration File
 *
 * This file configures the Quzz debugging tool for your React Server Components.
 * Learn more: https://github.com/onurhan1337/quzz
 */

module.exports = {
  // Log level: "debug" | "info" | "warn" | "error"
  logLevel: "%s",

  // Output format: "pretty" | "compact" | "json"
  outputFormat: "%s",

  // Performance monitoring configuration
  performance: {
    // Warn threshold in milliseconds
    warnThreshold: %d,
  },
%s%s
};
`, config.LogLevel, config.OutputFormat, config.Performance.WarnThreshold, componentFilter, sensitiveKeys)
}

func generateTSConfig(config QuzzConfig) string {
	sensitiveKeys := ""
	if len(config.SensitiveKeys) > 0 {
		sensitiveKeys = fmt.Sprintf(`  sensitiveKeys: %v,`, formatStringArray(config.SensitiveKeys))
	}

	componentFilter := ""
	if config.ComponentFilter != "" {
		componentFilter = fmt.Sprintf("  componentFilter: /%s/,\n", config.ComponentFilter)
	}

	return fmt.Sprintf(`/**
 * Quzz Configuration File
 *
 * This file configures the Quzz debugging tool for your React Server Components.
 * Learn more: https://github.com/onurhan1337/quzz
 */

import type { QuzzConfig } from 'quzz';

const config: QuzzConfig = {
  // Log level: "debug" | "info" | "warn" | "error"
  logLevel: "%s",

  // Output format: "pretty" | "compact" | "json"
  outputFormat: "%s",

  // Performance monitoring configuration
  performance: {
    // Warn threshold in milliseconds
    warnThreshold: %d,
  },
%s%s
};

export default config;
`, config.LogLevel, config.OutputFormat, config.Performance.WarnThreshold, componentFilter, sensitiveKeys)
}

func formatStringArray(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	result := "["
	for i, s := range arr {
		result += fmt.Sprintf(`"%s"`, s)
		if i < len(arr)-1 {
			result += ", "
		}
	}
	result += "]"
	return result
}

func DetectProjectType() string {
	files := []string{"package.json", "tsconfig.json", "next.config.js", "next.config.ts"}

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}

	return ""
}
