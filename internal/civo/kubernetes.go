package civo

import (
	"fmt"

	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeKubernetesClusters deletes all Kubernetes clusters associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeKubernetesClusters() error {
	c.logger.Printf("listing Kubernetes clusters")

	clusters, err := c.client.ListKubernetesClusters()
	if err != nil {
		return fmt.Errorf("unable to list Kubernetes clusters: %w", err)
	}

	c.logger.Printf("found %d clusters", len(clusters.Items))

	for _, cluster := range clusters.Items {
		c.logger.Printf("found cluster: name: %q - ID: %q", cluster.Name, cluster.ID)

		clusterVolumes, err := c.client.ListVolumesForCluster(cluster.ID)
		if err != nil {
			return fmt.Errorf("unable to list volumes for cluster %q: %w", cluster.ID, err)
		}

		for _, volume := range clusterVolumes {
			if c.nuke {
				c.logger.Printf("deleting volume %q for cluster %q", volume.ID, cluster.ID)
				res, err := c.client.DeleteVolume(volume.ID)
				if err != nil {
					return fmt.Errorf("unable to delete cluster %q volume %q: %w", cluster.ID, volume.ID, err)
				}

				if res.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %s when deleting volume %q: %s", res.ErrorCode, volume.ID, res.ErrorDetails)
				}

				outputwriter.WriteStdoutf("deleted volume %q for cluster %q", volume.ID, cluster.ID)
			} else {
				c.logger.Printf("refusing to delete volume %q for cluster %q: nuke is not enabled", volume.ID, cluster.ID)
			}
		}

		if c.nuke {
			c.logger.Printf("deleting cluster %q", cluster.ID)
			res, err := c.client.DeleteKubernetesCluster(cluster.ID)
			if err != nil {
				return fmt.Errorf("unable to delete cluster %q: %w", cluster.ID, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %s when deleting cluster %q: %s", res.ErrorCode, cluster.ID, res.ErrorDetails)
			}

			outputwriter.WriteStdoutf("deleted cluster %q", cluster.ID)
		} else {
			c.logger.Printf("refusing to delete cluster %q: nuke is not enabled", cluster.ID)
		}

		network, err := c.client.FindNetwork(cluster.Name)
		if err != nil {
			return fmt.Errorf("unable to find network for cluster %q: %w", cluster.ID, err)
		}

		if c.nuke {
			c.logger.Printf("deleting network %q", network.ID)
			res, err := c.client.DeleteNetwork(network.ID)
			if err != nil {
				return fmt.Errorf("unable to delete cluster network %q: %w", cluster.ID, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %s when deleting network %q: %s", res.ErrorCode, network.ID, res.ErrorDetails)
			}

			outputwriter.WriteStdoutf("deleted network %q", network.ID)
		} else {
			c.logger.Printf("refusing to delete network %q: nuke is not enabled", network.ID)
		}
	}

	return nil
}
