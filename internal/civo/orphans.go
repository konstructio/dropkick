package civo

import (
	"fmt"

	"github.com/civo/civogo"
	"golang.org/x/sync/errgroup"
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

	// fetch orphaned SSH keys
	orphanedSSHKeys, err := c.getOrphanedSSHKeys(nodes)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned SSH keys: %w", err)
	}

	// fetch orphaned networks
	orphanedNetworks, err := c.getOrphanedNetworks(nodes)
	if err != nil {
		return fmt.Errorf("unable to fetch orphaned networks: %w", err)
	}

	// create parallel executions for deleting orphaned resources
	eg := errgroup.Group{}

	// delete orphaned volumes
	eg.Go(func() error {
		return c.deleteVolumes(orphanedVolumes)
	})

	// delete orphaned SSH keys
	eg.Go(func() error {
		return c.deleteSSHKeys(orphanedSSHKeys)
	})

	// delete orphaned networks
	eg.Go(func() error {
		return c.deleteNetworks(orphanedNetworks)
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("unable to delete orphaned resources: %w", err)
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

		res, err := c.client.DeleteVolume(volume.ID)
		if err != nil {
			return fmt.Errorf("unable to delete volume %s: %w", volume.Name, err)
		}

		if res.ErrorCode != "200" {
			return fmt.Errorf("Civo returned an error code %q when deleting volume %s: %s", res.ErrorCode, volume.Name, res.ErrorDetails)
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

		res, err := c.client.DeleteSSHKey(key.ID)
		if err != nil {
			return fmt.Errorf("unable to delete SSH key %s: %w", key.Name, err)
		}

		if res.ErrorCode != "200" {
			return fmt.Errorf("Civo returned an error code %q when deleting SSH key %s: %s", res.ErrorCode, key.Name, res.ErrorDetails)
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

		res, err := c.client.DeleteNetwork(network.ID)
		if err != nil {
			return fmt.Errorf("unable to delete network %s: %w", network.Name, err)
		}

		if res.ErrorCode != "200" {
			return fmt.Errorf("Civo returned an error code %q when deleting network %s: %s", res.ErrorCode, network.Name, res.ErrorDetails)
		}
		c.logger.Infof("deleted network %s", network.Name)
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
