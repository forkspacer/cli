package module

import (
	"context"

	"github.com/forkspacer/api-server/pkg/services/forkspacer"
	batchv1 "github.com/forkspacer/forkspacer/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service wraps api-server service and adds missing methods
type Service struct {
	apiService *forkspacer.ForkspacerModuleService
	client     client.Client // For Get and other operations
}

// NewService creates a new module service wrapper
func NewService() (*Service, error) {
	// Create api-server service (for Create, Delete, List, Update)
	apiService, err := forkspacer.NewForkspacerModuleService()
	if err != nil {
		return nil, err
	}

	// Create our own client for operations not supported by api-server
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
		apiService: apiService,
		client:     k8sClient,
	}, nil
}

// Create delegates to api-server service
func (s *Service) Create(ctx context.Context, moduleIn forkspacer.ModuleCreateIn) (*batchv1.Module, error) {
	return s.apiService.Create(ctx, moduleIn)
}

// Delete delegates to api-server service
func (s *Service) Delete(ctx context.Context, name string, namespace *string) error {
	return s.apiService.Delete(ctx, name, namespace)
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

// Update delegates to api-server service
func (s *Service) Update(ctx context.Context, updateIn forkspacer.ModuleUpdateIn) (*batchv1.Module, error) {
	return s.apiService.Update(ctx, updateIn)
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
) (*batchv1.Module, error) {
	module := &batchv1.Module{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.ModuleSpec{
			Source: batchv1.ModuleSource{
				ExistingHelmRelease: &batchv1.ModuleSourceExistingHelmReleaseRef{
					Name:      helmReleaseName,
					Namespace: helmReleaseNamespace,
				},
			},
			Workspace: batchv1.ModuleWorkspaceReference{
				Name:      workspaceName,
				Namespace: workspaceNamespace,
			},
			Hibernated: &hibernated,
		},
	}

	err := s.client.Create(ctx, module)
	return module, err
}
