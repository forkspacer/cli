package workspace

import (
	"context"
	"fmt"
	"time"

	batchv1 "github.com/forkspacer/forkspacer/api/v1"
	"github.com/spf13/cobra"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/printer"
	"github.com/forkspacer/cli/pkg/styles"
	"github.com/forkspacer/cli/pkg/validation"
	workspaceService "github.com/forkspacer/cli/pkg/workspace"
)

var (
	createConnectionType   string
	createHibernationSched string
	createWakeSched        string
	createFromWorkspace    string
	createMigrateData      bool
	createWait             bool
)

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new workspace",
	Long: `Create a new Forkspacer workspace with optional auto-hibernation.

A workspace is an isolated Kubernetes environment that can:
  • Automatically hibernate on a schedule to save costs
  • Be forked from existing workspaces
  • Manage multiple modules (applications)

Examples:
  # Create a simple workspace
  forkspacer workspace create dev-env

  # Create with auto-hibernation
  forkspacer workspace create dev-env \
    --hibernation-schedule "0 18 * * *" \
    --wake-schedule "0 8 * * *"

  # Fork from existing workspace
  forkspacer workspace create staging \
    --from production \
    --migrate-data`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	createCmd.Flags().StringVar(&createConnectionType, "connection", "in-cluster",
		"Connection type (local|in-cluster|kubeconfig)")
	createCmd.Flags().StringVar(&createHibernationSched, "hibernation-schedule", "",
		"Hibernation cron schedule (e.g., '0 18 * * *' for 6 PM daily)")
	createCmd.Flags().StringVar(&createWakeSched, "wake-schedule", "",
		"Wake cron schedule (e.g., '0 8 * * *' for 8 AM daily)")
	createCmd.Flags().StringVar(&createFromWorkspace, "from", "",
		"Fork from existing workspace")
	createCmd.Flags().BoolVar(&createMigrateData, "migrate-data", false,
		"Migrate PV data when forking (requires --from)")
	createCmd.Flags().BoolVar(&createWait, "wait", false,
		"Wait for workspace to become ready")
}

func runCreate(c *cobra.Command, args []string) error {
	name := args[0]
	namespace := cmd.GetNamespace()

	// Print header
	fmt.Println()
	fmt.Println(styles.TitleStyle.Render(fmt.Sprintf("%s Creating workspace %s", styles.SymbolSparkles, name)))
	fmt.Println()

	// Step 1: Validate name
	sp := printer.NewSpinner("Validating workspace name")
	sp.Start()
	time.Sleep(200 * time.Millisecond) // Brief pause for UX

	if err := validation.ValidateDNS1123Subdomain(name); err != nil {
		sp.Stop()
		return formatValidationError(name, err)
	}
	sp.Success("Workspace name is valid")

	// Step 2: Validate hibernation schedule if provided
	if createHibernationSched != "" {
		sp = printer.NewSpinner("Validating hibernation schedule")
		sp.Start()
		time.Sleep(200 * time.Millisecond)

		if err := validation.ValidateCronSchedule(createHibernationSched); err != nil {
			sp.Stop()
			return formatCronError(createHibernationSched, err)
		}
		sp.Success("Hibernation schedule is valid")
	}

	if createWakeSched != "" {
		if createHibernationSched == "" {
			return fmt.Errorf("--wake-schedule requires --hibernation-schedule")
		}

		sp = printer.NewSpinner("Validating wake schedule")
		sp.Start()
		time.Sleep(200 * time.Millisecond)

		if err := validation.ValidateCronSchedule(createWakeSched); err != nil {
			sp.Stop()
			return formatCronError(createWakeSched, err)
		}
		sp.Success("Wake schedule is valid")
	}

	// Step 3: Connect to cluster and create service
	sp = printer.NewSpinner("Connecting to Kubernetes cluster")
	sp.Start()

	ctx := context.Background()
	service, err := workspaceService.NewService()
	if err != nil {
		sp.Error("Failed to connect to cluster")
		return fmt.Errorf("kubernetes connection failed: %w", err)
	}
	sp.Success("Connected to cluster")

	// Step 4: Build workspace input
	workspaceIn := buildWorkspaceInput(name, namespace)

	// Step 5: Create workspace using api-server service
	sp = printer.NewSpinner("Creating workspace resource")
	sp.Start()

	workspace, err := service.Create(ctx, workspaceIn)
	if err != nil {
		sp.Error("Failed to create workspace")
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	sp.Success("Workspace resource created")

	// Step 6: Wait for ready (optional)
	if createWait {
		sp = printer.NewSpinner("Waiting for workspace to become ready")
		sp.Start()

		if err := waitForWorkspaceReady(ctx, service, name, namespace, 2*time.Minute); err != nil {
			sp.Error("Workspace did not become ready")
			return err
		}
		sp.Success("Workspace is ready")
	}

	// Print summary
	printSuccessSummary(workspace)

	return nil
}

func buildWorkspaceInput(name, namespace string) workspaceService.WorkspaceCreateInput {
	workspaceIn := workspaceService.WorkspaceCreateInput{
		Name:           name,
		Namespace:      namespace,
		Hibernated:     false,
		ConnectionType: createConnectionType,
	}

	// Add auto-hibernation if specified
	if createHibernationSched != "" {
		workspaceIn.AutoHibernation = &workspaceService.AutoHibernationInput{
			Enabled:  true,
			Schedule: createHibernationSched,
		}
		if createWakeSched != "" {
			workspaceIn.AutoHibernation.WakeSchedule = &createWakeSched
		}
	}

	// Add fork reference if specified
	if createFromWorkspace != "" {
		workspaceIn.From = &workspaceService.FromWorkspaceInput{
			Name:      createFromWorkspace,
			Namespace: namespace,
		}
		// TODO: Wire up createMigrateData flag when API server supports it
		// Currently the flag is defined but not yet supported in WorkspaceCreateInput
		_ = createMigrateData
	}

	return workspaceIn
}

func waitForWorkspaceReady(ctx context.Context, service *workspaceService.Service, name, namespace string, timeout time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for workspace to become ready")
		case <-ticker.C:
			workspace, err := service.Get(ctx, name, namespace)
			if err != nil {
				continue // Workspace might not exist yet, keep waiting
			}

			if workspace.Status.Ready {
				return nil
			}

			// Check if workspace is in a failed state
			if workspace.Status.Phase == "failed" {
				if workspace.Status.Message != nil {
					return fmt.Errorf("workspace failed: %s", *workspace.Status.Message)
				}
				return fmt.Errorf("workspace entered failed state")
			}
		}
	}
}

func printSuccessSummary(workspace *batchv1.Workspace) {
	fmt.Println()
	fmt.Println(styles.Divider())
	fmt.Println()

	fmt.Printf("%s  %s\n", styles.Key("Name:"), styles.Value(workspace.Name))
	fmt.Printf("%s  %s\n", styles.Key("Namespace:"), styles.Value(workspace.Namespace))
	fmt.Printf("%s  %s\n", styles.Key("Type:"), styles.Value(string(workspace.Spec.Type)))

	if workspace.Spec.AutoHibernation != nil && workspace.Spec.AutoHibernation.Enabled {
		fmt.Printf("%s  %s\n", styles.Key("Hibernation:"), styles.Value("enabled"))
		fmt.Printf("  %s  %s\n", styles.Key("Sleep:"), styles.Value(workspace.Spec.AutoHibernation.Schedule))
		if workspace.Spec.AutoHibernation.WakeSchedule != nil {
			fmt.Printf("  %s  %s\n", styles.Key("Wake:"), styles.Value(*workspace.Spec.AutoHibernation.WakeSchedule))
		}
	} else {
		fmt.Printf("%s  %s\n", styles.Key("Hibernation:"), styles.Value("disabled"))
	}

	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("Next steps:"))
	fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code(fmt.Sprintf("forkspacer workspace get %s", workspace.Name)))
	fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code(fmt.Sprintf("forkspacer module deploy redis --workspace %s", workspace.Name)))

	fmt.Println()
	fmt.Println(styles.MutedStyle.Render("Documentation: https://forkspacer.com/docs/workspaces"))
	fmt.Println()
}

func formatValidationError(name string, err error) error {
	msg := fmt.Sprintf("\n%s\n\n", styles.Error("Invalid workspace name"))
	msg += fmt.Sprintf("  The name %s doesn't meet DNS-1123 requirements.\n\n", styles.Code(name))
	msg += fmt.Sprintf("  %s\n", styles.Key("Requirements:"))
	for _, req := range validation.DNS1123Requirements() {
		msg += fmt.Sprintf("    %s %s\n", styles.SymbolBullet, req)
	}
	msg += fmt.Sprintf("\n  %s\n", styles.Key("Valid examples:"))
	for _, example := range validation.DNS1123Examples() {
		msg += fmt.Sprintf("    %s %s\n", styles.SymbolBullet, styles.Code(example))
	}
	msg += fmt.Sprintf("\n  %s\n", styles.Key("Try:"))
	msg += fmt.Sprintf("    %s\n", styles.Code("forkspacer workspace create dev-env"))

	return fmt.Errorf("%s", msg)
}

func formatCronError(schedule string, err error) error {
	msg := fmt.Sprintf("\n%s\n\n", styles.Error("Invalid cron schedule"))
	msg += fmt.Sprintf("  The schedule %s is not valid.\n\n", styles.Code(schedule))
	msg += fmt.Sprintf("  %s\n", styles.Key("Common schedules:"))
	for schedule, desc := range validation.CronExamples() {
		msg += fmt.Sprintf("    %s %-20s %s\n", styles.SymbolBullet, styles.Code(schedule), desc)
	}
	msg += fmt.Sprintf("\n  %s https://crontab.guru\n", styles.Key("Learn more:"))

	return fmt.Errorf("%s", msg)
}
