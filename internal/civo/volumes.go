package civo

import (
	"fmt"
)

// NukeVolumes deletes all volumes associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeVolumes() error {
	volumes, err := c.client.ListVolumes()
	if err != nil {
		return fmt.Errorf("unable to list volumes: %w", err)
	}

	for _, volume := range volumes {
		if c.nuke {
			res, err := c.client.DeleteVolume(volume.ID)
			if err != nil {
				return fmt.Errorf("unable to delete volume %s: %w", volume.ID, err)
			}

			if res.ErrorCode != "200" {
				return fmt.Errorf("Civo returned an error code %s when deleting volume %s: %s", res.ErrorCode, volume.ID, res.ErrorDetails)
			}
		} else {
			c.logger.Printf("refusing to delete volume %s: nuke is not enabled", volume.ID)
		}
	}
	return nil
}
