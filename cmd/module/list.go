package module

import (
	"context"
	"fmt"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/module"
	"github.com/forkspacer/cli/pkg/printer"
	"github.com/forkspacer/cli/pkg/styles"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all modules",
	Long: `List all Forkspacer modules in the specified namespace.

Examples:
  # List modules in default namespace
  forkspacer module list

  # List modules in specific namespace
  forkspacer module list -n production`,
	RunE: runList,
}

func init() {
	moduleCmd.AddCommand(listCmd)
}

func runList(c *cobra.Command, args []string) error {
	namespace := cmd.GetNamespace()
	ctx := context.Background()

	service, err := module.NewService()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	modules, err := service.List(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to list modules: %w", err)
	}

	if len(modules.Items) == 0 {
		fmt.Println()
		fmt.Println(styles.MutedStyle.Render(fmt.Sprintf("No modules found in namespace '%s'", namespace)))
		fmt.Println()
		return nil
	}

	// Print table
	fmt.Println()
	table := printer.NewTable([]string{"NAME", "NAMESPACE", "WORKSPACE", "PHASE", "LAST ACTIVITY"})

	for _, mod := range modules.Items {
		workspace := fmt.Sprintf("%s/%s", mod.Spec.Workspace.Namespace, mod.Spec.Workspace.Name)
		lastActivity := "never"
		if mod.Status.LastActivity != nil {
			lastActivity = mod.Status.LastActivity.Format("2006-01-02 15:04:05")
		}

		table.AddRow([]string{
			mod.Name,
			mod.Namespace,
			workspace,
			string(mod.Status.Phase),
			lastActivity,
		})
	}

	table.Render()
	fmt.Println()
	fmt.Printf(styles.MutedStyle.Render("Total: %d module(s)"), len(modules.Items))
	fmt.Println()
	fmt.Println()

	return nil
}
