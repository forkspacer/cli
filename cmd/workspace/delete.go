package workspace

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/k8s"
	"github.com/forkspacer/cli/pkg/printer"
	"github.com/forkspacer/cli/pkg/styles"
)

var (
	deleteForce bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a workspace",
	Long: `Delete a Forkspacer workspace and all its modules.

WARNING: This will delete all modules associated with this workspace.

Examples:
  # Delete workspace (with confirmation)
  forkspacer workspace delete dev-env

  # Delete without confirmation
  forkspacer workspace delete dev-env --force`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false,
		"Skip confirmation prompt")
}

func runDelete(c *cobra.Command, args []string) error {
	name := args[0]
	namespace := cmd.GetNamespace()

	ctx := context.Background()
	client, err := k8s.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	// Check if workspace exists
	workspace, err := client.GetWorkspace(ctx, name, namespace)
	if err != nil {
		return err
	}

	// Confirm deletion unless --force
	if !deleteForce {
		fmt.Println()
		fmt.Println(styles.WarningStyle.Render(fmt.Sprintf("⚠  About to delete workspace: %s/%s", namespace, name)))
		fmt.Println()
		fmt.Println(styles.MutedStyle.Render("This will also delete all modules in this workspace."))
		fmt.Println()
		fmt.Print("Continue? (y/N): ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println()
			fmt.Println(styles.Info("Deletion cancelled"))
			fmt.Println()
			return nil
		}
	}

	// Delete workspace
	sp := printer.NewSpinner("Deleting workspace")
	sp.Start()

	if err := client.DeleteWorkspace(ctx, workspace.Name, workspace.Namespace); err != nil {
		sp.Error("Failed to delete workspace")
		return err
	}

	sp.Success(fmt.Sprintf("Workspace %s deleted successfully", name))
	fmt.Println()

	return nil
}
