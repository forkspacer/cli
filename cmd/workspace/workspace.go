package workspace

import (
	"github.com/spf13/cobra"
	"github.com/forkspacer/cli/cmd"
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
