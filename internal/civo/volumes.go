package civo

import (
	"fmt"

	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeVolumes deletes all volumes associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
//
//nolint:dupl // similar functions due to upstream packaging
func (c *Civo) NukeVolumes() error {
	c.logger.Infof("listing volumes")

	volumes, err := c.client.ListVolumes()
	if err != nil {
		return fmt.Errorf("unable to list volumes: %w", err)
	}

	c.logger.Infof("found %d volumes", len(volumes))

	for _, volume := range volumes {
		c.logger.Infof("found volume %q - ID: %q", volume.Name, volume.ID)

		if c.nameFilter != nil && !c.nameFilter.MatchString(volume.Name) {
			c.logger.Warnf("skipping volume %q: name does not match filter", volume.Name)
			continue
		}

		if c.nuke {
			c.logger.Infof("deleting volume %q", volume.Name)

			res, err := c.client.DeleteVolume(volume.ID)
			if err != nil {
				return fmt.Errorf("unable to delete volume %s: %w", volume.Name, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %q when deleting volume %s: %s", res.ErrorCode, volume.Name, res.ErrorDetails)
			}
			outputwriter.WriteStdoutf("deleted volume %s", volume.Name)
		} else {
			c.logger.Warnf("refusing to delete volume %s: nuke is not enabled", volume.Name)
		}
	}
	return nil
}
