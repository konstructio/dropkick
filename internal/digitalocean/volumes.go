package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

func (d *DigitalOcean) NukeVolumes() error {
	d.logger.Infof("listing volumes")

	volumes, _, err := d.client.Storage.ListVolumes(d.context, &godo.ListVolumeParams{})
	if err != nil {
		return fmt.Errorf("unable to list volumes: %w", err)
	}

	d.logger.Infof("found %d volumes", len(volumes))

	for _, volume := range volumes {
		d.logger.Infof("found volume %q", volume.ID)

		if d.nuke {
			d.logger.Infof("deleting volume %q", volume.ID)
			_, err := d.client.Storage.DeleteVolume(d.context, volume.ID)
			if err != nil {
				return fmt.Errorf("unable to delete volume %s: %w", volume.ID, err)
			}
			outputwriter.WriteStdoutf("deleted volume %q", volume.ID)
		} else {
			d.logger.Warnf("refusing to delete volume %s: nuke is not enabled", volume.ID)
		}
	}

	return nil
}
