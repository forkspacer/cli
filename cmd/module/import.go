package module

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/forkspacer/cli/cmd"
	"github.com/forkspacer/cli/pkg/module"
	"github.com/forkspacer/cli/pkg/styles"
)

type chartSourceType string

const (
	chartSourceGit    chartSourceType = "Git Repository"
	chartSourcePublic chartSourceType = "Public Chart Repository"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Interactively import an existing Helm release",
	Long: `Interactively import an existing Helm release to be managed by Forkspacer.

This command will guide you through:
  • Selecting a namespace
  • Choosing a Helm release from that namespace
  • Providing chart source information (Git or public chart repo)
  • Configuring workspace association

Examples:
  # Start interactive import
  forkspacer import`,
	RunE: runImport,
}

func init() {
	// Add to root command directly (forkspacer import)
	cmd.GetRootCmd().AddCommand(importCmd)
}

type importConfig struct {
	moduleName            string
	namespace             string
	helmRelease           string
	helmReleaseNamespace  string
	workspace             string
	workspaceNamespace    string
	chartSourceType       chartSourceType
	gitRepo               string
	gitPath               string
	gitRevision           string
	gitAuthSecret         string
	gitAuthSecretNS       string
	publicChartRepo       string
	publicChartName       string
	publicChartVersion    string
	chartRepoAuthSecret   string
	chartRepoAuthSecretNS string
	hibernated            bool
}

func runImport(c *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize Kubernetes client
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	config := &importConfig{
		namespace: cmd.GetNamespace(),
	}

	// Step 1: Select namespace
	fmt.Println()
	fmt.Println(styles.TitleStyle.Render("Import Helm Release"))
	fmt.Println()

	namespaces, err := getNamespaces(ctx, clientset)
	if err != nil {
		return fmt.Errorf("failed to get namespaces: %w", err)
	}

	nsOptions := make([]huh.Option[string], len(namespaces))
	for i, ns := range namespaces {
		nsOptions[i] = huh.NewOption(ns, ns)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select namespace").
				Options(nsOptions...).
				Value(&config.helmReleaseNamespace),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Step 2: Select Helm release
	releases, err := getHelmReleases(ctx, config.helmReleaseNamespace)
	if err != nil {
		return fmt.Errorf("failed to get Helm releases: %w", err)
	}

	if len(releases) == 0 {
		return fmt.Errorf("no Helm releases found in namespace %s", config.helmReleaseNamespace)
	}

	releaseOptions := make([]huh.Option[string], len(releases))
	for i, release := range releases {
		releaseOptions[i] = huh.NewOption(release, release)
	}

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Helm release").
				Options(releaseOptions...).
				Value(&config.helmRelease),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Step 3: Choose chart source type
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[chartSourceType]().
				Title("Select chart source type").
				Options(
					huh.NewOption("Git Repository", chartSourceGit),
					huh.NewOption("Public Chart Repository", chartSourcePublic),
				).
				Value(&config.chartSourceType),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Step 4: Get chart source details
	if config.chartSourceType == chartSourceGit {
		form = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Git Repository URL").
					Placeholder("https://github.com/org/repo").
					Value(&config.gitRepo).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("git repository URL is required")
						}
						return nil
					}),
				huh.NewInput().
					Title("Chart Path in Repository").
					Placeholder("charts/app or helm").
					Value(&config.gitPath).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("chart path is required")
						}
						return nil
					}),
				huh.NewInput().
					Title("Git Revision").
					Placeholder("main, master, v1.0.0, etc.").
					Value(&config.gitRevision).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("git revision is required")
						}
						return nil
					}),
				huh.NewInput().
					Title("Auth Secret Name (optional)").
					Description("Leave empty for public repos, provide secret name for private repos").
					Placeholder("github-token-secret").
					Value(&config.gitAuthSecret),
				huh.NewInput().
					Title("Auth Secret Namespace (optional)").
					Description("Namespace where the auth secret is located").
					Placeholder("default").
					Value(&config.gitAuthSecretNS),
			),
		)
	} else {
		form = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Chart Repository URL").
					Placeholder("https://charts.helm.sh/stable").
					Value(&config.publicChartRepo).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("chart repository URL is required")
						}
						return nil
					}),
				huh.NewInput().
					Title("Chart Name").
					Placeholder("nginx, postgresql, etc.").
					Value(&config.publicChartName).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("chart name is required")
						}
						return nil
					}),
				huh.NewInput().
					Title("Chart Version").
					Placeholder("1.0.0").
					Value(&config.publicChartVersion).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("chart version is required")
						}
						return nil
					}),
				huh.NewInput().
					Title("Auth Secret Name (optional)").
					Description("Leave empty for public repos, provide secret name for private repos").
					Placeholder("my-harbor-secret").
					Value(&config.chartRepoAuthSecret),
				huh.NewInput().
					Title("Auth Secret Namespace (optional)").
					Description("Namespace where the auth secret is located").
					Placeholder("default").
					Value(&config.chartRepoAuthSecretNS),
			),
		)
	}

	if err := form.Run(); err != nil {
		return err
	}

	// Step 5: Module configuration
	config.moduleName = config.helmRelease // Default to helm release name
	config.namespace = cmd.GetNamespace()

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Module Name").
				Placeholder(config.helmRelease).
				Value(&config.moduleName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("module name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Module Namespace").
				Placeholder(config.namespace).
				Value(&config.namespace),
			huh.NewInput().
				Title("Workspace Name").
				Placeholder("my-workspace").
				Value(&config.workspace).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("workspace name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Workspace Namespace").
				Placeholder(config.namespace).
				Value(&config.workspaceNamespace),
			huh.NewConfirm().
				Title("Import in hibernated state?").
				Value(&config.hibernated),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Default values
	if config.namespace == "" {
		config.namespace = cmd.GetNamespace()
	}
	if config.workspaceNamespace == "" {
		config.workspaceNamespace = config.namespace
	}

	// Step 6: Create the module
	return createModuleFromConfig(ctx, config)
}

func getNamespaces(ctx context.Context, clientset *kubernetes.Clientset) ([]string, error) {
	nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces := make([]string, 0, len(nsList.Items))
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}

func getHelmReleases(ctx context.Context, namespace string) ([]string, error) {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// List secrets with Helm label
	secrets, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "owner=helm",
	})
	if err != nil {
		return nil, err
	}

	// Extract unique release names
	releaseMap := make(map[string]bool)
	for _, secret := range secrets.Items {
		if name, ok := secret.Labels["name"]; ok {
			releaseMap[name] = true
		}
	}

	releases := make([]string, 0, len(releaseMap))
	for release := range releaseMap {
		releases = append(releases, release)
	}

	return releases, nil
}

func createModuleFromConfig(ctx context.Context, config *importConfig) error {
	fmt.Println()
	fmt.Println(styles.TitleStyle.Render(fmt.Sprintf("%s Creating module %s", styles.SymbolSparkles, config.moduleName)))
	fmt.Println()

	service, err := module.NewService()
	if err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	var moduleResource interface{}

	if config.chartSourceType == chartSourceGit {
		// Set default namespace for auth secret if not provided
		authSecretNS := config.gitAuthSecretNS
		if authSecretNS == "" && config.gitAuthSecret != "" {
			authSecretNS = config.namespace
		}

		moduleResource, err = service.CreateExistingHelmRelease(
			ctx,
			config.moduleName,
			config.namespace,
			config.helmRelease,
			config.helmReleaseNamespace,
			config.workspace,
			config.workspaceNamespace,
			config.hibernated,
			config.gitRepo,
			config.gitPath,
			config.gitRevision,
			config.gitAuthSecret,
			authSecretNS,
		)
	} else {
		// Set default namespace for auth secret if not provided
		authSecretNS := config.chartRepoAuthSecretNS
		if authSecretNS == "" && config.chartRepoAuthSecret != "" {
			authSecretNS = config.namespace
		}

		moduleResource, err = service.CreateExistingHelmReleaseWithChartRepo(
			ctx,
			config.moduleName,
			config.namespace,
			config.helmRelease,
			config.helmReleaseNamespace,
			config.workspace,
			config.workspaceNamespace,
			config.hibernated,
			config.publicChartRepo,
			config.publicChartName,
			config.publicChartVersion,
			config.chartRepoAuthSecret,
			authSecretNS,
		)
	}

	if err != nil {
		return fmt.Errorf("failed to create module: %w", err)
	}

	// Print success
	fmt.Println()
	fmt.Println(styles.SuccessStyle.Render("✓ Module created successfully"))
	fmt.Println()
	fmt.Println(styles.SubtitleStyle.Render("Next steps:"))
	fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code(fmt.Sprintf("forkspacer module get %s", config.moduleName)))
	fmt.Printf("  %s %s\n", styles.SymbolArrow, styles.Code(fmt.Sprintf("forkspacer workspace get %s", config.workspace)))
	fmt.Println()

	_ = moduleResource // Suppress unused variable warning
	return nil
}
