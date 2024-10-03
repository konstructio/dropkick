package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/civo/sdk"
)

// NukeEverything deletes all resources associated with the Civo account it's
// targeting in the given Civo region.
func (c *Civo) NukeEverything(ctx context.Context) error {
	// The order in which these resources are deleted matter. We start with those
	// resources that have dependencies. In Civo, certain resources won't delete
	// their dependencies (for example, deleting an Instance that is on a Network
	// won't delete the Network because it could be shared with other Instances).
	// So we delete first those resources that could cascade delete other resources.

	// We start by deleting all Load Balancers, which tie to Firewalls and Kubernetes Clusters.
	if err := c.client.Each(ctx, sdk.LoadBalancer{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete load balancers: %w", err)
	}

	// Kubernetes clusters depend on volumes (PVCs), but we will delete the PVCs
	// after the clusters to clean them all.
	if err := c.client.Each(ctx, sdk.KubernetesCluster{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete Kubernetes clusters: %w", err)
	}

	// Then we delete the instances, which might also have volumes attached.
	if err := c.client.Each(ctx, sdk.Instance{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete instances: %w", err)
	}

	// Now we delete the volumes.
	if err := c.client.Each(ctx, sdk.Volume{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete volumes: %w", err)
	}

	// And we also delete the SSH keys, including those now orphaned by instances.
	if err := c.client.Each(ctx, sdk.SSHKey{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete SSH keys: %w", err)
	}

	// Then we delete object stores, which will leave their credentials orphaned.
	if err := c.client.Each(ctx, sdk.ObjectStore{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete object stores: %w", err)
	}

	// Now we delete the object store credentials.
	if err := c.client.Each(ctx, sdk.ObjectStoreCredential{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete object store credentials: %w", err)
	}

	// Firewalls are deleted next, since we need them to be gone before deleting networks.
	if err := c.client.Each(ctx, sdk.Firewall{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete firewalls: %w", err)
	}

	// And finally, we delete the networks.
	if err := c.client.Each(ctx, sdk.Network{}, c.deleteIterator(ctx)); err != nil {
		return fmt.Errorf("unable to delete networks: %w", err)
	}

	return nil
}
