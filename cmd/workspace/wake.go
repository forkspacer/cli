package workspace

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/k8s"
	"github.com/forkspacer/cli/pkg/printer"
	"github.com/forkspacer/cli/pkg/styles"
)

var wakeCmd = &cobra.Command{
	Use:   "wake [name]",
	Short: "Wake up a hibernated workspace",
	Long: `Wake up a hibernated workspace to restore all modules.

This will:
  • Restore all modules to their previous state
  • Resume all services
  • Make the workspace active again

Examples:
  # Wake up a workspace
  forkspacer workspace wake dev-env

  # Wake workspace in specific namespace
  forkspacer workspace wake staging -n production`,
	Args: cobra.ExactArgs(1),
	RunE: runWake,
}

func runWake(c *cobra.Command, args []string) error {
	name := args[0]
	namespace := cmd.GetNamespace()

	ctx := context.Background()
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	// Get workspace
	sp := printer.NewSpinner("Fetching workspace")
	sp.Start()

	workspace, err := k8sClient.GetWorkspace(ctx, name, namespace)
	if err != nil {
		sp.Error("Failed to fetch workspace")
		return err
	}

	if workspace.Spec.Hibernated == nil || !*workspace.Spec.Hibernated {
		sp.Stop()
		fmt.Println()
		fmt.Println(styles.Info(fmt.Sprintf("Workspace %s is already awake", name)))
		fmt.Println()
		return nil
	}

	sp.Success("Workspace found")

	// Patch workspace to set hibernated=false
	sp = printer.NewSpinner("Waking up workspace")
	sp.Start()

	patch := client.MergeFrom(workspace.DeepCopy())
	awake := false
	workspace.Spec.Hibernated = &awake

	if err := k8sClient.PatchWorkspace(ctx, workspace, patch); err != nil {
		sp.Error("Failed to wake workspace")
		return err
	}

	sp.Success(fmt.Sprintf("Workspace %s is now awake", name))

	fmt.Println()
	fmt.Println(styles.MutedStyle.Render("All modules in this workspace will scale back to their original state."))
	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("Check status:"))
	fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code(fmt.Sprintf("forkspacer workspace get %s", name)))
	fmt.Println()

	return nil
}
