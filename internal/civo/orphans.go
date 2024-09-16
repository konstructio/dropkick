package civo

import (
	"fmt"

	"github.com/civo/civogo"
)

// NukeOrphanedResources deletes all subresources not in use by a compute
// instance in the Civo account targeted. It returns an error if the deletion
// process encounters any issues. The resources targeted by this function are:
// - Volumes
// - SSH keys
// - Networks
func (c *Civo) NukeOrphanedResources() error {
	c.logger.Infof("orphaned resources enabled: looking for volumes, SSH keys, and networks not in use by any instances")

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

// deleteVolumes deletes all volumes in the provided list. It returns an error
// if the deletion process encounters any issues.
func (c *Civo) deleteVolumes(volumes []civogo.Volume) error {
	for _, volume := range volumes {
		if !c.nuke {
			c.logger.Warnf("refusing to delete volume %q: nuke is not enabled", volume.Name)
			continue
		}

		c.logger.Infof("deleting volume %q", volume.Name)

		if _, err := c.client.DeleteVolume(volume.ID); err != nil {
			return fmt.Errorf("unable to delete volume %s: %w", volume.Name, err)
		}

		c.logger.Infof("deleted volume %s", volume.Name)
	}

	return nil
}

// deleteSSHKeys deletes all SSH keys in the provided list. It returns an error
// if the deletion process encounters any issues.
func (c *Civo) deleteSSHKeys(keys []civogo.SSHKey) error {
	for _, key := range keys {
		if !c.nuke {
			c.logger.Warnf("refusing to delete SSH key %s: nuke is not enabled", key.Name)
			continue
		}

		c.logger.Infof("deleting SSH key %q", key.Name)

		if _, err := c.client.DeleteSSHKey(key.ID); err != nil {
			return fmt.Errorf("unable to delete SSH key %s: %w", key.Name, err)
		}

		c.logger.Infof("deleted SSH key %s", key.Name)
	}

	return nil
}

// deleteNetworks deletes all networks in the provided list. It returns an error
// if the deletion process encounters any issues.
func (c *Civo) deleteNetworks(networks []civogo.Network) error {
	for _, network := range networks {
		if !c.nuke {
			c.logger.Warnf("refusing to delete network %s: nuke is not enabled", network.Name)
			continue
		}

		c.logger.Infof("deleting network %q", network.Name)

		if _, err := c.client.DeleteNetwork(network.ID); err != nil {
			return fmt.Errorf("unable to delete network %s: %w", network.Name, err)
		}

		c.logger.Infof("deleted network %s", network.Name)
	}

	return nil
}

// deleteFirewalls deletes all firewalls in the provided list. It returns an error
// if the deletion process encounters any issues.
func (c *Civo) deleteFirewalls(firewalls []*civogo.Firewall) error {
	for _, firewall := range firewalls {
		if !c.nuke {
			c.logger.Warnf("refusing to delete firewall %s: nuke is not enabled", firewall.Name)
			continue
		}

		c.logger.Infof("deleting firewall %q", firewall.Name)

		_, err := c.client.DeleteFirewall(firewall.ID)
		if err != nil {
			return fmt.Errorf("unable to delete firewall %s: %w", firewall.Name, err)
		}

		c.logger.Infof("deleted firewall %s", firewall.Name)
	}

	return nil
}

// getAllNodes fetches all nodes in the Civo account. Since the results are
// paginated, this function will fetch all pages of nodes and return them as a
// single list. It returns an error if the fetching process encounters any
// issues.
func (c *Civo) getAllNodes() ([]civogo.Instance, error) {
	var nodes []civogo.Instance

	perPage := 100
	for page := 1; ; page++ {
		c.logger.Infof("listing page %d of nodes", page)

		nodesPage, err := c.client.ListInstances(page, perPage)
		if err != nil {
			return nil, fmt.Errorf("unable to list nodes: %w", err)
		}

		c.logger.Infof("found %d nodes on page %d", len(nodesPage.Items), page)
		nodes = append(nodes, nodesPage.Items...)

		if nodesPage.Pages == page {
			break
		}
	}

	c.logger.Infof("found %d nodes on all pages", len(nodes))
	return nodes, nil
}

// getOrphanedVolumes fetches all volumes that are not attached to any node
// instance instead of relying if they are referenced by a node instance. It
// returns an error if the fetching process encounters any issues.
func (c *Civo) getOrphanedVolumes() ([]civogo.Volume, error) {
	c.logger.Infof("listing volumes")

	volumes, err := c.client.ListVolumes()
	if err != nil {
		return nil, fmt.Errorf("unable to list volumes: %w", err)
	}

	newVolumeList := make([]civogo.Volume, 0, len(volumes))
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
func (c *Civo) getOrphanedSSHKeys(nodes []civogo.Instance) ([]civogo.SSHKey, error) {
	var keys []civogo.SSHKey

	c.logger.Infof("listing SSH keys")

	keys, err := c.client.ListSSHKeys()
	if err != nil {
		return nil, fmt.Errorf("unable to list SSH keys: %w", err)
	}

	c.logger.Infof("found %d SSH keys", len(keys))

	// iterate over all keys and check if they are associated with any nodes
	orphanedKeys := make([]civogo.SSHKey, 0, len(keys))
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
func (c *Civo) getOrphanedNetworks(nodes []civogo.Instance) ([]civogo.Network, error) {
	c.logger.Infof("listing networks")

	networks, err := c.client.ListNetworks()
	if err != nil {
		return nil, fmt.Errorf("unable to list networks: %w", err)
	}

	c.logger.Infof("found %d networks", len(networks))

	// iterate over all networks and check if they are associated with any nodes
	orphanedNetworks := make([]civogo.Network, 0, len(networks))
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
func (c *Civo) getOrphanedFirewalls() ([]*civogo.Firewall, error) {
	c.logger.Infof("listing firewalls")

	firewalls, err := c.client.ListFirewalls()
	if err != nil {
		return nil, fmt.Errorf("unable to list firewalls: %w", err)
	}

	c.logger.Infof("found %d firewalls", len(firewalls))

	// iterate over all firewalls and check if they are associated with any nodes
	orphanedFirewalls := make([]*civogo.Firewall, 0, len(firewalls))
	for _, firewall := range firewalls {
		if firewall.ClusterCount > 0 || firewall.InstanceCount > 0 || firewall.LoadBalancerCount > 0 {
			c.logger.Warnf(
				"skipping firewall %q: it is associated with %d clusters, %d instances, and %d load balancers",
				firewall.Name, firewall.ClusterCount, firewall.InstanceCount, firewall.LoadBalancerCount,
			)
			continue
		}

		c.logger.Infof("found orphaned firewall %q - ID: %q", firewall.Name, firewall.ID)
		orphanedFirewalls = append(orphanedFirewalls, &firewall)
	}

	return orphanedFirewalls, nil
}
