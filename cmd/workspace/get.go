package workspace

import (
	"context"
	"fmt"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/styles"
	workspaceService "github.com/forkspacer/cli/pkg/workspace"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get workspace details",
	Long: `Display detailed information about a specific workspace.

Examples:
  # Get workspace in default namespace
  forkspacer workspace get dev-env

  # Get workspace in specific namespace
  forkspacer workspace get dev-env -n production`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func runGet(c *cobra.Command, args []string) error {
	name := args[0]
	namespace := cmd.GetNamespace()

	ctx := context.Background()
	service, err := workspaceService.NewService()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	workspace, err := service.Get(ctx, name, namespace)
	if err != nil {
		return err
	}

	// Print detailed workspace info
	fmt.Println()
	fmt.Println(styles.TitleStyle.Render(fmt.Sprintf("Workspace: %s", workspace.Name)))
	fmt.Println()

	// Metadata
	fmt.Println(styles.KeyStyle.Render("Metadata"))
	fmt.Printf("  %s  %s\n", styles.Key("Name:"), styles.Value(workspace.Name))
	fmt.Printf("  %s  %s\n", styles.Key("Namespace:"), styles.Value(workspace.Namespace))
	fmt.Printf("  %s  %s\n", styles.Key("UID:"), styles.Value(string(workspace.UID)))
	fmt.Printf("  %s  %s\n", styles.Key("Created:"), styles.Value(workspace.CreationTimestamp.Format("2006-01-02 15:04:05")))
	fmt.Println()

	// Spec
	fmt.Println(styles.KeyStyle.Render("Specification"))
	fmt.Printf("  %s  %s\n", styles.Key("Type:"), styles.Value(string(workspace.Spec.Type)))
	fmt.Printf("  %s  %s\n", styles.Key("Connection:"), styles.Value(string(workspace.Spec.Connection.Type)))

	hibernated := "false"
	if workspace.Spec.Hibernated != nil && *workspace.Spec.Hibernated {
		hibernated = "true"
	}
	fmt.Printf("  %s  %s\n", styles.Key("Hibernated:"), styles.Value(hibernated))

	if workspace.Spec.AutoHibernation != nil {
		fmt.Println()
		fmt.Println(styles.KeyStyle.Render("Auto-Hibernation"))
		fmt.Printf("  %s  %t\n", styles.Key("Enabled:"), workspace.Spec.AutoHibernation.Enabled)
		if workspace.Spec.AutoHibernation.Enabled {
			fmt.Printf("  %s  %s\n", styles.Key("Sleep Schedule:"), styles.Value(workspace.Spec.AutoHibernation.Schedule))
			if workspace.Spec.AutoHibernation.WakeSchedule != nil {
				fmt.Printf("  %s  %s\n", styles.Key("Wake Schedule:"), styles.Value(*workspace.Spec.AutoHibernation.WakeSchedule))
			}
		}
	}

	if workspace.Spec.From != nil {
		fmt.Println()
		fmt.Println(styles.KeyStyle.Render("Forked From"))
		fmt.Printf("  %s  %s\n", styles.Key("Workspace:"), styles.Value(workspace.Spec.From.Name))
		fmt.Printf("  %s  %s\n", styles.Key("Namespace:"), styles.Value(workspace.Spec.From.Namespace))
	}

	fmt.Println()

	// Status
	statusStyle := styles.ValueStyle
	switch workspace.Status.Phase {
	case "ready":
		statusStyle = styles.SuccessStyle
	case "failed":
		statusStyle = styles.ErrorStyle
	case "hibernated":
		statusStyle = styles.WarningStyle
	}

	fmt.Println(styles.KeyStyle.Render("Status"))
	fmt.Printf("  %s  %s\n", styles.Key("Phase:"), statusStyle.Render(string(workspace.Status.Phase)))
	fmt.Printf("  %s  %t\n", styles.Key("Ready:"), workspace.Status.Ready)

	if workspace.Status.LastActivity != nil {
		fmt.Printf("  %s  %s\n", styles.Key("Last Activity:"),
			styles.Value(workspace.Status.LastActivity.Format("2006-01-02 15:04:05")))
	}

	if workspace.Status.HibernatedAt != nil {
		fmt.Printf("  %s  %s\n", styles.Key("Hibernated At:"),
			styles.Value(workspace.Status.HibernatedAt.Format("2006-01-02 15:04:05")))
	}

	if workspace.Status.Message != nil {
		fmt.Printf("  %s  %s\n", styles.Key("Message:"), styles.Value(*workspace.Status.Message))
	}

	fmt.Println()

	return nil
}
