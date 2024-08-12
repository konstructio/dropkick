package civo

import (
	"fmt"

	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeNetworks deletes all networks associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeNetworks() error {
	c.logger.Printf("listing networks")

	networks, err := c.client.ListNetworks()
	if err != nil {
		return fmt.Errorf("unable to list networks: %w", err)
	}

	c.logger.Printf("found %d networks", len(networks))

	for _, network := range networks {
		c.logger.Printf("found network: name: %q - ID: %q", network.Name, network.ID)

		if c.nuke {
			c.logger.Printf("deleting network %q", network.ID)
			res, err := c.client.DeleteNetwork(network.ID)
			if err != nil {
				return fmt.Errorf("unable to delete network %q: %w", network.ID, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %s when deleting network %q: %s", res.ErrorCode, network.ID, res.ErrorDetails)
			}

			outputwriter.WriteStdout("deleted network %q", network.ID)
		} else {
			c.logger.Printf("refusing to delete network %q: nuke is not enabled", network.ID)
		}
	}

	return nil
}
