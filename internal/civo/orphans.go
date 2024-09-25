package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/civov2"
)

// NukeOrphanedResources deletes all subresources not in use by a compute
// instance in the Civo account targeted. It returns an error if the deletion
// process encounters any issues. The resources targeted by this function are:
// - Volumes
// - Object store credentials
// - SSH keys
// - Networks
// - Firewalls
func (c *Civo) NukeOrphanedResources() error {
	c.logger.Infof("orphaned resources enabled: looking for volumes, object store credentials, SSH keys, networks and firewalls not in use by any instances")

	// fetch all nodes first, we'll need them to check for orphaned resources
	c.logger.Infof("fetching all instances")
	nodes, err := c.getAllNodes()
	if err != nil {
		return fmt.Errorf("unable to fetch nodes: %w", err)
	}

	c.logger.Infof("found %d instances", len(nodes))

	// fetch orphaned volumes
	orphanedVolumes, err := c.getOrphanedVolumes()
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned volumes: %w", err)
	}

	c.logger.Infof("found %d orphaned volumes", len(orphanedVolumes))

	if err := c.deleteVolumes(orphanedVolumes); err != nil {
		return fmt.Errorf("unable to delete orphaned volumes: %w", err)
	}

	// fetch orphaned object store credentials
	orphanedObjectStoreCredentials, err := c.getOrphanedObjectStoreCredentials()
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned object store credentials: %w", err)
	}

	c.logger.Infof("found %d orphaned object store credentials", len(orphanedObjectStoreCredentials))

	if err := c.deleteObjectStoreCredentials(orphanedObjectStoreCredentials); err != nil {
		return fmt.Errorf("unable to delete orphaned object store credentials: %w", err)
	}

	// fetch orphaned SSH keys
	orphanedSSHKeys, err := c.getOrphanedSSHKeys(nodes)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned SSH keys: %w", err)
	}

	c.logger.Infof("found %d orphaned SSH keys", len(orphanedSSHKeys))

	if err := c.deleteSSHKeys(orphanedSSHKeys); err != nil {
		return fmt.Errorf("unable to delete orphaned SSH keys: %w", err)
	}

	// fetch orphaned networks
	orphanedNetworks, err := c.getOrphanedNetworks(nodes)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned networks: %w", err)
	}

	c.logger.Infof("found %d orphaned networks", len(orphanedNetworks))

	if err := c.deleteNetworks(orphanedNetworks); err != nil {
		return fmt.Errorf("unable to delete orphaned networks: %w", err)
	}

	// fetch orphaned firewalls
	orphanedFirewalls, err := c.getOrphanedFirewalls()
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned firewalls: %w", err)
	}

	c.logger.Infof("found %d orphaned firewalls", len(orphanedFirewalls))

	if err := c.deleteFirewalls(orphanedFirewalls); err != nil {
		return fmt.Errorf("unable to delete orphaned firewalls: %w", err)
	}

	return nil
}

// getOrphanedObjectStoreCredentials fetches all object store then the object
// store credentials and compares them against each other. If a credential is
// not used by any store, it is considered orphaned. It returns an error if the
// fetching process encounters any issues.
func (c *Civo) getOrphanedObjectStoreCredentials() ([]civov2.ObjectStoreCredential, error) {
	c.logger.Infof("listing object stores")

	objectStores, err := c.client.GetAllObjectStores(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to list object stores: %w", err)
	}

	c.logger.Infof("found %d object stores", len(objectStores))

	credentials, err := c.client.GetAllObjectStoreCredentials(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to list object store credentials: %w", err)
	}

	c.logger.Infof("found %d object store credentials", len(credentials))

	// iterate over all object stores and check if they are associated with any credentials
	orphanedCredentials := make([]civov2.ObjectStoreCredential, 0)
	for _, credential := range credentials {
		var found bool

		// iterate through the object stores finding if they use the current credential
		for _, objectStore := range objectStores {
			if objectStore.Credentials.ID == credential.ID {
				found = true
				break
			}
		}

		if !found {
			c.logger.Infof("found orphaned object store credential %q - ID: %q", credential.Name, credential.ID)
			orphanedCredentials = append(orphanedCredentials, credential)
		}
	}

	return orphanedCredentials, nil
}

// deleteObjectStoreCredentials deletes all object store credentials in the
// provided list. It returns an error if the deletion process encounters any
// issues.
func (c *Civo) deleteObjectStoreCredentials(credentials []civov2.ObjectStoreCredential) error {
	for _, credential := range credentials {
		if !c.nuke {
			c.logger.Warnf("refusing to delete object store credential %q: nuke is not enabled", credential.Name)
			continue
		}

		c.logger.Infof("deleting object store credential %q", credential.Name)

		if err := c.client.DeleteObjectStoreCredential(context.Background(), credential.ID); err != nil {
			return fmt.Errorf("unable to delete object store credential %s: %w", credential.Name, err)
		}

		c.logger.Infof("deleted object store credential %q", credential.Name)
	}

	return nil
}

// deleteVolumes deletes all volumes in the provided list. It returns an error
// if the deletion process encounters any issues.
func (c *Civo) deleteVolumes(volumes []civov2.Volume) error {
	for _, volume := range volumes {
		if !c.nuke {
			c.logger.Warnf("refusing to delete volume %q: nuke is not enabled", volume.Name)
			continue
		}

		c.logger.Infof("deleting volume %q", volume.Name)

		if err := c.client.DeleteVolume(context.Background(), volume.ID); err != nil {
			return fmt.Errorf("unable to delete volume %q: %w", volume.Name, err)
		}

		c.logger.Infof("deleted volume %s", volume.Name)
	}

	return nil
}

// deleteSSHKeys deletes all SSH keys in the provided list. It returns an error
// if the deletion process encounters any issues.
func (c *Civo) deleteSSHKeys(keys []civov2.SSHKey) error {
	for _, key := range keys {
		if !c.nuke {
			c.logger.Warnf("refusing to delete SSH key %q: nuke is not enabled", key.Name)
			continue
		}

		c.logger.Infof("deleting SSH key %q", key.Name)

		if err := c.client.DeleteSSHKey(context.Background(), key.ID); err != nil {
			return fmt.Errorf("unable to delete SSH key %q: %w", key.Name, err)
		}

		c.logger.Infof("deleted SSH key %q", key.Name)
	}

	return nil
}

// deleteNetworks deletes all networks in the provided list. It returns an error
// if the deletion process encounters any issues.
func (c *Civo) deleteNetworks(networks []civov2.Network) error {
	for _, network := range networks {
		if !c.nuke {
			c.logger.Warnf("refusing to delete network %q: nuke is not enabled", network.Name)
			continue
		}

		c.logger.Infof("deleting network %q", network.Name)

		if err := c.client.DeleteNetwork(context.Background(), network.ID); err != nil {
			return fmt.Errorf("unable to delete network %s: %w", network.Name, err)
		}

		c.logger.Infof("deleted network %q", network.Name)
	}

	return nil
}

// deleteFirewalls deletes all firewalls in the provided list. It returns an error
// if the deletion process encounters any issues.
func (c *Civo) deleteFirewalls(firewalls []civov2.Firewall) error {
	for _, firewall := range firewalls {
		if !c.nuke {
			c.logger.Warnf("refusing to delete firewall %q: nuke is not enabled", firewall.Name)
			continue
		}

		c.logger.Infof("deleting firewall %q", firewall.Name)

		err := c.client.DeleteFirewall(context.Background(), firewall.ID)
		if err != nil {
			return fmt.Errorf("unable to delete firewall %q: %w", firewall.Name, err)
		}

		c.logger.Infof("deleted firewall %q", firewall.Name)
	}

	return nil
}

// getAllNodes fetches all nodes in the Civo account. Since the results are
// paginated, this function will fetch all pages of nodes and return them as a
// single list. It returns an error if the fetching process encounters any
// issues.
func (c *Civo) getAllNodes() ([]civov2.Instance, error) {
	c.logger.Infof("listing nodes")

	nodes, err := c.client.GetAllInstances(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to list nodes: %w", err)
	}

	c.logger.Infof("found %d nodes on all pages", len(nodes))
	return nodes, nil
}

// getOrphanedVolumes fetches all volumes that are not attached to any node
// instance instead of relying if they are referenced by a node instance. It
// returns an error if the fetching process encounters any issues.
func (c *Civo) getOrphanedVolumes() ([]civov2.Volume, error) {
	c.logger.Infof("listing volumes")

	volumes, err := c.client.GetAllVolumes(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to list volumes: %w", err)
	}

	newVolumeList := make([]civov2.Volume, 0, len(volumes))

	for _, volume := range volumes {
		if volume.Status == volumeStatusAttached {
			c.logger.Warnf("skipping volume %q: it is attached to the node instance with ID %q", volume.Name, volume.InstanceID)
			continue
		}

		c.logger.Infof("found orphaned volume (not attached) %q - ID: %q", volume.Name, volume.ID)
		newVolumeList = append(newVolumeList, volume)
	}

	return newVolumeList, nil
}

// getOrphanedSSHKeys fetches all SSH keys then compares them against the
// provided list of nodes to determine if they are associated with any of
// them. It returns an error if the fetching process encounters any issues.
func (c *Civo) getOrphanedSSHKeys(nodes []civov2.Instance) ([]civov2.SSHKey, error) {
	c.logger.Infof("listing SSH keys")

	keys, err := c.client.GetAllSSHKeys(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to list SSH keys: %w", err)
	}

	c.logger.Infof("found %d SSH keys", len(keys))

	// iterate over all keys and check if they are associated with any nodes
	orphanedKeys := make([]civov2.SSHKey, 0, len(keys))
	for _, key := range keys {
		var found bool

		// iterate through the nodes finding if they use the current key
		for _, node := range nodes {
			if node.SSHKeyID == key.ID {
				found = true
				break
			}
		}

		if !found {
			c.logger.Infof("found orphaned SSH key %q - ID: %q", key.Name, key.ID)
			orphanedKeys = append(orphanedKeys, key)
		}
	}

	return orphanedKeys, nil
}

// getOrphanedNetworks fetches all networks then compares them against the
// provided list of nodes to determine if they are associated with any of
// them. It returns an error if the fetching process encounters any issues.
func (c *Civo) getOrphanedNetworks(nodes []civov2.Instance) ([]civov2.Network, error) {
	c.logger.Infof("listing networks")

	networks, err := c.client.GetAllNetworks(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to list networks: %w", err)
	}

	c.logger.Infof("found %d networks", len(networks))

	// iterate over all networks and check if they are associated with any nodes
	orphanedNetworks := make([]civov2.Network, 0, len(networks))
	for _, network := range networks {
		var found bool

		// iterate through the nodes finding if they use the current network
		for _, node := range nodes {
			if node.NetworkID == network.ID {
				found = true
				break
			}
		}

		if !found {
			c.logger.Infof("found orphaned network %q - ID: %q", network.Name, network.ID)
			orphanedNetworks = append(orphanedNetworks, network)
		}
	}

	return orphanedNetworks, nil
}

// getOrphanedFirewalls fetches all firewalls then checks if they are associated
// with any node instance, cluster, or load balancer. It returns an error if the
// fetching process encounters any issues.
func (c *Civo) getOrphanedFirewalls() ([]civov2.Firewall, error) {
	c.logger.Infof("listing firewalls")

	firewalls, err := c.client.GetAllFirewalls(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to list firewalls: %w", err)
	}

	c.logger.Infof("found %d firewalls", len(firewalls))

	// iterate over all firewalls and check if they are associated with any nodes
	orphanedFirewalls := make([]civov2.Firewall, 0, len(firewalls))
	for _, firewall := range firewalls {
		if firewall.ClusterCount > 0 || firewall.InstanceCount > 0 || firewall.LoadBalancerCount > 0 {
			c.logger.Warnf(
				"skipping firewall %q: it is associated with %d clusters, %d instances, and %d load balancers",
				firewall.Name, firewall.ClusterCount, firewall.InstanceCount, firewall.LoadBalancerCount,
			)
			continue
		}

		c.logger.Infof("found orphaned firewall %q - ID: %q", firewall.Name, firewall.ID)
		orphanedFirewalls = append(orphanedFirewalls, firewall)
	}

	return orphanedFirewalls, nil
}
