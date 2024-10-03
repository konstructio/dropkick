package sdk

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

// DeleteAll removes all resources of a given type from the Civo API.
func (c *Client) DeleteAll(ctx context.Context, v APIResource, conditionFunc func(APIResource) bool) error {
	switch r := v.(type) {
	case Instance:
		return nuke(ctx, c, func(i Instance) bool { return conditionFunc(i) })
	case Firewall:
		return nuke(ctx, c, func(f Firewall) bool { return conditionFunc(f) })
	case Volume:
		return nuke(ctx, c, func(v Volume) bool { return conditionFunc(v) })
	case KubernetesCluster:
		return nuke(ctx, c, func(k KubernetesCluster) bool { return conditionFunc(k) })
	case Network:
		return nuke(ctx, c, func(n Network) bool { return conditionFunc(n) })
	case ObjectStore:
		return nuke(ctx, c, func(o ObjectStore) bool { return conditionFunc(o) })
	case ObjectStoreCredential:
		return nuke(ctx, c, func(o ObjectStoreCredential) bool { return conditionFunc(o) })
	case LoadBalancer:
		return nuke(ctx, c, func(l LoadBalancer) bool { return conditionFunc(l) })
	case SSHKey:
		return nuke(ctx, c, func(s SSHKey) bool { return conditionFunc(s) })
	default:
		return fmt.Errorf("unsupported resource type: %T", r)
	}
}

// nuke removes all resources of a given type from the Civo API.
func nuke[T Resource](ctx context.Context, c Civoer, filterFunc func(T) bool) error {
	var res T

	resources, err := getAll[T](ctx, c)
	if err != nil {
		return fmt.Errorf("unable to list %s: %w", res.GetResourceType(), err)
	}

	for _, r := range resources {
		if filterFunc(r) {
			if err := deleteResource(ctx, c, r); err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete removes a resource from the Civo API.
func (c *Client) Delete(ctx context.Context, resource APIResource) error {
	switch r := resource.(type) {
	case Instance:
		return deleteResource(ctx, c, r)
	case Firewall:
		return deleteResource(ctx, c, r)
	case Volume:
		return deleteResource(ctx, c, r)
	case KubernetesCluster:
		return deleteResource(ctx, c, r)
	case Network:
		return deleteResource(ctx, c, r)
	case ObjectStore:
		return deleteResource(ctx, c, r)
	case ObjectStoreCredential:
		return deleteResource(ctx, c, r)
	case LoadBalancer:
		return deleteResource(ctx, c, r)
	case SSHKey:
		return deleteResource(ctx, c, r)
	default:
		return fmt.Errorf("unsupported resource type: %T", r)
	}
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
