package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

const volumeStatusAttached = "attached"

// NukeVolumes deletes all volumes associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeVolumes() error {
	c.logger.Infof("listing volumes")

	volumes, err := c.client.GetAllVolumes(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list volumes: %w", err)
	}

	c.logger.Infof("found %d volumes", len(volumes))

	for _, volume := range volumes {
		c.logger.Infof("found volume %q - ID: %q (attached? %v)", volume.Name, volume.ID, volume.Status == volumeStatusAttached)

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(volume.Name, c.nameFilter) {
			c.logger.Warnf("skipping volume %q: name does not match filter", volume.Name)
			continue
		}

		if c.nuke {
			if volume.Status == volumeStatusAttached {
				c.logger.Warnf("refusing to delete volume %s: it is attached to the node instance with ID %q", volume.Name, volume.InstanceID)
				continue
			}

			c.logger.Infof("deleting volume %q", volume.Name)

			if err := c.client.DeleteVolume(context.Background(), volume.ID); err != nil {
				return fmt.Errorf("unable to delete volume %s: %w", volume.Name, err)
			}

			outputwriter.WriteStdoutf("deleted volume %s", volume.Name)
		} else {
			c.logger.Warnf("refusing to delete volume %s: nuke is not enabled", volume.Name)
		}
	}
	return nil
}
