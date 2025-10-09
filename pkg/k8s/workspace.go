package k8s

import (
	"context"
	"fmt"

	batchv1 "github.com/forkspacer/forkspacer/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListWorkspaces lists all workspaces in the specified namespace
func (c *Client) ListWorkspaces(ctx context.Context, namespace string) (*batchv1.WorkspaceList, error) {
	workspaces := &batchv1.WorkspaceList{}
	opts := []client.ListOption{}

	if namespace != "" {
		opts = append(opts, client.InNamespace(namespace))
	}

	if err := c.List(ctx, workspaces, opts...); err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", err)
	}

	return workspaces, nil
}

// GetWorkspace retrieves a specific workspace
func (c *Client) GetWorkspace(ctx context.Context, name, namespace string) (*batchv1.Workspace, error) {
	workspace := &batchv1.Workspace{}
	key := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	if err := c.Get(ctx, key, workspace); err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	return workspace, nil
}

// CreateWorkspace creates a new workspace
func (c *Client) CreateWorkspace(ctx context.Context, workspace *batchv1.Workspace) error {
	if err := c.Create(ctx, workspace); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}
	return nil
}

// DeleteWorkspace deletes a workspace
func (c *Client) DeleteWorkspace(ctx context.Context, name, namespace string) error {
	workspace := &batchv1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	if err := c.Delete(ctx, workspace); err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	return nil
}

// PatchWorkspace patches a workspace
func (c *Client) PatchWorkspace(ctx context.Context, workspace *batchv1.Workspace, patch client.Patch) error {
	if err := c.Patch(ctx, workspace, patch); err != nil {
		return fmt.Errorf("failed to patch workspace: %w", err)
	}
	return nil
}

// WorkspaceExists checks if a workspace exists
func (c *Client) WorkspaceExists(ctx context.Context, name, namespace string) (bool, error) {
	_, err := c.GetWorkspace(ctx, name, namespace)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
