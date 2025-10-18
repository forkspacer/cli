package module

import (
	"context"
	"fmt"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/module"
	"github.com/forkspacer/cli/pkg/styles"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a module",
	Long: `Delete a Forkspacer module.

This will remove the module resource and uninstall the associated application.

Examples:
  # Delete a module
  forkspacer module delete my-module

  # Delete module in specific namespace
  forkspacer module delete my-module -n production`,
	Args: cobra.ExactArgs(1),
	RunE: runDelete,
}

func init() {
	moduleCmd.AddCommand(deleteCmd)
}

func runDelete(c *cobra.Command, args []string) error {
	name := args[0]
	namespace := cmd.GetNamespace()
	ctx := context.Background()

	service, err := module.NewService()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	fmt.Println()
	fmt.Printf("%s Deleting module %s in namespace %s...\n",
		styles.SymbolWarning,
		styles.Code(name),
		styles.Code(namespace))

	if err := service.Delete(ctx, name, &namespace); err != nil {
		return fmt.Errorf("failed to delete module: %w", err)
	}

	fmt.Println()
	fmt.Printf("%s Module %s deleted successfully\n", styles.SymbolSuccess, styles.Value(name))
	fmt.Println()

	return nil
}