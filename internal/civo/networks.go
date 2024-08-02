package civo

import (
	"fmt"
)

// NukeNetworks deletes all networks associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeNetworks() error {
	networks, err := c.client.ListNetworks()
	if err != nil {
		return fmt.Errorf("unable to list networks: %w", err)
	}

	for _, network := range networks {
		if c.nuke {
			res, err := c.client.DeleteNetwork(network.ID)
			if err != nil {
				return fmt.Errorf("unable to delete network %q: %w", network.ID, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %s when deleting network %q: %s", res.ErrorCode, network.ID, res.ErrorDetails)
			}
		} else {
			c.logger.Printf("refusing to delete network %q: nuke is not enabled", network.ID)
		}
	}

	return nil
}
