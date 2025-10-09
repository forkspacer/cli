package k8s

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "github.com/forkspacer/forkspacer/api/v1"
)

// Client wraps the Kubernetes client with helper methods
type Client struct {
	client.Client
	Scheme  *runtime.Scheme
	Context string // Current kubectl context name
}

// NewClient creates a new Kubernetes client with Forkspacer CRDs registered
func NewClient() (*Client, error) {
	// Create scheme and register types
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add client-go scheme: %w", err)
	}
	if err := batchv1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add Forkspacer CRDs: %w", err)
	}

	// Get Kubernetes config (supports kubeconfig and in-cluster)
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes config: %w", err)
	}

	// Create controller-runtime client
	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Get current context name
	contextName := "unknown"
	if config := ctrl.GetConfigOrDie(); config != nil {
		// Try to extract context from kubeconfig
		contextName = "current"
	}

	return &Client{
		Client:  k8sClient,
		Scheme:  scheme,
		Context: contextName,
	}, nil
}

// CheckOperatorInstalled verifies that Forkspacer operator is installed
func (c *Client) CheckOperatorInstalled(ctx context.Context) error {
	// Try to list workspaces - if CRD doesn't exist, this will fail
	workspaces := &batchv1.WorkspaceList{}
	if err := c.List(ctx, workspaces, client.Limit(1)); err != nil {
		return fmt.Errorf("Forkspacer operator not found: %w", err)
	}
	return nil
}
