package workspace

import (
	"context"

	"github.com/forkspacer/api-server/pkg/services/forkspacer"
	"github.com/forkspacer/api-server/pkg/utils"
	batchv1 "github.com/forkspacer/forkspacer/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service wraps api-server service and adds missing Get() method
// TODO: Remove this wrapper when api-server adds Get() method
type Service struct {
	apiService *forkspacer.ForkspacerWorkspaceService
	client     client.Client // Only for Get operation
}

// NewService creates a new workspace service wrapper
func NewService() (*Service, error) {
	// Create api-server service (for Create, Delete, List, Update)
	apiService, err := forkspacer.NewForkspacerWorkspaceService()
	if err != nil {
		return nil, err
	}

	// Create our own client for Get operation
	// (api-server's client field is not exported)
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
func (s *Service) Create(ctx context.Context, workspaceIn forkspacer.WorkspaceCreateIn) (*batchv1.Workspace, error) {
	return s.apiService.Create(ctx, workspaceIn)
}

// Delete delegates to api-server service
func (s *Service) Delete(ctx context.Context, name string, namespace *string) error {
	return s.apiService.Delete(ctx, name, namespace)
}

// List workspaces with optional namespace filtering
// TODO: Add namespace parameter to api-server's List() method
func (s *Service) List(ctx context.Context, namespace string) (*batchv1.WorkspaceList, error) {
	// Use our client for namespace filtering (api-server doesn't support this yet)
	workspaces := &batchv1.WorkspaceList{}

	var opts []client.ListOption
	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	err := s.client.List(ctx, workspaces, opts...)
	return workspaces, err
}

// Update delegates to api-server service
func (s *Service) Update(ctx context.Context, updateIn forkspacer.WorkspaceUpdateIn) (*batchv1.Workspace, error) {
	return s.apiService.Update(ctx, updateIn)
}

// Get fetches a single workspace
// TODO: Remove when api-server adds Get() method
func (s *Service) Get(ctx context.Context, name, namespace string) (*batchv1.Workspace, error) {
	workspace := &batchv1.Workspace{}
	err := s.client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, workspace)
	return workspace, err
}

// SetHibernation is a helper to set hibernation state
// Uses api-server's Update with retry logic
func (s *Service) SetHibernation(ctx context.Context, name, namespace string, hibernated bool) (*batchv1.Workspace, error) {
	return s.apiService.Update(ctx, forkspacer.WorkspaceUpdateIn{
		Name:       name,
		Namespace:  &namespace,
		Hibernated: utils.ToPtr(hibernated),
	})
}
