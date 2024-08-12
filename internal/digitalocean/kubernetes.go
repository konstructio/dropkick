package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeKubernetesClusters deletes all Kubernetes clusters associated with the
// DigitalOcean client. It returns an error if the deletion process encounters
// any issues.
func (d *DigitalOcean) NukeKubernetesClusters() error {
	page := 1

	for {
		d.logger.Printf("listing Kubernetes clusters for page %d", page)

		clusters, res, err := d.client.Kubernetes.List(d.context, &godo.ListOptions{
			Page: page, PerPage: 50,
		})
		if err != nil {
			return fmt.Errorf("unable to list Kubernetes clusters for page %d: %w", page, err)
		}

		d.logger.Printf("found %d clusters", len(clusters))

		for _, cluster := range clusters {
			d.logger.Printf("found cluster: name: %q - ID: %q", cluster.Name, cluster.ID)

			// Delete the Kubernetes cluster
			if d.nuke {
				d.logger.Printf("deleting cluster %q", cluster.ID)
				_, err := d.client.Kubernetes.Delete(d.context, cluster.ID)
				if err != nil {
					return fmt.Errorf("unable to delete cluster %q: %w", cluster.ID, err)
				}
				outputwriter.WriteStdoutf("deleted cluster %q", cluster.ID)
			} else {
				d.logger.Printf("refusing to delete cluster %q: nuke is not enabled", cluster.ID)
			}

			foo, _, err := d.client.Kubernetes.ListAssociatedResourcesForDeletion(d.context, cluster.ID)
			if err != nil {
				return fmt.Errorf("unable to list associated resources for cluster %q: %w", cluster.ID, err)
			}

			d.logger.Printf("found %d load balancers, %d volumes, and %d volume snapshots for cluster %q", len(foo.LoadBalancers), len(foo.Volumes), len(foo.VolumeSnapshots), cluster.Name)

			// Delete load balancers associated with this cluster
			for _, loadbalancer := range foo.LoadBalancers {
				if d.nuke {
					d.logger.Printf("deleting loadbalancer %q for cluster %q", loadbalancer.ID, cluster.ID)
					_, err := d.client.LoadBalancers.Delete(d.context, loadbalancer.ID)
					if err != nil {
						return fmt.Errorf("unable to delete cluster %q loadbalancer %q: %w", cluster.ID, loadbalancer.ID, err)
					}
					outputwriter.WriteStdoutf("deleted loadbalancer %q for cluster %q", loadbalancer.ID, cluster.ID)
				} else {
					d.logger.Printf("refusing to delete loadbalancer %q for cluster %q: nuke is not enabled", loadbalancer.ID, cluster.ID)
				}
			}

			// Delete volumes associated with this cluster
			for _, volume := range foo.Volumes {
				if d.nuke {
					d.logger.Printf("deleting volume %q for cluster %q", volume.ID, cluster.ID)
					_, err := d.client.Storage.DeleteVolume(d.context, volume.ID)
					if err != nil {
						return fmt.Errorf("unable to delete cluster %q volume %q: %w", cluster.ID, volume.ID, err)
					}
					outputwriter.WriteStdoutf("deleted volume %q for cluster %q", volume.ID, cluster.ID)
				} else {
					d.logger.Printf("refusing to delete volume %q for cluster %q: nuke is not enabled", volume.ID, cluster.ID)
				}
			}

			// Delete volume snapshots associated with this cluster
			for _, snapshot := range foo.VolumeSnapshots {
				if d.nuke {
					d.logger.Printf("deleting volume snapshot %q for cluster %q", snapshot.ID, cluster.ID)
					_, err := d.client.Snapshots.Delete(d.context, snapshot.ID)
					if err != nil {
						return fmt.Errorf("unable to delete cluster %q volume snapshot %q: %w", cluster.ID, snapshot.ID, err)
					}
					outputwriter.WriteStdoutf("deleted volume snapshot %q for cluster %q", snapshot.ID, cluster.ID)
				} else {
					d.logger.Printf("refusing to delete volume snapshot %q for cluster %q: nuke is not enabled", snapshot.ID, cluster.ID)
				}
			}
		}

		// Exit if we've reached the last page.
		if res.Links == nil || res.Links.IsLastPage() {
			break
		}
	}

	return nil
}
