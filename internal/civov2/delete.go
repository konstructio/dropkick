package civov2

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

// DeleteInstance deletes an instance.
func (client *Client) DeleteInstance(ctx context.Context, id string) error {
	return delete(client, ctx, "/v2/instances", id)
}

// DeleteFirewall deletes a firewall.
func (client *Client) DeleteFirewall(ctx context.Context, id string) error {
	return delete(client, ctx, "/v2/firewalls", id)
}

// DeleteVolume deletes a volume.
func (client *Client) DeleteVolume(ctx context.Context, id string) error {
	return delete(client, ctx, "/v2/volumes", id)
}

// DeleteKubernetesCluster deletes a Kubernetes cluster.
func (client *Client) DeleteKubernetesCluster(ctx context.Context, id string) error {
	return delete(client, ctx, "/v2/kubernetes/clusters", id)
}

// DeleteNetwork deletes a network.
func (client *Client) DeleteNetwork(ctx context.Context, id string) error {
	return delete(client, ctx, "/v2/networks", id)
}

// DeleteObjectStore deletes an object store.
func (client *Client) DeleteObjectStore(ctx context.Context, id string) error {
	return delete(client, ctx, "/v2/objectstores", id)
}

// DeleteObjectStoreCredential deletes an object store credential.
func (client *Client) DeleteObjectStoreCredential(ctx context.Context, id string) error {
	return delete(client, ctx, "/v2/objectstore/credentials", id)
}

// DeleteSSHKey deletes an SSH key.
func (client *Client) DeleteSSHKey(ctx context.Context, id string) error {
	return delete(client, ctx, "/v2/sshkeys", id)
}

// delete is a helper function to delete an item via a HTTP DELETE request
// to the Civo API.
func delete(client *Client, ctx context.Context, endpoint, id string) error {
	params := map[string]string{
		"region": client.region,
	}

	var output struct {
		ID     string `json:"id"`
		Result string `json:"result"`
	}

	fullpath := path.Join(endpoint, id)
	if err := client.requester.doCivo(ctx, fullpath, http.MethodDelete, nil, &output, params); err != nil {
		return fmt.Errorf("unable to delete item: %w", err)
	}

	return nil
}
