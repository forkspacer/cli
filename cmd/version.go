package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/forkspacer/cli/pkg/styles"
)

var (
	// Set via ldflags during build
	version   = "dev"
	gitCommit = "none"
	buildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Println(styles.TitleStyle.Render("Forkspacer CLI"))
	fmt.Println()
	fmt.Printf("%s  %s\n", styles.Key("Version:"), styles.Value(version))
	fmt.Printf("%s  %s\n", styles.Key("Git Commit:"), styles.Value(gitCommit))
	fmt.Printf("%s  %s\n", styles.Key("Build Date:"), styles.Value(buildDate))
	fmt.Printf("%s  %s\n", styles.Key("Go Version:"), styles.Value(runtime.Version()))
	fmt.Printf("%s  %s/%s\n", styles.Key("Platform:"), styles.Value(runtime.GOOS), styles.Value(runtime.GOARCH))
}
