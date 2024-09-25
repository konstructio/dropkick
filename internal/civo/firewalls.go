package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeFirewalls deletes all firewalls associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeFirewalls() error {
	c.logger.Infof("listing firewalls")

	firewalls, err := c.client.GetAllFirewalls(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list firewalls: %w", err)
	}

	c.logger.Infof("found %d firewalls", len(firewalls))

	for _, firewall := range firewalls {
		c.logger.Infof("found firewall: name: %q - ID: %q", firewall.Name, firewall.ID)

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(firewall.Name, c.nameFilter) {
			c.logger.Warnf("skipping firewall %q: name does not match filter", firewall.Name)
			continue
		}

		if c.nuke {
			c.logger.Infof("deleting firewall %q", firewall.Name)

			err := c.client.DeleteFirewall(context.Background(), firewall.ID)
			if err != nil {
				return fmt.Errorf("unable to delete firewall %q: %w", firewall.Name, err)
			}

			outputwriter.WriteStdoutf("deleted firewall %q", firewall.Name)
		} else {
			c.logger.Warnf("refusing to delete firewall %q: nuke is not enabled", firewall.Name)
		}
	}

	return nil
}
