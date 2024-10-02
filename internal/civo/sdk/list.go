package sdk

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Client) GetInstances(ctx context.Context) ([]Instance, error) {
	return getAll[Instance](ctx, c)
}

func (c *Client) GetFirewalls(ctx context.Context) ([]Firewall, error) {
	return getAll[Firewall](ctx, c)
}

func (c *Client) GetVolumes(ctx context.Context) ([]Volume, error) {
	return getAll[Volume](ctx, c)
}

func (c *Client) GetKubernetesClusters(ctx context.Context) ([]KubernetesCluster, error) {
	return getAll[KubernetesCluster](ctx, c)
}

func (c *Client) GetNetworks(ctx context.Context) ([]Network, error) {
	return getAll[Network](ctx, c)
}

func (c *Client) GetObjectStores(ctx context.Context) ([]ObjectStore, error) {
	return getAll[ObjectStore](ctx, c)
}

func (c *Client) GetObjectStoreCredentials(ctx context.Context) ([]ObjectStoreCredential, error) {
	return getAll[ObjectStoreCredential](ctx, c)
}

func (c *Client) GetLoadBalancers(ctx context.Context) ([]LoadBalancer, error) {
	return getAll[LoadBalancer](ctx, c)
}

func (c *Client) GetSSHKeys(ctx context.Context) ([]SSHKey, error) {
	return getAll[SSHKey](ctx, c)
}

// PaginatedResponse is a helper struct to unmarshal paginated responses from
// the Civo API.
type PaginatedResponse[T Resource] struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Pages   int `json:"pages"`
	Items   []T `json:"items"`
}

// GetAll returns all resources of a given type.
func getAll[T Resource](ctx context.Context, c Civoer) ([]T, error) {
	var item T

	if item.IsSinglePaged() {
		return getSinglePage[T](ctx, c, item.GetAPIEndpoint())
	}

	return getPaginated[T](ctx, c, item.GetAPIEndpoint())
}

// getPaginated is a helper function to get results off an API endpoint that
// supports pagination using "page" and "perPage" query parameters.
func getPaginated[T Resource](ctx context.Context, c Civoer, endpoint string) ([]T, error) {
	var totalItems []T

	for page := 1; ; page++ {
		params := map[string]string{
			"page":     strconv.Itoa(page),
			"per_page": "100",
			"region":   c.GetRegion(),
		}

		var resp PaginatedResponse[T]
		err := c.Do(ctx, endpoint, "GET", &resp, params)
		if err != nil {
			return nil, fmt.Errorf("unable to get items: %w", err)
		}

		if resp.Page == page {
			totalItems = append(totalItems, resp.Items...)
		}

		// Break if we're on the last page or if the page number doesn't match
		// the expected page number (Civo returns page 1 if you request a page
		// that overflows).
		if resp.Page >= resp.Pages || resp.Page != page {
			break
		}
	}

	return totalItems, nil
}

// getSinglePage is a helper function to get results off an API endpoint that
// does not support pagination.
func getSinglePage[T Resource](ctx context.Context, c Civoer, endpoint string) ([]T, error) {
	var resp []T

	params := map[string]string{"region": c.GetRegion()}

	err := c.Do(ctx, endpoint, "GET", &resp, params)
	if err != nil {
		return nil, fmt.Errorf("unable to get items: %w", err)
	}

	return resp, nil
}
