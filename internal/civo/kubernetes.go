package civo

import (
	"errors"
	"fmt"

	"github.com/civo/civogo"
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

		if c.nameFilter != nil && !c.nameFilter.MatchString(cluster.Name) {
			c.logger.Warnf("skipping cluster %q: name does not match filter", cluster.Name)
			continue
		}

		clusterVolumes, err := c.client.ListVolumesForCluster(cluster.ID)
		if err != nil {
			return fmt.Errorf("unable to list volumes for cluster %q: %w", cluster.Name, err)
		}

		for _, volume := range clusterVolumes {
			c.logger.Infof("found volume: name: %q - ID: %q", volume.Name, volume.ID)

			if c.nameFilter != nil && !c.nameFilter.MatchString(volume.Name) {
				c.logger.Warnf("skipping volume %q: name does not match filter", volume.Name)
				continue
			}

			if c.nuke {
				c.logger.Infof("deleting volume %q for cluster %q", volume.Name, cluster.Name)
				res, err := c.client.DeleteVolume(volume.ID)
				if err != nil {
					return fmt.Errorf("unable to delete cluster %q volume %q: %w", cluster.Name, volume.Name, err)
				}

				if res.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %q when deleting volume %q: %s", res.ErrorCode, volume.Name, res.ErrorDetails)
				}

				outputwriter.WriteStdoutf("deleted volume %q for cluster %q", volume.Name, cluster.Name)
			} else {
				c.logger.Warnf("refusing to delete volume %q for cluster %q: nuke is not enabled", volume.Name, cluster.Name)
			}
		}

		if c.nuke {
			c.logger.Infof("deleting cluster %q", cluster.Name)
			res, err := c.client.DeleteKubernetesCluster(cluster.ID)
			if err != nil {
				return fmt.Errorf("unable to delete cluster %q: %w", cluster.Name, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %q when deleting cluster %q: %s", res.ErrorCode, cluster.Name, res.ErrorDetails)
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

		if c.nameFilter != nil && !c.nameFilter.MatchString(network.Name) {
			c.logger.Warnf("skipping network %q: name does not match filter", network.Name)
			continue
		}

		if c.nuke {
			c.logger.Infof("deleting network %q", network.Name)
			res, err := c.client.DeleteNetwork(network.ID)
			if err != nil {
				return fmt.Errorf("unable to delete cluster network %q: %w", cluster.Name, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %q when deleting network %q: %s", res.ErrorCode, network.Name, res.ErrorDetails)
			}

			outputwriter.WriteStdoutf("deleted network %q", network.Name)
		} else {
			c.logger.Warnf("refusing to delete network %q: nuke is not enabled", network.Name)
		}
	}

	return nil
}
