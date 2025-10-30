package cmd

import (
	"fmt"

	"github.com/onurhan1337/quzz-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version, commit hash, and build date of the Quzz CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.TitleStyle.Render("Quzz CLI"))
		fmt.Printf("%s %s\n", ui.LabelStyle.Render("Version:"), ui.ValueStyle.Render(version))
		fmt.Printf("%s %s\n", ui.LabelStyle.Render("Commit:"), ui.ValueStyle.Render(commit))
		fmt.Printf("%s %s\n", ui.LabelStyle.Render("Built:"), ui.ValueStyle.Render(date))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
