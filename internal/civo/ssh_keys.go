//nolint:dupl // similar functions due to upstream packaging
package civo

import (
	"fmt"

	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeSSHKeys deletes all SSH keys associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeSSHKeys() error {
	c.logger.Infof("listing SSH keys")

	keys, err := c.client.ListSSHKeys()
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

			res, err := c.client.DeleteSSHKey(key.ID)
			if err != nil {
				return fmt.Errorf("unable to delete SSH key %s: %w", key.Name, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %q when deleting SSH key %s: %s", res.ErrorCode, key.Name, res.ErrorDetails)
			}
			outputwriter.WriteStdoutf("deleted SSH key %s", key.Name)
		} else {
			c.logger.Warnf("refusing to delete SSH key %s: nuke is not enabled", key.Name)
		}
	}
	return nil
}
