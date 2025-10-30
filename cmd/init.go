package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/onurhan1337/quzz-cli/internal/config"
	"github.com/onurhan1337/quzz-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	skipPrompts bool
	useTS       bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Quzz configuration for your project",
	Long: `Initialize Quzz configuration with an interactive setup wizard.

This command will guide you through creating a quzz.config.js or quzz.config.ts
file with recommended settings for your React Server Components project.

Examples:
  # Interactive setup
  quzz init

  # Skip prompts and use defaults (JavaScript)
  quzz init --skip-prompts

  # Generate TypeScript config
  quzz init --typescript`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVar(&skipPrompts, "skip-prompts", false, "Skip prompts and use default configuration")
	initCmd.Flags().BoolVar(&useTS, "typescript", false, "Generate TypeScript config file")
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println(ui.TitleStyle.Render("🚀 Quzz Configuration Setup"))
	fmt.Println()

	existingFile := checkExistingConfig()
	if existingFile != "" {
		fmt.Printf("%s Configuration file already exists: %s\n", ui.WarningStyle.Render("⚠"), existingFile)
		prompt := promptui.Prompt{
			Label:     "Overwrite existing configuration",
			IsConfirm: true,
		}
		result, err := prompt.Run()
		if err != nil || result != "y" {
			fmt.Println(ui.InfoStyle.Render("Setup cancelled."))
			return nil
		}
	}

	projectType := config.DetectProjectType()
	if projectType != "" {
		fmt.Printf("%s Detected project type: %s\n", ui.SuccessStyle.Render("✓"), projectType)
	}

	cfg := config.QuzzConfig{
		LogLevel:     "info",
		OutputFormat: "compact",
		Performance: config.PerformanceConfig{
			WarnThreshold: 500,
		},
	}

	format := config.ConfigFormatJS
	if useTS {
		format = config.ConfigFormatTS
	}

	if !skipPrompts {
		var err error
		cfg, format, err = interactiveSetup(cfg, projectType)
		if err != nil {
			return fmt.Errorf("%s setup failed: %w", ui.ErrorStyle.Render("Error:"), err)
		}
	}

	filePath, err := config.GenerateConfigFile(format, cfg)
	if err != nil {
		return fmt.Errorf("%s failed to generate config: %w", ui.ErrorStyle.Render("Error:"), err)
	}

	fmt.Println()
	fmt.Printf("%s Configuration file created: %s\n", ui.SuccessStyle.Render("✓"), ui.HighlightStyle.Render(filePath))
	fmt.Println()
	fmt.Println(ui.InfoStyle.Render("Next steps:"))
	fmt.Println("  1. Review and customize your configuration")
	fmt.Println("  2. Wrap your React Server Components with quzz()")
	fmt.Println("  3. Run your application and view traces with: quzz visualize traces.json")
	fmt.Println()
	fmt.Println(ui.DimStyle.Render("Learn more: https://github.com/onurhan1337/quzz"))

	return nil
}

func interactiveSetup(defaultCfg config.QuzzConfig, projectType string) (config.QuzzConfig, config.ConfigFormat, error) {
	cfg := defaultCfg

	formatPrompt := promptui.Select{
		Label: "Select configuration file format",
		Items: []string{"JavaScript (.js)", "TypeScript (.ts)"},
	}
	formatIdx, _, err := formatPrompt.Run()
	if err != nil {
		return cfg, config.ConfigFormatJS, err
	}

	format := config.ConfigFormatJS
	if formatIdx == 1 {
		format = config.ConfigFormatTS
	}

	logLevelPrompt := promptui.Select{
		Label: "Select log level",
		Items: []string{"debug", "info", "warn", "error"},
		CursorPos: 1,
	}
	_, logLevel, err := logLevelPrompt.Run()
	if err != nil {
		return cfg, format, err
	}
	cfg.LogLevel = logLevel

	outputFormatPrompt := promptui.Select{
		Label: "Select output format",
		Items: []string{"compact", "pretty", "json"},
	}
	_, outputFormat, err := outputFormatPrompt.Run()
	if err != nil {
		return cfg, format, err
	}
	cfg.OutputFormat = outputFormat

	thresholdPrompt := promptui.Prompt{
		Label:   "Performance warning threshold (ms)",
		Default: "500",
		Validate: func(input string) error {
			var val int
			_, err := fmt.Sscanf(input, "%d", &val)
			if err != nil || val < 0 {
				return fmt.Errorf("please enter a valid positive number")
			}
			return nil
		},
	}
	thresholdStr, err := thresholdPrompt.Run()
	if err != nil {
		return cfg, format, err
	}
	fmt.Sscanf(thresholdStr, "%d", &cfg.Performance.WarnThreshold)

	regexPrompt := promptui.Prompt{
		Label:   "Component filter regex (optional, e.g., ^(Blog|Product))",
		Default: "",
	}
	regex, err := regexPrompt.Run()
	if err != nil {
		return cfg, format, err
	}
	if regex != "" {
		cfg.ComponentFilter = regex
	}

	sensitivePrompt := promptui.Prompt{
		Label:   "Sensitive keys to redact (comma-separated, optional)",
		Default: "apiKey,secretToken,password",
	}
	sensitiveStr, err := sensitivePrompt.Run()
	if err != nil {
		return cfg, format, err
	}
	if sensitiveStr != "" {
		cfg.SensitiveKeys = parseSensitiveKeys(sensitiveStr)
	}

	return cfg, format, nil
}

func checkExistingConfig() string {
	files := []string{"quzz.config.js", "quzz.config.ts"}
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}
	return ""
}

func parseSensitiveKeys(input string) []string {
	if input == "" {
		return nil
	}

	var keys []string
	current := ""
	for _, char := range input {
		if char == ',' {
			if current != "" {
				keys = append(keys, current)
				current = ""
			}
		} else if char != ' ' {
			current += string(char)
		}
	}
	if current != "" {
		keys = append(keys, current)
	}
	return keys
}
