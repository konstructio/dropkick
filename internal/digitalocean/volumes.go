package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
)

func (d *DigitalOcean) NukeVolumes() error {
	volumes, _, err := d.client.Storage.ListVolumes(d.context, &godo.ListVolumeParams{})
	if err != nil {
		return fmt.Errorf("unable to list volumes: %w", err)
	}

	for _, volume := range volumes {
		if d.nuke {
			_, err := d.client.Storage.DeleteVolume(d.context, volume.ID)
			if err != nil {
				return fmt.Errorf("unable to delete volume %s: %w", volume.ID, err)
			}
		} else {
			d.logger.Printf("refusing to delete volume %s: nuke is not enabled", volume.ID)
		}
	}

	return nil
}
