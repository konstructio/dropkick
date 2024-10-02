package sdk

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

// DeleteInstance removes an instance from the Civo API.
func (c *Client) DeleteInstance(ctx context.Context, instance Instance) error {
	return deleteResource(ctx, c, instance)
}

// DeleteFirewall removes a firewall from the Civo API.
func (c *Client) DeleteFirewall(ctx context.Context, firewall Firewall) error {
	return deleteResource(ctx, c, firewall)
}

// DeleteVolume removes a volume from the Civo API.
func (c *Client) DeleteVolume(ctx context.Context, volume Volume) error {
	return deleteResource(ctx, c, volume)
}

// DeleteKubernetesCluster removes a Kubernetes cluster from the Civo API.
func (c *Client) DeleteKubernetesCluster(ctx context.Context, cluster KubernetesCluster) error {
	return deleteResource(ctx, c, cluster)
}

// DeleteNetwork removes a network from the Civo API.
func (c *Client) DeleteNetwork(ctx context.Context, network Network) error {
	return deleteResource(ctx, c, network)
}

// DeleteObjectStore removes an object store from the Civo API.
func (c *Client) DeleteObjectStore(ctx context.Context, objstore ObjectStore) error {
	return deleteResource(ctx, c, objstore)
}

// DeleteObjectStoreCredential removes an object store credential from the Civo API.
func (c *Client) DeleteObjectStoreCredential(ctx context.Context, objstorecred ObjectStoreCredential) error {
	return deleteResource(ctx, c, objstorecred)
}

// DeleteLoadBalancer removes a load balancer from the Civo API.
func (c *Client) DeleteLoadBalancer(ctx context.Context, lb LoadBalancer) error {
	return deleteResource(ctx, c, lb)
}

// DeleteSSHKey removes an SSH key from the Civo API.
func (c *Client) DeleteSSHKey(ctx context.Context, sshkey SSHKey) error {
	return deleteResource(ctx, c, sshkey)
}

// Delete removes a resource from the Civo API. The resource must have
// a non-empty ID.
func deleteResource[T Resource](ctx context.Context, c Civoer, resource T) error {
	if resource.GetID() == "" {
		return fmt.Errorf("the ID field in the resource %s is empty", resource.GetResourceType())
	}

	params := map[string]string{"region": c.GetRegion()}

	fullpath := path.Join(resource.GetAPIEndpoint(), resource.GetID())
	if err := c.Do(ctx, fullpath, http.MethodDelete, nil, params); err != nil {
		return fmt.Errorf("unable to delete item: %w", err)
	}

	return nil
}
