package workspace

import (
	"context"

	batchv1 "github.com/forkspacer/forkspacer/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service provides operations for managing workspaces
type Service struct {
	client client.Client
}

// WorkspaceCreateInput defines the input for creating a workspace
type WorkspaceCreateInput struct {
	Name            string
	Namespace       string
	Hibernated      bool
	ConnectionType  string
	AutoHibernation *AutoHibernationInput
	From            *FromWorkspaceInput
}

// AutoHibernationInput defines auto-hibernation configuration
type AutoHibernationInput struct {
	Enabled      bool
	Schedule     string
	WakeSchedule *string
}

// FromWorkspaceInput defines workspace forking configuration
type FromWorkspaceInput struct {
	Name      string
	Namespace string
}

// NewService creates a new workspace service
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

// Create creates a new workspace
func (s *Service) Create(ctx context.Context, input WorkspaceCreateInput) (*batchv1.Workspace, error) {
	workspace := &batchv1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.Name,
			Namespace: input.Namespace,
		},
		Spec: batchv1.WorkspaceSpec{
			Type:       batchv1.WorkspaceTypeKubernetes,
			Hibernated: input.Hibernated,
			Connection: batchv1.WorkspaceConnection{
				Type: batchv1.WorkspaceConnectionType(input.ConnectionType),
			},
		},
	}

	// Add auto-hibernation if specified
	if input.AutoHibernation != nil {
		workspace.Spec.AutoHibernation = &batchv1.WorkspaceAutoHibernation{
			Enabled:  input.AutoHibernation.Enabled,
			Schedule: input.AutoHibernation.Schedule,
		}
		if input.AutoHibernation.WakeSchedule != nil {
			workspace.Spec.AutoHibernation.WakeSchedule = input.AutoHibernation.WakeSchedule
		}
	}

	// Add fork reference if specified
	if input.From != nil {
		workspace.Spec.From = &batchv1.WorkspaceFromReference{
			Name:      input.From.Name,
			Namespace: input.From.Namespace,
		}
	}

	err := s.client.Create(ctx, workspace)
	return workspace, err
}

// Delete deletes a workspace
func (s *Service) Delete(ctx context.Context, name string, namespace *string) error {
	ns := "default"
	if namespace != nil {
		ns = *namespace
	}

	workspace := &batchv1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	return s.client.Delete(ctx, workspace)
}

// List workspaces with optional namespace filtering
func (s *Service) List(ctx context.Context, namespace string) (*batchv1.WorkspaceList, error) {
	workspaces := &batchv1.WorkspaceList{}

	var opts []client.ListOption
	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	err := s.client.List(ctx, workspaces, opts...)
	return workspaces, err
}

// Get fetches a single workspace
func (s *Service) Get(ctx context.Context, name, namespace string) (*batchv1.Workspace, error) {
	workspace := &batchv1.Workspace{}
	err := s.client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, workspace)
	return workspace, err
}

// SetHibernation updates the hibernation state of a workspace
func (s *Service) SetHibernation(ctx context.Context, name, namespace string, hibernated bool) (*batchv1.Workspace, error) {
	workspace, err := s.Get(ctx, name, namespace)
	if err != nil {
		return nil, err
	}

	workspace.Spec.Hibernated = hibernated

	err = s.client.Update(ctx, workspace)
	return workspace, err
}
