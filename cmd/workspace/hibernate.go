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

var hibernateCmd = &cobra.Command{
	Use:   "hibernate [name]",
	Short: "Hibernate a workspace",
	Long: `Hibernate a workspace to save resources.

Hibernation will:
  • Scale down all modules to zero replicas
  • Preserve all data and configuration
  • Save cluster resources and costs

The workspace can be woken up later with 'forkspacer workspace wake'.

Examples:
  # Hibernate a workspace
  forkspacer workspace hibernate dev-env

  # Hibernate workspace in specific namespace
  forkspacer workspace hibernate staging -n production`,
	Args: cobra.ExactArgs(1),
	RunE: runHibernate,
}

func runHibernate(c *cobra.Command, args []string) error {
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

	if workspace.Spec.Hibernated != nil && *workspace.Spec.Hibernated {
		sp.Stop()
		fmt.Println()
		fmt.Println(styles.Info(fmt.Sprintf("Workspace %s is already hibernated", name)))
		fmt.Println()
		return nil
	}

	sp.Success("Workspace found")

	// Patch workspace to set hibernated=true
	sp = printer.NewSpinner("Hibernating workspace")
	sp.Start()

	patch := client.MergeFrom(workspace.DeepCopy())
	hibernated := true
	workspace.Spec.Hibernated = &hibernated

	if err := k8sClient.PatchWorkspace(ctx, workspace, patch); err != nil {
		sp.Error("Failed to hibernate workspace")
		return err
	}

	sp.Success(fmt.Sprintf("Workspace %s hibernated successfully", name))

	fmt.Println()
	fmt.Println(styles.MutedStyle.Render("All modules in this workspace will scale down to zero replicas."))
	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("To wake up:"))
	fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code(fmt.Sprintf("forkspacer workspace wake %s", name)))
	fmt.Println()

	return nil
}
