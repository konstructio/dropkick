package civov2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
)

// GetInstance returns an instance by ID.
func GetInstance(ctx context.Context, client *Client, id string) (*Instance, error) {
	return getByID[Instance](client, ctx, "/v2/instances", id)
}

// GetFirewall returns a firewall by ID.
func GetFirewall(ctx context.Context, client *Client, id string) (*Firewall, error) {
	return getByID[Firewall](client, ctx, "/v2/firewalls", id)
}

// GetVolume returns a volume by ID.
func GetVolume(ctx context.Context, client *Client, id string) (*Volume, error) {
	return getByID[Volume](client, ctx, "/v2/volumes", id)
}

// GetKubernetesCluster returns a Kubernetes cluster by ID.
func GetKubernetesCluster(ctx context.Context, client *Client, id string) (*KubernetesCluster, error) {
	return getByID[KubernetesCluster](client, ctx, "/v2/kubernetes/clusters", id)
}

// GetNetwork returns a network by ID.
func GetNetwork(ctx context.Context, client *Client, id string) (*Network, error) {
	return getByID[Network](client, ctx, "/v2/networks", id)
}

// GetObjectStore returns an object store by ID.
func GetObjectStore(ctx context.Context, client *Client, id string) (*ObjectStore, error) {
	return getByID[ObjectStore](client, ctx, "/v2/objectstores", id)
}

// GetObjectStoreCredential returns an object store credential by ID.
func GetObjectStoreCredential(ctx context.Context, client *Client, id string) (*ObjectStoreCredential, error) {
	return getByID[ObjectStoreCredential](client, ctx, "/v2/objectstore/credentials", id)
}

// GetSSHKey returns an SSH key by ID.
func GetSSHKey(ctx context.Context, client *Client, id string) (*SSHKey, error) {
	return getByID[SSHKey](client, ctx, "/v2/sshkeys", id)
}

// getByID is a helper function to get an item by ID from the Civo API.
func getByID[T any](client *Client, ctx context.Context, endpoint, id string) (*T, error) {
	var output T

	params := map[string]string{"region": client.region}

	fullpath := path.Join(endpoint, id)
	if err := client.requester.doCivo(ctx, fullpath, http.MethodGet, nil, &output, params); err != nil {
		if errors.Is(err, &HTTPError{Code: http.StatusNotFound}) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("unable to get item: %w", err)
	}

	return &output, nil
}
