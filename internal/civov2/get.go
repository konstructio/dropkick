package civov2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
)

// GetInstance returns an instance by ID.
func (c *Client) GetInstance(ctx context.Context, id string) (*Instance, error) {
	return getByID[Instance](ctx, c, "/v2/instances", id)
}

// GetFirewall returns a firewall by ID.
func (c *Client) GetFirewall(ctx context.Context, id string) (*Firewall, error) {
	return getByID[Firewall](ctx, c, "/v2/firewalls", id)
}

// GetVolume returns a volume by ID.
func (c *Client) GetVolume(ctx context.Context, id string) (*Volume, error) {
	return getByID[Volume](ctx, c, "/v2/volumes", id)
}

// GetKubernetesCluster returns a Kubernetes cluster by ID.
func (c *Client) GetKubernetesCluster(ctx context.Context, id string) (*KubernetesCluster, error) {
	return getByID[KubernetesCluster](ctx, c, "/v2/kubernetes/clusters", id)
}

// GetNetwork returns a network by ID.
func (c *Client) GetNetwork(ctx context.Context, id string) (*Network, error) {
	return getByID[Network](ctx, c, "/v2/networks", id)
}

// GetObjectStore returns an object store by ID.
func (c *Client) GetObjectStore(ctx context.Context, id string) (*ObjectStore, error) {
	return getByID[ObjectStore](ctx, c, "/v2/objectstores", id)
}

// GetObjectStoreCredential returns an object store credential by ID.
func (c *Client) GetObjectStoreCredential(ctx context.Context, id string) (*ObjectStoreCredential, error) {
	return getByID[ObjectStoreCredential](ctx, c, "/v2/objectstore/credentials", id)
}

// GetSSHKey returns an SSH key by ID.
func (c *Client) GetSSHKey(ctx context.Context, id string) (*SSHKey, error) {
	return getByID[SSHKey](ctx, c, "/v2/sshkeys", id)
}

// getByID is a helper function to get an item by ID from the Civo API.
func getByID[T any](ctx context.Context, client *Client, endpoint, id string) (*T, error) {
	var output T

	params := map[string]string{"region": client.region}

	fullpath := path.Join(endpoint, id)
	if err := client.requester.doCivo(ctx, fullpath, http.MethodGet, &output, params); err != nil {
		if errors.Is(err, &HTTPError{Code: http.StatusNotFound}) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("unable to get item: %w", err)
	}

	return &output, nil
}
