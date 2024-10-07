package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/civo/sdk"
)

// NukeOrphanedResources deletes all subresources not in use by a compute
// instance in the Civo account targeted. It returns an error if the deletion
// process encounters any issues. The resources targeted by this function are:
// - Load Balancers
// - Volumes
// - Object store credentials
// - SSH keys
// - Networks
// - Firewalls
func (c *Civo) NukeOrphanedResources(ctx context.Context) error {
	// fetch all nodes first, we'll need them to check for orphaned resources
	c.logger.Infof("fetching all instances")
	nodes, err := c.client.GetInstances(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch instances: %w", err)
	}

	// fetch also all volumes to check for networks connected to them
	c.logger.Infof("fetching all volumes")
	volumes, err := c.client.GetVolumes(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch volumes: %w", err)
	}

	// fetch orphaned load balancers
	orphanedLBs, err := c.getOrphanedLoadBalancers(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned load balancers: %w", err)
	}

	if err := nukeSlice(ctx, c, orphanedLBs); err != nil {
		return fmt.Errorf("unable to delete orphaned load balancers: %w", err)
	}

	// fetch orphaned volumes
	orphanedVolumes := c.getOrphanedVolumes(volumes)
	if err := nukeSlice(ctx, c, orphanedVolumes); err != nil {
		return fmt.Errorf("unable to delete orphaned volumes: %w", err)
	}

	// fetch orphaned object store credentials
	orphanedObjectStoreCredentials, err := c.getOrphanedObjectStoreCredentials(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned object store credentials: %w", err)
	}

	if err := nukeSlice(ctx, c, orphanedObjectStoreCredentials); err != nil {
		return fmt.Errorf("unable to delete orphaned object store credentials: %w", err)
	}

	// fetch orphaned SSH keys
	orphanedSSHKeys, err := c.getOrphanedSSHKeys(ctx, nodes)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned SSH keys: %w", err)
	}

	if err := nukeSlice(ctx, c, orphanedSSHKeys); err != nil {
		return fmt.Errorf("unable to delete orphaned SSH keys: %w", err)
	}

	// fetch orphaned networks
	orphanedNetworks, err := c.getOrphanedNetworks(ctx, nodes, volumes)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned networks: %w", err)
	}

	if err := nukeSlice(ctx, c, orphanedNetworks); err != nil {
		return fmt.Errorf("unable to delete orphaned networks: %w", err)
	}

	// fetch orphaned firewalls
	orphanedFirewalls, err := c.getOrphanedFirewalls(ctx)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned firewalls: %w", err)
	}

	if err := nukeSlice(ctx, c, orphanedFirewalls); err != nil {
		return fmt.Errorf("unable to delete orphaned firewalls: %w", err)
	}

	return nil
}

// nukeSlice deletes all resources in the provided slice. It returns an error if
// the deletion process encounters any issues.
func nukeSlice[T sdk.Resource](ctx context.Context, c *Civo, resources []T) error {
	for _, resource := range resources {
		if err := c.deleteIterator(ctx)(resource); err != nil {
			return err
		}
	}

	return nil
}

// getOrphanedObjectStoreCredentials fetches all object store then the object
// store credentials and compares them against each other. If a credential is
// not used by any store, it is considered orphaned. It returns an error if the
// fetching process encounters any issues.
func (c *Civo) getOrphanedObjectStoreCredentials(ctx context.Context) ([]sdk.ObjectStoreCredential, error) {
	c.logger.Infof("listing object stores")

	objectStores, err := c.client.GetObjectStores(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list object stores: %w", err)
	}

	credentials, err := c.client.GetObjectStoreCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list object store credentials: %w", err)
	}

	c.logger.Infof("found %d object store credentials", len(credentials))

	// iterate over all object stores and check if they are associated with any credentials
	orphanedCredentials := make([]sdk.ObjectStoreCredential, 0)
	for _, credential := range credentials {
		var found bool

		// iterate through the object stores finding if they use the current credential
		for _, objectStore := range objectStores {
			// on a GET request for object stores, only
			if objectStore.Credentials.ID == credential.CredentialID {
				c.logger.Warnf("skipping object store credential %q: it is associated with the object store with ID %q", credential.Name, objectStore.ID)
				found = true
				break
			}
		}

		if !found {
			c.logger.Infof("found orphaned object store credential %q - ID: %q", credential.Name, credential.ID)
			orphanedCredentials = append(orphanedCredentials, credential)
		}
	}

	c.logger.Infof("got %d object store credentials of which %d are orphan", len(credentials), len(orphanedCredentials))
	return orphanedCredentials, nil
}

func (c *Civo) getOrphanedLoadBalancers(ctx context.Context) ([]sdk.LoadBalancer, error) {
	c.logger.Infof("listing load balancers")

	lbs, err := c.client.GetLoadBalancers(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list load balancers: %w", err)
	}

	// iterate over all load balancers and check if they are associated with any nodes
	orphanedLBs := make([]sdk.LoadBalancer, 0, len(lbs))
	for _, lb := range lbs {
		if lb.ClusterID != "" {
			c.logger.Warnf("skipping load balancer %q: it is associated with the cluster with ID %q", lb.Name, lb.ClusterID)
			continue
		}

		if lb.FirewallID != "" {
			c.logger.Warnf("skipping load balancer %q: it is associated with the firewall with ID %q", lb.Name, lb.FirewallID)
			continue
		}

		c.logger.Infof("found orphaned load balancer %q - ID: %q", lb.Name, lb.ID)
		orphanedLBs = append(orphanedLBs, lb)
	}

	c.logger.Infof("found %d load balancers, %d of which are orphaned", len(lbs), len(orphanedLBs))
	return orphanedLBs, nil
}

// getOrphanedVolumes fetches all volumes that are not attached to any node
// instance instead of relying if they are referenced by a node instance. It
// returns an error if the fetching process encounters any issues.
func (c *Civo) getOrphanedVolumes(volumes []sdk.Volume) []sdk.Volume {
	newVolumeList := make([]sdk.Volume, 0, len(volumes))

	for _, volume := range volumes {
		if volume.Status == "attached" {
			c.logger.Warnf("skipping volume %q: it is attached to the node instance with ID %q", volume.Name, volume.InstanceID)
			continue
		}

		c.logger.Infof("found orphaned volume (not attached) %q - ID: %q", volume.Name, volume.ID)
		newVolumeList = append(newVolumeList, volume)
	}

	c.logger.Infof("found %d volumes, %d of which are orphaned", len(volumes), len(newVolumeList))
	return newVolumeList
}

// getOrphanedSSHKeys fetches all SSH keys then compares them against the
// provided list of nodes to determine if they are associated with any of
// them. It returns an error if the fetching process encounters any issues.
func (c *Civo) getOrphanedSSHKeys(ctx context.Context, nodes []sdk.Instance) ([]sdk.SSHKey, error) {
	c.logger.Infof("listing SSH keys")

	keys, err := c.client.GetSSHKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list SSH keys: %w", err)
	}

	// iterate over all keys and check if they are associated with any nodes
	orphanedKeys := make([]sdk.SSHKey, 0, len(keys))
	for _, key := range keys {
		var found bool

		// iterate through the nodes finding if they use the current key
		for _, node := range nodes {
			if node.SSHKeyID == key.ID {
				c.logger.Warnf("skipping SSH key %q: it is associated with the node instance with ID %q", key.Name, node.ID)
				found = true
				break
			}
		}

		if !found {
			c.logger.Infof("found orphaned SSH key %q - ID: %q", key.Name, key.ID)
			orphanedKeys = append(orphanedKeys, key)
		}
	}

	c.logger.Infof("found %d ssh keys, %d of which are orphaned", len(keys), len(orphanedKeys))
	return orphanedKeys, nil
}

// getOrphanedNetworks fetches all networks then compares them against the
// provided list of nodes to determine if they are associated with any of
// them. It returns an error if the fetching process encounters any issues.
func (c *Civo) getOrphanedNetworks(ctx context.Context, nodes []sdk.Instance, volumes []sdk.Volume) ([]sdk.Network, error) {
	c.logger.Infof("listing networks")

	networks, err := c.client.GetNetworks(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list networks: %w", err)
	}

	// iterate over all networks and check if they are associated with any nodes
	orphanedNetworks := make([]sdk.Network, 0, len(networks))
	for _, network := range networks {
		found := false

		// check if network name is "default", if so, skip it
		if network.Default {
			c.logger.Warnf("skipping network %q: it is the default network", network.Name)
			continue
		}

		// iterate through the nodes finding if they use the current network
		for _, node := range nodes {
			if node.NetworkID == network.ID {
				c.logger.Warnf("skipping network %q: it is associated with the node instance with ID %q", network.Name, node.ID)
				found = true
				break
			}
		}

		// iterate through the volumes finding if they use the current network
		for _, volume := range volumes {
			if volume.NetworkID == network.ID {
				c.logger.Warnf("skipping network %q: it is associated with the volume with ID %q", network.Name, volume.ID)
				found = true
				break
			}
		}

		if !found {
			c.logger.Infof("found orphaned network %q - ID: %q", network.Name, network.ID)
			orphanedNetworks = append(orphanedNetworks, network)
		}
	}

	c.logger.Infof("found %d networks, %d of which are orphaned", len(networks), len(orphanedNetworks))
	return orphanedNetworks, nil
}

// getOrphanedFirewalls fetches all firewalls then checks if they are associated
// with any node instance, cluster, or load balancer. It returns an error if the
// fetching process encounters any issues.
func (c *Civo) getOrphanedFirewalls(ctx context.Context) ([]sdk.Firewall, error) {
	c.logger.Infof("listing firewalls")

	firewalls, err := c.client.GetFirewalls(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list firewalls: %w", err)
	}

	// iterate over all firewalls and check if they are associated with any nodes
	orphanedFirewalls := make([]sdk.Firewall, 0, len(firewalls))
	for _, firewall := range firewalls {
		if firewall.ClusterCount > 0 || firewall.InstanceCount > 0 || firewall.LoadBalancerCount > 0 {
			c.logger.Warnf(
				"skipping firewall %q: it is associated with %d clusters, %d instances, and %d load balancers",
				firewall.Name, firewall.ClusterCount, firewall.InstanceCount, firewall.LoadBalancerCount,
			)
			continue
		}

		if firewall.NetworkID != "" {
			c.logger.Warnf("skipping firewall %q: it is associated with the network with ID %q", firewall.Name, firewall.NetworkID)
			continue
		}

		c.logger.Infof("found orphaned firewall %q - ID: %q", firewall.Name, firewall.ID)
		orphanedFirewalls = append(orphanedFirewalls, firewall)
	}

	c.logger.Infof("found %d firewalls, %d of which are orphaned", len(firewalls), len(orphanedFirewalls))
	return orphanedFirewalls, nil
}
