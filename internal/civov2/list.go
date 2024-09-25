package civov2

import (
	"context"
	"fmt"
	"strconv"
)

// GetAllInstances returns all instances in the Civo account.
func (c *Client) GetAllInstances(ctx context.Context) ([]Instance, error) {
	// currently, instances endpoint is paginated
	return getPaginated[Instance](ctx, c, "/v2/instances")
}

// GetAllVolumes returns all volumes in the Civo account.
func (c *Client) GetAllFirewalls(ctx context.Context) ([]Firewall, error) {
	// currently, firewalls endpoint is NOT paginated
	return getSinglePage[Firewall](ctx, c, "/v2/firewalls")
}

// GetAllVolumes returns all volumes in the Civo account.
func (c *Client) GetAllVolumes(ctx context.Context) ([]Volume, error) {
	// currently, volumes endpoint is NOT paginated
	return getSinglePage[Volume](ctx, c, "/v2/volumes")
}

// GetAllKubernetesClusters returns all Kubernetes clusters in the Civo account.
func (c *Client) GetAllKubernetesClusters(ctx context.Context) ([]KubernetesCluster, error) {
	// currently, kubernetes clusters endpoint is paginated
	return getPaginated[KubernetesCluster](ctx, c, "/v2/kubernetes/clusters")
}

// GetAllNetworks returns all networks in the Civo account.
func (c *Client) GetAllNetworks(ctx context.Context) ([]Network, error) {
	// currently, networks endpoint is NOT paginated
	return getSinglePage[Network](ctx, c, "/v2/networks")
}

// GetAllObjectStores returns all object stores in the Civo account.
func (c *Client) GetAllObjectStores(ctx context.Context) ([]ObjectStore, error) {
	// currently, object stores endpoint is paginated
	return getPaginated[ObjectStore](ctx, c, "/v2/objectstores")
}

// GetAllObjectStoreCredentials returns all object store credentials in the Civo account.
func (c *Client) GetAllObjectStoreCredentials(ctx context.Context) ([]ObjectStoreCredential, error) {
	// currently, object store credentials endpoint is paginated
	return getPaginated[ObjectStoreCredential](ctx, c, "/v2/objectstore/credentials")
}

// GetAllSSHKeys returns all SSH keys in the Civo account.
func (c *Client) GetAllSSHKeys(ctx context.Context) ([]SSHKey, error) {
	// currently, SSH keys endpoint is NOT paginated
	return getSinglePage[SSHKey](ctx, c, "/v2/sshkeys")
}

// getPaginated is a helper function to get results off an API endpoint that
// supports pagination using "page" and "perPage" query parameters.
func getPaginated[T any](ctx context.Context, c *Client, endpoint string) ([]T, error) {
	var totalItems []T

	for page := 1; ; page++ {
		params := map[string]string{
			"page":    strconv.Itoa(page),
			"perPage": "100",
			"region":  c.region,
		}

		var resp struct {
			Page    int `json:"page"`
			PerPage int `json:"per_page"`
			Pages   int `json:"pages"`
			Items   []T `json:"items"`
		}

		err := c.requester.doCivo(ctx, endpoint, "GET", &resp, params)
		if err != nil {
			return nil, fmt.Errorf("unable to get items: %w", err)
		}

		totalItems = append(totalItems, resp.Items...)

		if resp.Page >= resp.Pages {
			break
		}
	}

	return totalItems, nil
}

// getSinglePage is a helper function to get results off an API endpoint that
// does not support pagination.
func getSinglePage[T any](ctx context.Context, c *Client, endpoint string) ([]T, error) {
	var resp []T

	err := c.requester.doCivo(ctx, endpoint, "GET", &resp, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get items: %w", err)
	}

	return resp, nil
}
