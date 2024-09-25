package civov2

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

// DeleteInstance deletes an instance.
func (c *Client) DeleteInstance(ctx context.Context, id string) error {
	return deleteItem(ctx, c, "/v2/instances", id)
}

// DeleteFirewall deletes a firewall.
func (c *Client) DeleteFirewall(ctx context.Context, id string) error {
	return deleteItem(ctx, c, "/v2/firewalls", id)
}

// DeleteVolume deletes a volume.
func (c *Client) DeleteVolume(ctx context.Context, id string) error {
	return deleteItem(ctx, c, "/v2/volumes", id)
}

// DeleteKubernetesCluster deletes a Kubernetes cluster.
func (c *Client) DeleteKubernetesCluster(ctx context.Context, id string) error {
	return deleteItem(ctx, c, "/v2/kubernetes/clusters", id)
}

// DeleteNetwork deletes a network.
func (c *Client) DeleteNetwork(ctx context.Context, id string) error {
	return deleteItem(ctx, c, "/v2/networks", id)
}

// DeleteObjectStore deletes an object store.
func (c *Client) DeleteObjectStore(ctx context.Context, id string) error {
	return deleteItem(ctx, c, "/v2/objectstores", id)
}

// DeleteObjectStoreCredential deletes an object store credential.
func (c *Client) DeleteObjectStoreCredential(ctx context.Context, id string) error {
	return deleteItem(ctx, c, "/v2/objectstore/credentials", id)
}

// DeleteSSHKey deletes an SSH key.
func (c *Client) DeleteSSHKey(ctx context.Context, id string) error {
	return deleteItem(ctx, c, "/v2/sshkeys", id)
}

// deleteItem is a helper function to deleteItem an item via a HTTP DELETE request
// to the Civo API.
func deleteItem(ctx context.Context, client *Client, endpoint, id string) error {
	params := map[string]string{
		"region": client.region,
	}

	var output struct {
		ID     string `json:"id"`
		Result string `json:"result"`
	}

	fullpath := path.Join(endpoint, id)
	if err := client.requester.doCivo(ctx, fullpath, http.MethodDelete, &output, params); err != nil {
		return fmt.Errorf("unable to delete item: %w", err)
	}

	return nil
}
