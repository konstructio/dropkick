//nolint:dupl // similar functions due to upstream packaging
package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeSSHKeys deletes all SSH keys associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeSSHKeys() error {
	c.logger.Infof("listing SSH keys")

	keys, err := c.client.GetAllSSHKeys(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list SSH keys: %w", err)
	}

	c.logger.Infof("found %d SSH keys", len(keys))

	for _, key := range keys {
		c.logger.Infof("found SSH key %q - ID: %q", key.Name, key.ID)

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(key.Name, c.nameFilter) {
			c.logger.Warnf("skipping SSH key %q: name does not match filter", key.Name)
			continue
		}

		if c.nuke {
			c.logger.Infof("deleting SSH key %q", key.Name)

			if err := c.client.DeleteSSHKey(context.Background(), key.ID); err != nil {
				return fmt.Errorf("unable to delete SSH key %s: %w", key.Name, err)
			}

			outputwriter.WriteStdoutf("deleted SSH key %s", key.Name)
		} else {
			c.logger.Warnf("refusing to delete SSH key %s: nuke is not enabled", key.Name)
		}
	}
	return nil
}
