package module

import (
	"context"
	"fmt"
	"time"

	batchv1 "github.com/forkspacer/forkspacer/api/v1"
	"github.com/spf13/cobra"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/module"
	"github.com/forkspacer/cli/pkg/printer"
	"github.com/forkspacer/cli/pkg/styles"
	"github.com/forkspacer/cli/pkg/validation"
)

var (
	addHelmRelease          string
	addHelmReleaseNamespace string
	addWorkspace            string
	addWorkspaceNamespace   string
	addHibernated           bool
	addWait                 bool
	addChartGitRepo         string
	addChartGitPath         string
	addChartGitRevision     string
)

var addCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add an existing Helm release as a Forkspacer module",
	Long: `Add an existing Helm release to be managed by Forkspacer.

This creates a Module resource that references an existing Helm release,
allowing Forkspacer to manage its lifecycle (hibernation, forking, etc.).

The added module will:
  • Reference the existing Helm release
  • Be associated with a workspace
  • Support hibernation and lifecycle management

Examples:
  # Add a Helm release from the default namespace
  forkspacer module add my-module \
    --helm-release my-release \
    --workspace dev-env \
    --chart-git-repo https://github.com/org/repo \
    --chart-git-path charts/app

  # Add from a specific namespace
  forkspacer module add my-module \
    --helm-release my-release \
    --helm-release-namespace apps \
    --workspace dev-env \
    --workspace-namespace workspaces \
    --chart-git-repo https://github.com/org/repo \
    --chart-git-path charts/app

  # Add in hibernated state
  forkspacer module add my-module \
    --helm-release my-release \
    --workspace dev-env \
    --chart-git-repo https://github.com/org/repo \
    --chart-git-path charts/app \
    --hibernated`,
	Args: validateAddArgs,
	RunE: runAdd,
}

func validateAddArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("accepts 1 arg(s), received 0\n\nUsage:\n  forkspacer module add <name> --helm-release <release> --workspace <workspace> --chart-git-repo <repo> --chart-git-path <path>\n\nExample:\n  forkspacer module add my-module --helm-release my-release --workspace dev-env --chart-git-repo https://github.com/org/repo --chart-git-path charts/app")
	}
	if len(args) > 1 {
		return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
	}
	return nil
}

func init() {
	addCmd.Flags().StringVar(&addHelmRelease, "helm-release", "",
		"Name of the existing Helm release to add (required)")
	addCmd.Flags().StringVar(&addHelmReleaseNamespace, "helm-release-namespace", "default",
		"Namespace of the Helm release")
	addCmd.Flags().StringVar(&addWorkspace, "workspace", "",
		"Workspace to associate this module with (required)")
	addCmd.Flags().StringVar(&addWorkspaceNamespace, "workspace-namespace", "",
		"Namespace of the workspace (defaults to module namespace)")
	addCmd.Flags().BoolVar(&addHibernated, "hibernated", false,
		"Add in hibernated state")
	addCmd.Flags().BoolVar(&addWait, "wait", false,
		"Wait for module to become ready")

	// ChartSource Git flags
	addCmd.Flags().StringVar(&addChartGitRepo, "chart-git-repo", "",
		"Git repository URL for the Helm chart source (required)")
	addCmd.Flags().StringVar(&addChartGitPath, "chart-git-path", "",
		"Path to chart directory in the Git repository (required)")
	addCmd.Flags().StringVar(&addChartGitRevision, "chart-git-revision", "main",
		"Git revision (branch, tag, or commit)")

	addCmd.MarkFlagRequired("helm-release")
	addCmd.MarkFlagRequired("workspace")
	addCmd.MarkFlagRequired("chart-git-repo")
	addCmd.MarkFlagRequired("chart-git-path")

	moduleCmd.AddCommand(addCmd)
}

func runAdd(c *cobra.Command, args []string) error {
	name := args[0]
	namespace := cmd.GetNamespace()

	// Default workspace namespace to module namespace if not specified
	if addWorkspaceNamespace == "" {
		addWorkspaceNamespace = namespace
	}

	// Print header
	fmt.Println()
	fmt.Println(styles.TitleStyle.Render(fmt.Sprintf("%s Adding module %s", styles.SymbolSparkles, name)))
	fmt.Println()

	// Step 1: Validate module name
	sp := printer.NewSpinner("Validating module name")
	sp.Start()
	time.Sleep(200 * time.Millisecond) // Brief pause for UX

	if err := validation.ValidateDNS1123Subdomain(name); err != nil {
		sp.Stop()
		return formatAddValidationError(name, err)
	}
	sp.Success("Module name is valid")

	// Step 2: Validate Helm release name
	sp = printer.NewSpinner("Validating Helm release name")
	sp.Start()
	time.Sleep(200 * time.Millisecond)

	if err := validation.ValidateDNS1123Subdomain(addHelmRelease); err != nil {
		sp.Stop()
		return formatAddValidationError(addHelmRelease, err)
	}
	sp.Success("Helm release name is valid")

	// Step 3: Connect to cluster and create service
	sp = printer.NewSpinner("Connecting to Kubernetes cluster")
	sp.Start()

	ctx := context.Background()
	service, err := module.NewService()
	if err != nil {
		sp.Error("Failed to connect to cluster")
		return fmt.Errorf("kubernetes connection failed: %w", err)
	}
	sp.Success("Connected to cluster")

	// Step 4: Create module resource
	sp = printer.NewSpinner("Creating module resource")
	sp.Start()

	moduleResource, err := service.CreateExistingHelmRelease(
		ctx,
		name,
		namespace,
		addHelmRelease,
		addHelmReleaseNamespace,
		addWorkspace,
		addWorkspaceNamespace,
		addHibernated,
		addChartGitRepo,
		addChartGitPath,
		addChartGitRevision,
	)
	if err != nil {
		sp.Error("Failed to create module")
		return fmt.Errorf("failed to create module: %w", err)
	}
	sp.Success("Module resource created")

	// Step 5: Wait for ready (optional)
	if addWait {
		sp = printer.NewSpinner("Waiting for module to become ready")
		sp.Start()

		if err := waitForModuleReady(ctx, service, name, namespace, 2*time.Minute); err != nil {
			sp.Error("Module did not become ready")
			return err
		}
		sp.Success("Module is ready")
	}

	// Print summary
	printAddSummary(moduleResource)

	return nil
}

func waitForModuleReady(ctx context.Context, service *module.Service, name, namespace string, timeout time.Duration) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-timeoutCh:
			return fmt.Errorf("timeout waiting for module to become ready")
		case <-ticker.C:
			mod, err := service.Get(ctx, name, namespace)
			if err != nil {
				continue // Module might not exist yet, keep waiting
			}

			if mod.Status.Phase == batchv1.ModulePhaseReady {
				return nil
			}

			// Check if module is in a failed state
			if mod.Status.Phase == batchv1.ModulePhaseFailed {
				if mod.Status.Message != nil {
					return fmt.Errorf("module failed: %s", *mod.Status.Message)
				}
				return fmt.Errorf("module entered failed state")
			}
		}
	}
}

func printAddSummary(mod *batchv1.Module) {
	fmt.Println()
	fmt.Println(styles.Divider())
	fmt.Println()

	fmt.Printf("%s  %s\n", styles.Key("Name:"), styles.Value(mod.Name))
	fmt.Printf("%s  %s\n", styles.Key("Namespace:"), styles.Value(mod.Namespace))

	if mod.Spec.Source.ExistingHelmRelease != nil {
		fmt.Printf("%s  %s\n", styles.Key("Source:"), styles.Value("existing-helm-release"))
		fmt.Printf("  %s  %s\n", styles.Key("Release:"), styles.Value(mod.Spec.Source.ExistingHelmRelease.Name))
		if mod.Spec.Source.ExistingHelmRelease.Namespace != "" {
			fmt.Printf("  %s  %s\n", styles.Key("Release Namespace:"), styles.Value(mod.Spec.Source.ExistingHelmRelease.Namespace))
		}
	}

	fmt.Printf("%s  %s/%s\n",
		styles.Key("Workspace:"),
		styles.Value(mod.Spec.Workspace.Namespace),
		styles.Value(mod.Spec.Workspace.Name))

	hibernatedStatus := "active"
	if mod.Spec.Hibernated {
		hibernatedStatus = "hibernated"
	}
	fmt.Printf("%s  %s\n", styles.Key("State:"), styles.Value(hibernatedStatus))

	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("Next steps:"))
	fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code(fmt.Sprintf("forkspacer module get %s", mod.Name)))
	fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code(fmt.Sprintf("forkspacer workspace get %s", mod.Spec.Workspace.Name)))

	fmt.Println()
	fmt.Println(styles.MutedStyle.Render("Documentation: https://forkspacer.com/docs/modules"))
	fmt.Println()
}

func formatAddValidationError(name string, err error) error {
	msg := fmt.Sprintf("\n%s\n\n", styles.Error("Invalid name"))
	msg += fmt.Sprintf("  The name %s doesn't meet DNS-1123 requirements.\n\n", styles.Code(name))
	msg += fmt.Sprintf("  %s\n", styles.Key("Requirements:"))
	for _, req := range validation.DNS1123Requirements() {
		msg += fmt.Sprintf("    %s %s\n", styles.SymbolBullet, req)
	}
	msg += fmt.Sprintf("\n  %s\n", styles.Key("Valid examples:"))
	for _, example := range validation.DNS1123Examples() {
		msg += fmt.Sprintf("    %s %s\n", styles.SymbolBullet, styles.Code(example))
	}

	return fmt.Errorf("%s", msg)
}
