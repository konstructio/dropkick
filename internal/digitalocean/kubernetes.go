package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
)

// NukeKubernetesClusters deletes all Kubernetes clusters associated with the
// DigitalOcean client. It returns an error if the deletion process encounters
// any issues.
func (c *DigitalOcean) NukeKubernetesClusters() error {
	page := 1
	for {
		clusters, res, err := c.client.Kubernetes.List(c.context, &godo.ListOptions{
			Page: page, PerPage: 50,
		})
		if err != nil {
			return fmt.Errorf("unable to list Kubernetes clusters for page %d: %w", page, err)
		}

		for _, cluster := range clusters {
			c.logger.Printf("found cluster: name: %q - ID: %q", cluster.Name, cluster.ID)

			if c.nuke {
				_, err := c.client.Kubernetes.Delete(c.context, cluster.ID)
				if err != nil {
					return fmt.Errorf("unable to delete cluster %q: %w", cluster.ID, err)
				}
			} else {
				c.logger.Printf("refusing to delete cluster %q: nuke is not enabled", cluster.ID)
			}
		}

		// Exit if we've reached the last page.
		if res.Links == nil || res.Links.IsLastPage() {
			break
		}
	}

	return nil
}
