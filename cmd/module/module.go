package module

import (
	"github.com/forkspacer/cli/cmd"
	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Manage Forkspacer modules",
	Long: `Manage Forkspacer modules (applications) within workspaces.

Modules represent applications deployed to workspaces. They can be:
  • Imported from existing Helm releases
  • Deployed from Helm charts
  • Managed through their lifecycle (install, uninstall, hibernate)

Examples:
  # Import existing Helm release
  forkspacer import my-module --helm-release my-release --workspace dev-env

  # List modules
  forkspacer module list

  # Get module details
  forkspacer module get my-module`,
}

func init() {
	cmd.GetRootCmd().AddCommand(moduleCmd)
}