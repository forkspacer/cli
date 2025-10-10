package cmd

import (
	"fmt"
	"os"

	"github.com/forkspacer/cli/pkg/styles"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	namespace string
	output    string
	verbose   bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "forkspacer",
	Short: "Manage Forkspacer workspaces and modules",
	Long: styles.TitleStyle.Render("Forkspacer CLI") + "\n\n" +
		"A cloud-native Kubernetes operator for dynamic workspace lifecycle management.\n" +
		"Create, manage, and hibernate ephemeral development environments at scale.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default",
		"Kubernetes namespace")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table",
		"Output format (table|json|yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Enable verbose output")

	// Custom help template with better styling
	rootCmd.SetHelpTemplate(getHelpTemplate())
}

// GetRootCmd returns the root command for use in subcommands
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func getHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

// HandleError provides consistent error formatting
func HandleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "\n"+styles.Error(err.Error()))
		os.Exit(1)
	}
}

// GetNamespace returns the configured namespace
func GetNamespace() string {
	return namespace
}

// GetOutput returns the configured output format
func GetOutput() string {
	return output
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return verbose
}
