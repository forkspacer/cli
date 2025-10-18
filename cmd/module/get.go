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

	if mod.Spec.Source.ExistingHelmRelease != nil {
		fmt.Printf("%s  %s\n", styles.Key("Type:"), styles.Value("existing-helm-release"))
		fmt.Printf("%s  %s\n", styles.Key("Release Name:"), styles.Value(mod.Spec.Source.ExistingHelmRelease.Name))
		if mod.Spec.Source.ExistingHelmRelease.Namespace != "" {
			fmt.Printf("%s  %s\n", styles.Key("Release Namespace:"), styles.Value(mod.Spec.Source.ExistingHelmRelease.Namespace))
		}
	} else if mod.Spec.Source.HttpURL != nil {
		fmt.Printf("%s  %s\n", styles.Key("Type:"), styles.Value("http"))
		fmt.Printf("%s  %s\n", styles.Key("URL:"), styles.Value(*mod.Spec.Source.HttpURL))
	} else if mod.Spec.Source.Raw != nil {
		fmt.Printf("%s  %s\n", styles.Key("Type:"), styles.Value("raw"))
	} else if mod.Spec.Source.ConfigMap != nil {
		fmt.Printf("%s  %s\n", styles.Key("Type:"), styles.Value("configmap"))
		fmt.Printf("%s  %s/%s\n", styles.Key("ConfigMap:"),
			styles.Value(mod.Spec.Source.ConfigMap.Namespace),
			styles.Value(mod.Spec.Source.ConfigMap.Name))
	}

	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("State"))
	hibernatedStatus := "active"
	if mod.Spec.Hibernated != nil && *mod.Spec.Hibernated {
		hibernatedStatus = "hibernated"
	}
	fmt.Printf("%s  %s\n", styles.Key("Hibernated:"), styles.Value(hibernatedStatus))

	fmt.Println()

	return nil
}