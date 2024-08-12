package civo

import (
	"fmt"

	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeVolumes deletes all volumes associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeVolumes() error {
	c.logger.Printf("listing volumes")

	volumes, err := c.client.ListVolumes()
	if err != nil {
		return fmt.Errorf("unable to list volumes: %w", err)
	}

	c.logger.Printf("found %d volumes", len(volumes))

	for _, volume := range volumes {
		c.logger.Printf("found volume %q", volume.ID)

		if c.nuke {
			c.logger.Printf("deleting volume %q", volume.ID)
			res, err := c.client.DeleteVolume(volume.ID)
			if err != nil {
				return fmt.Errorf("unable to delete volume %s: %w", volume.ID, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %s when deleting volume %s: %s", res.ErrorCode, volume.ID, res.ErrorDetails)
			}
			outputwriter.WriteStdoutf("deleted volume %s", volume.ID)
		} else {
			c.logger.Printf("refusing to delete volume %s: nuke is not enabled", volume.ID)
		}
	}
	return nil
}
