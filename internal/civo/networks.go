//nolint:dupl // similar functions due to upstream packaging
package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeNetworks deletes all networks associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeNetworks() error {
	c.logger.Infof("listing networks")

	networks, err := c.client.GetAllNetworks(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list networks: %w", err)
	}

	c.logger.Infof("found %d networks", len(networks))

	for _, network := range networks {
		c.logger.Infof("found network: name: %q - ID: %q", network.Name, network.ID)

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(network.Name, c.nameFilter) {
			c.logger.Warnf("skipping network %q: name does not match filter", network.Name)
			continue
		}

		if c.nuke {
			c.logger.Infof("deleting network %q", network.Name)

			if err := c.client.DeleteNetwork(context.Background(), network.ID); err != nil {
				return fmt.Errorf("unable to delete network %q: %w", network.Name, err)
			}

			outputwriter.WriteStdoutf("deleted network %q", network.Name)
		} else {
			c.logger.Warnf("refusing to delete network %q: nuke is not enabled", network.Name)
		}
	}

	return nil
}
