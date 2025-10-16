package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/forkspacer/cli/pkg/styles"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	SilenceUsage: true,
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

	// Register namespace flag completion
	rootCmd.RegisterFlagCompletionFunc("namespace", namespaceCompletionFunc)

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

// namespaceCompletionFunc provides dynamic completion for namespace flag
func namespaceCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Create k8s client
	restConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	k8sClient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// List all namespaces
	ctx := context.Background()
	namespaces := &corev1.NamespaceList{}
	if err := k8sClient.List(ctx, namespaces); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Extract namespace names
	var names []string
	for _, ns := range namespaces.Items {
		names = append(names, ns.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
