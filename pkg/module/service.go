package module

import (
	"context"

	batchv1 "github.com/forkspacer/forkspacer/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service provides operations for managing modules
type Service struct {
	client client.Client
}

// NewService creates a new module service
func NewService() (*Service, error) {
	restConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := batchv1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	k8sClient, err := client.New(restConfig, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return &Service{
		client: k8sClient,
	}, nil
}

// Delete deletes a module
func (s *Service) Delete(ctx context.Context, name string, namespace *string) error {
	ns := "default"
	if namespace != nil {
		ns = *namespace
	}

	module := &batchv1.Module{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	return s.client.Delete(ctx, module)
}

// List modules with optional namespace filtering
func (s *Service) List(ctx context.Context, namespace string) (*batchv1.ModuleList, error) {
	modules := &batchv1.ModuleList{}

	var opts []client.ListOption
	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	err := s.client.List(ctx, modules, opts...)
	return modules, err
}

// Get fetches a single module
func (s *Service) Get(ctx context.Context, name, namespace string) (*batchv1.Module, error) {
	module := &batchv1.Module{}
	err := s.client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, module)
	return module, err
}

// CreateExistingHelmRelease creates a module that imports an existing Helm release
func (s *Service) CreateExistingHelmRelease(
	ctx context.Context,
	name string,
	namespace string,
	helmReleaseName string,
	helmReleaseNamespace string,
	workspaceName string,
	workspaceNamespace string,
	hibernated bool,
	chartSourceGitRepo string,
	chartSourceGitPath string,
	chartSourceGitRevision string,
	authSecretName string,
	authSecretNamespace string,
) (*batchv1.Module, error) {
	chartGit := &batchv1.ModuleSpecHelmChartGit{
		Repo:     chartSourceGitRepo,
		Path:     chartSourceGitPath,
		Revision: chartSourceGitRevision,
	}

	// Add auth if secret is provided
	if authSecretName != "" {
		chartGit.Auth = &batchv1.ModuleSpecHelmChartGitAuth{
			HTTPSSecretRef: &batchv1.ModuleSpecHelmChartGitAuthSecret{
				Name:      authSecretName,
				Namespace: authSecretNamespace,
			},
		}
	}

	module := &batchv1.Module{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.ModuleSpec{
			Helm: &batchv1.ModuleSpecHelm{
				ExistingRelease: &batchv1.ModuleSpecHelmExistingRelease{
					Name:      helmReleaseName,
					Namespace: helmReleaseNamespace,
				},
				Chart: batchv1.ModuleSpecHelmChart{
					Git: chartGit,
				},
			},
			Workspace: batchv1.ModuleWorkspaceReference{
				Name:      workspaceName,
				Namespace: workspaceNamespace,
			},
			Hibernated: hibernated,
		},
	}

	err := s.client.Create(ctx, module)
	return module, err
}

// CreateExistingHelmReleaseWithChartRepo creates a module that imports an existing Helm release using a chart repository
func (s *Service) CreateExistingHelmReleaseWithChartRepo(
	ctx context.Context,
	name string,
	namespace string,
	helmReleaseName string,
	helmReleaseNamespace string,
	workspaceName string,
	workspaceNamespace string,
	hibernated bool,
	chartRepoURL string,
	chartName string,
	chartVersion string,
	authSecretName string,
	authSecretNamespace string,
) (*batchv1.Module, error) {
	chartRepo := &batchv1.ModuleSpecHelmChartRepo{
		URL:     chartRepoURL,
		Chart:   chartName,
		Version: nil,
	}

	if chartVersion != "" {
		chartRepo.Version = &chartVersion
	}

	// Add auth if secret is provided
	if authSecretName != "" {
		chartRepo.Auth = &batchv1.ModuleSpecHelmChartRepoAuth{
			Name:      authSecretName,
			Namespace: authSecretNamespace,
		}
	}

	module := &batchv1.Module{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.ModuleSpec{
			Helm: &batchv1.ModuleSpecHelm{
				ExistingRelease: &batchv1.ModuleSpecHelmExistingRelease{
					Name:      helmReleaseName,
					Namespace: helmReleaseNamespace,
				},
				Chart: batchv1.ModuleSpecHelmChart{
					Repo: chartRepo,
				},
			},
			Workspace: batchv1.ModuleWorkspaceReference{
				Name:      workspaceName,
				Namespace: workspaceNamespace,
			},
			Hibernated: hibernated,
		},
	}

	err := s.client.Create(ctx, module)
	return module, err
}
