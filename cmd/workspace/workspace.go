package workspace

import (
	"context"

	"github.com/forkspacer/cli/cmd"
	workspaceService "github.com/forkspacer/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

// WorkspaceCmd represents the workspace command
var WorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage Forkspacer workspaces",
	Long: `Manage Forkspacer workspaces - isolated Kubernetes environments with lifecycle management.

Workspaces provide:
  • Environment isolation
  • Auto-hibernation scheduling
  • Resource management
  • Workspace forking`,
	Aliases: []string{"ws"},
}

func init() {
	cmd.GetRootCmd().AddCommand(WorkspaceCmd)

	// Add subcommands
	WorkspaceCmd.AddCommand(createCmd)
	WorkspaceCmd.AddCommand(listCmd)
	WorkspaceCmd.AddCommand(getCmd)
	WorkspaceCmd.AddCommand(deleteCmd)
	WorkspaceCmd.AddCommand(hibernateCmd)
	WorkspaceCmd.AddCommand(wakeCmd)
}

// workspaceNameCompletion provides dynamic completion for workspace names
func workspaceNameCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Only complete the first argument (workspace name)
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Create workspace service
	ctx := context.Background()
	service, err := workspaceService.NewService()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Get namespace from flag
	namespace := cmd.Flag("namespace").Value.String()

	// List workspaces in namespace
	workspaces, err := service.List(ctx, namespace)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Extract workspace names
	var names []string
	for _, ws := range workspaces.Items {
		names = append(names, ws.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
