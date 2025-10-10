package workspace

import (
	"context"
	"fmt"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/printer"
	"github.com/forkspacer/cli/pkg/styles"
	workspaceService "github.com/forkspacer/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

var (
	listAllNamespaces bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all workspaces",
	Aliases: []string{"ls"},
	Long: `List all Forkspacer workspaces in the current or specified namespace.

Examples:
  # List workspaces in default namespace
  forkspacer workspace list

  # List in specific namespace
  forkspacer workspace list -n production

  # List across all namespaces
  forkspacer workspace list --all-namespaces`,
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVarP(&listAllNamespaces, "all-namespaces", "A", false,
		"List workspaces across all namespaces")
}

func runList(c *cobra.Command, args []string) error {
	namespace := cmd.GetNamespace()
	if listAllNamespaces {
		namespace = "" // Empty means all namespaces
	}

	ctx := context.Background()
	service, err := workspaceService.NewService()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	workspaces, err := service.List(ctx, namespace)
	if err != nil {
		return err
	}

	if len(workspaces.Items) == 0 {
		fmt.Println()
		fmt.Println(styles.MutedStyle.Render("No workspaces found"))
		fmt.Println()
		fmt.Println(styles.SubtitleStyle.Render("Get started:"))
		fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code("forkspacer workspace create dev-env"))
		fmt.Println()
		return nil
	}

	// Print table
	fmt.Println()
	table := printer.NewTable([]string{"NAME", "NAMESPACE", "PHASE", "READY", "HIBERNATED", "LAST ACTIVITY"})

	for _, ws := range workspaces.Items {
		hibernated := "false"
		if ws.Spec.Hibernated != nil && *ws.Spec.Hibernated {
			hibernated = "true"
		}

		lastActivity := "never"
		if ws.Status.LastActivity != nil {
			lastActivity = ws.Status.LastActivity.Format("2006-01-02 15:04:05")
		}

		// Color code the phase
		phase := string(ws.Status.Phase)
		ready := fmt.Sprintf("%t", ws.Status.Ready)

		table.AddRow([]string{
			ws.Name,
			ws.Namespace,
			phase,
			ready,
			hibernated,
			lastActivity,
		})
	}

	table.Render()
	fmt.Println()
	fmt.Printf(styles.MutedStyle.Render("Total: %d workspace(s)"), len(workspaces.Items))
	fmt.Println()
	fmt.Println()

	return nil
}
