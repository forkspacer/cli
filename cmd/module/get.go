package module

import (
	"context"
	"fmt"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/module"
	"github.com/forkspacer/cli/pkg/styles"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get details of a module",
	Long: `Get detailed information about a specific Forkspacer module.

Examples:
  # Get module details
  forkspacer module get my-module

  # Get module in specific namespace
  forkspacer module get my-module -n production`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func init() {
	moduleCmd.AddCommand(getCmd)
}

func runGet(c *cobra.Command, args []string) error {
	name := args[0]
	namespace := cmd.GetNamespace()
	ctx := context.Background()

	service, err := module.NewService()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	mod, err := service.Get(ctx, name, namespace)
	if err != nil {
		return fmt.Errorf("failed to get module: %w", err)
	}

	// Print details
	fmt.Println()
	fmt.Println(styles.TitleStyle.Render(fmt.Sprintf("Module: %s", name)))
	fmt.Println()
	fmt.Println(styles.Divider())
	fmt.Println()

	fmt.Printf("%s  %s\n", styles.Key("Name:"), styles.Value(mod.Name))
	fmt.Printf("%s  %s\n", styles.Key("Namespace:"), styles.Value(mod.Namespace))
	fmt.Printf("%s  %s\n", styles.Key("Phase:"), styles.Value(string(mod.Status.Phase)))

	if mod.Status.Message != nil {
		fmt.Printf("%s  %s\n", styles.Key("Message:"), styles.Value(*mod.Status.Message))
	}

	if mod.Status.LastActivity != nil {
		fmt.Printf("%s  %s\n", styles.Key("Last Activity:"), styles.Value(mod.Status.LastActivity.Format("2006-01-02 15:04:05")))
	}

	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("Workspace"))
	fmt.Printf("%s  %s\n", styles.Key("Name:"), styles.Value(mod.Spec.Workspace.Name))
	fmt.Printf("%s  %s\n", styles.Key("Namespace:"), styles.Value(mod.Spec.Workspace.Namespace))

	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("Source"))

	if mod.Spec.Helm != nil {
		fmt.Printf("%s  %s\n", styles.Key("Type:"), styles.Value("helm"))

		if mod.Spec.Helm.ExistingRelease != nil {
			fmt.Printf("%s  %s\n", styles.Key("Existing Release:"), styles.Value(mod.Spec.Helm.ExistingRelease.Name))
			if mod.Spec.Helm.ExistingRelease.Namespace != "" {
				fmt.Printf("%s  %s\n", styles.Key("Release Namespace:"), styles.Value(mod.Spec.Helm.ExistingRelease.Namespace))
			}
		}

		if mod.Spec.Helm.Chart.Repo != nil {
			fmt.Printf("%s  %s\n", styles.Key("Chart Repo:"), styles.Value(mod.Spec.Helm.Chart.Repo.URL))
			fmt.Printf("%s  %s\n", styles.Key("Chart Name:"), styles.Value(mod.Spec.Helm.Chart.Repo.Chart))
			if mod.Spec.Helm.Chart.Repo.Version != nil {
				fmt.Printf("%s  %s\n", styles.Key("Chart Version:"), styles.Value(*mod.Spec.Helm.Chart.Repo.Version))
			}
		} else if mod.Spec.Helm.Chart.Git != nil {
			fmt.Printf("%s  %s\n", styles.Key("Git Repo:"), styles.Value(mod.Spec.Helm.Chart.Git.Repo))
			fmt.Printf("%s  %s\n", styles.Key("Git Path:"), styles.Value(mod.Spec.Helm.Chart.Git.Path))
			fmt.Printf("%s  %s\n", styles.Key("Git Revision:"), styles.Value(mod.Spec.Helm.Chart.Git.Revision))
		} else if mod.Spec.Helm.Chart.ConfigMap != nil {
			fmt.Printf("%s  %s/%s\n", styles.Key("ConfigMap:"),
				styles.Value(mod.Spec.Helm.Chart.ConfigMap.Namespace),
				styles.Value(mod.Spec.Helm.Chart.ConfigMap.Name))
		}
	} else if mod.Spec.Custom != nil {
		fmt.Printf("%s  %s\n", styles.Key("Type:"), styles.Value("custom"))
		fmt.Printf("%s  %s\n", styles.Key("Image:"), styles.Value(mod.Spec.Custom.Image))
	}

	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("State"))
	hibernatedStatus := "active"
	if mod.Spec.Hibernated {
		hibernatedStatus = "hibernated"
	}
	fmt.Printf("%s  %s\n", styles.Key("Hibernated:"), styles.Value(hibernatedStatus))

	fmt.Println()

	return nil
}
