package cmd

import (
	"fmt"

	"github.com/onurhan1337/quzz-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version = "v0.5.6"
	author  = "Onurhan Demir"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version, author, and build date of the Quzz CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintTitle("QUZZ CLI")
		fmt.Println()
		ui.PrintSectionItem("Version", ui.Highlight.Sprint(version))
		ui.PrintSectionItem("Author", ui.Value.Sprint(author))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
