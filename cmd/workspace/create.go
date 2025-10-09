package workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	batchv1 "github.com/forkspacer/forkspacer/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/k8s"
	"github.com/forkspacer/cli/pkg/printer"
	"github.com/forkspacer/cli/pkg/styles"
	"github.com/forkspacer/cli/pkg/validation"
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

	// Step 3: Connect to cluster
	sp = printer.NewSpinner("Connecting to Kubernetes cluster")
	sp.Start()

	ctx := context.Background()
	client, err := k8s.NewClient()
	if err != nil {
		sp.Error("Failed to connect to cluster")
		return fmt.Errorf("kubernetes connection failed: %w", err)
	}
	sp.Success(fmt.Sprintf("Connected to cluster (context: %s)", client.Context))

	// Step 4: Check if operator is installed
	sp = printer.NewSpinner("Checking Forkspacer operator installation")
	sp.Start()
	time.Sleep(200 * time.Millisecond)

	if err := client.CheckOperatorInstalled(ctx); err != nil {
		sp.Error("Forkspacer operator not found")
		return fmt.Errorf("operator not installed: %w\n\nInstall with: helm install forkspacer forkspacer/forkspacer", err)
	}
	sp.Success("Forkspacer operator is installed")

	// Step 5: Check if workspace already exists
	sp = printer.NewSpinner("Checking if workspace already exists")
	sp.Start()
	time.Sleep(200 * time.Millisecond)

	exists, err := client.WorkspaceExists(ctx, name, namespace)
	if err != nil {
		sp.Error("Failed to check workspace existence")
		return err
	}
	if exists {
		sp.Error("Workspace already exists")
		return fmt.Errorf("workspace %s/%s already exists\n\nUse: forkspacer workspace get %s", namespace, name, name)
	}
	sp.Success("Workspace name is available")

	// Step 6: Build workspace object
	workspace := buildWorkspace(name, namespace)

	// Step 7: Create workspace
	sp = printer.NewSpinner("Creating workspace resource")
	sp.Start()

	if err := client.CreateWorkspace(ctx, workspace); err != nil {
		sp.Error("Failed to create workspace")
		return err
	}
	sp.Success("Workspace resource created")

	// Step 8: Wait for ready (optional)
	if createWait {
		sp = printer.NewSpinner("Waiting for workspace to become ready")
		sp.Start()

		if err := waitForWorkspaceReady(ctx, client, name, namespace, 2*time.Minute); err != nil {
			sp.Error("Workspace did not become ready")
			return err
		}
		sp.Success("Workspace is ready")
	}

	// Print summary
	printSuccessSummary(workspace)

	return nil
}

func buildWorkspace(name, namespace string) *batchv1.Workspace {
	workspace := &batchv1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.WorkspaceSpec{
			Type: batchv1.WorkspaceTypeKubernetes,
			Connection: &batchv1.WorkspaceConnection{
				Type: batchv1.WorkspaceConnectionType(createConnectionType),
			},
		},
	}

	// Add auto-hibernation if specified
	if createHibernationSched != "" {
		workspace.Spec.AutoHibernation = &batchv1.WorkspaceAutoHibernation{
			Enabled:  true,
			Schedule: createHibernationSched,
		}
		if createWakeSched != "" {
			workspace.Spec.AutoHibernation.WakeSchedule = &createWakeSched
		}
	}

	// Add fork reference if specified
	if createFromWorkspace != "" {
		workspace.Spec.From = &batchv1.WorkspaceFromReference{
			Name:      createFromWorkspace,
			Namespace: namespace,
		}
	}

	return workspace
}

func waitForWorkspaceReady(ctx context.Context, client *k8s.Client, name, namespace string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		ws, err := client.GetWorkspace(ctx, name, namespace)
		if err != nil {
			return err
		}

		if ws.Status.Phase == batchv1.WorkspacePhaseReady && ws.Status.Ready {
			return nil
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("timeout waiting for workspace to become ready")
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

	return fmt.Errorf(msg)
}

func formatCronError(schedule string, err error) error {
	msg := fmt.Sprintf("\n%s\n\n", styles.Error("Invalid cron schedule"))
	msg += fmt.Sprintf("  The schedule %s is not valid.\n\n", styles.Code(schedule))
	msg += fmt.Sprintf("  %s\n", styles.Key("Common schedules:"))
	for schedule, desc := range validation.CronExamples() {
		msg += fmt.Sprintf("    %s %-20s %s\n", styles.SymbolBullet, styles.Code(schedule), desc)
	}
	msg += fmt.Sprintf("\n  %s https://crontab.guru\n", styles.Key("Learn more:"))

	return fmt.Errorf(msg)
}
