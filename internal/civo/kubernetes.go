package civo

import (
	"errors"
	"fmt"

	"github.com/civo/civogo"
	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeKubernetesClusters deletes all Kubernetes clusters associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeKubernetesClusters() error {
	c.logger.Infof("listing Kubernetes clusters")

	clusters, err := c.client.ListKubernetesClusters()
	if err != nil {
		return fmt.Errorf("unable to list Kubernetes clusters: %w", err)
	}

	c.logger.Infof("found %d clusters", len(clusters.Items))

	for _, cluster := range clusters.Items {
		c.logger.Infof("found cluster: name: %q - ID: %q", cluster.Name, cluster.ID)

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(cluster.Name, c.nameFilter) {
			c.logger.Warnf("skipping cluster %q: name does not match filter", cluster.Name)
			continue
		}

		clusterVolumes, err := c.client.ListVolumesForCluster(cluster.ID)
		if err != nil {
			return fmt.Errorf("unable to list volumes for cluster %q: %w", cluster.Name, err)
		}

		for _, volume := range clusterVolumes {
			c.logger.Infof("found volume: name: %q - ID: %q", volume.Name, volume.ID)

			// We don't filter by volumes since Volume names are specific to the PVC
			// that created them (so on a cluster called "foo", the PVC would be called
			// "pvc-82d40d15-5ce2-418d-bb81-09a0349ec975").
			// If we land in this section, it means a Kubernetes cluster was already
			// found by matching the name filter, so we can safely delete all volumes
			// associated with it.

			if c.nuke {
				c.logger.Infof("deleting volume %q for cluster %q", volume.Name, cluster.Name)

				if _, err := c.client.DeleteVolume(volume.ID); err != nil {
					return fmt.Errorf("unable to delete cluster %q volume %q: %w", cluster.Name, volume.Name, err)
				}

				outputwriter.WriteStdoutf("deleted volume %q for cluster %q", volume.Name, cluster.Name)
			} else {
				c.logger.Warnf("refusing to delete volume %q for cluster %q: nuke is not enabled", volume.Name, cluster.Name)
			}
		}

		if c.nuke {
			c.logger.Infof("deleting cluster %q", cluster.Name)

			if _, err := c.client.DeleteKubernetesCluster(cluster.ID); err != nil {
				return fmt.Errorf("unable to delete cluster %q: %w", cluster.Name, err)
			}

			outputwriter.WriteStdoutf("deleted cluster %q", cluster.Name)
		} else {
			c.logger.Warnf("refusing to delete cluster %q: nuke is not enabled", cluster.Name)
		}

		network, err := c.client.FindNetwork(cluster.ID)
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				c.logger.Warnf("no network found for cluster %q", cluster.Name)
				continue
			}

			return fmt.Errorf("unable to find network for cluster %q: %w", cluster.Name, err)
		}

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(network.Name, c.nameFilter) {
			c.logger.Warnf("skipping network %q: name does not match filter", network.Name)
			continue
		}

		if c.nuke {
			c.logger.Infof("deleting network %q", network.Name)

			if _, err := c.client.DeleteNetwork(network.ID); err != nil {
				return fmt.Errorf("unable to delete cluster network %q: %w", cluster.Name, err)
			}

			outputwriter.WriteStdoutf("deleted network %q", network.Name)
		} else {
			c.logger.Warnf("refusing to delete network %q: nuke is not enabled", network.Name)
		}
	}

	return nil
}
