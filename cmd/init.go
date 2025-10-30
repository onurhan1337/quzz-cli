package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/onurhan1337/quzz-cli/internal/config"
	"github.com/onurhan1337/quzz-cli/internal/ui"
	"github.com/onurhan1337/quzz-cli/internal/ui/wizard"
	"github.com/spf13/cobra"
)

type initFlags struct {
	skipPrompts bool
	useTS       bool
	minimal     bool
}

var iFlags initFlags

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Quzz configuration for your project",
	Long: `Initialize Quzz configuration with an interactive setup wizard.

Creates a quzz.config.js or quzz.config.ts file with configuration options.
Supports both minimal (zero-config) and full configuration modes.

Examples:
  quzz init
  quzz init --typescript
  quzz init --minimal
  quzz init --skip-prompts --typescript`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVar(&iFlags.skipPrompts, "skip-prompts", false, "Skip prompts, use defaults")
	initCmd.Flags().BoolVar(&iFlags.useTS, "typescript", false, "Generate TypeScript config")
	initCmd.Flags().BoolVar(&iFlags.minimal, "minimal", false, "Generate minimal config (zero-config)")
}

func runInit(cmd *cobra.Command, args []string) error {
	if err := checkExistingConfig(); err != nil {
		return err
	}

	detectAndShowProjectType()

	var cfg config.QuzzConfig
	var format config.ConfigFormat
	var minimal bool
	var filePath string

	if iFlags.skipPrompts {
		cfg, format, minimal = getDefaultOptions()
		var err error
		filePath, err = generateConfig(cfg, format, minimal)
		if err != nil {
			return err
		}
	} else {
		m := wizard.New()
		p := tea.NewProgram(m)
		result, err := p.Run()
		if err != nil {
			return err
		}

		if initResult := wizard.ExtractResult(result); initResult != nil {
			cfg = initResult.Config
			format = initResult.Format
			minimal = initResult.Minimal
			filePath = initResult.FilePath
		} else {
			ui.PrintInfoPanel("SETUP CANCELLED", []string{
				"Configuration setup was cancelled by user",
			})
			return fmt.Errorf("setup cancelled by user")
		}
	}

	printSuccess(filePath, minimal, cfg)
	return nil
}

func checkExistingConfig() error {
	existingFile := findExistingConfig()
	if existingFile == "" {
		return nil
	}

	ui.PrintWarningPanel("WARNING", []string{
		fmt.Sprintf("Configuration file already exists: %s", ui.Highlight.Sprint(existingFile)),
		"",
		"File will be overwritten if you continue.",
	})

	return nil
}

func findExistingConfig() string {
	files := []string{"quzz.config.js", "quzz.config.ts"}
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}
	return ""
}

func detectAndShowProjectType() {
	projectType := config.DetectProjectType()
	if projectType != "" {
		ui.PrintInfoPanel("INFO", []string{
			fmt.Sprintf("Detected project file: %s", ui.Highlight.Sprint(projectType)),
		})
	}
}

func getDefaultOptions() (config.QuzzConfig, config.ConfigFormat, bool) {
	format := config.ConfigFormatJS
	if iFlags.useTS {
		format = config.ConfigFormatTS
	}

	minimal := iFlags.minimal
	cfg := config.DefaultConfig()
	if minimal {
		cfg = config.MinimalConfig()
	}

	return cfg, format, minimal
}

func generateConfig(cfg config.QuzzConfig, format config.ConfigFormat, minimal bool) (string, error) {
	opts := config.GeneratorOptions{
		Format:  format,
		Config:  cfg,
		Minimal: minimal,
	}

	filePath, err := config.GenerateConfigFile(opts)
	if err != nil {
		ui.PrintWarningPanel("GENERATION ERROR", []string{
			fmt.Sprintf("Failed to generate config file: %v", err),
		})
		return "", err
	}

	return filePath, nil
}

func printSuccess(filePath string, minimal bool, cfg config.QuzzConfig) {
	mode := "Full"
	if minimal {
		mode = "Minimal"
	}

	ui.PrintSuccessPanel("CONFIGURATION CREATED", []string{
		fmt.Sprintf("File:   %s", ui.Highlight.Sprint(filePath)),
		fmt.Sprintf("Mode:   %s", ui.Highlight.Sprint(mode)),
		fmt.Sprintf("Format: %s", ui.Highlight.Sprint(cfg.OutputFormat)),
		fmt.Sprintf("Level:  %s", ui.Highlight.Sprint(cfg.LogLevel)),
	})

	if minimal {
		ui.PrintInfoPanel("INFO", []string{
			"Zero-config mode enabled",
			"Quzz will work out of the box with sensible defaults",
			"You can customize the config file later as needed",
		})
	}

	ui.PrintSection("NEXT STEPS")
	fmt.Println()
	ui.PrintListItem("1", "Review and customize your configuration")
	ui.PrintListItem("2", "Wrap your React Server Components with quzz()")
	ui.PrintListItem("3", "Run your application in development mode")
	ui.PrintListItem("4", "Visualize traces: "+ui.Highlight.Sprint("quzz visualize traces.json"))
	fmt.Println()

	ui.Dim.Printf("  Learn more: https://github.com/onurhan1337/quzz\n\n")
}
