package civo

import (
	"fmt"
)

func (c *CivoConfiguration) NukeVolumes(CivoCmdOptions *CivoCmdOptions) error {
	vols, err := c.Client.ListVolumes()
	if err != nil {
		fmt.Println("err getting volumes", err)
	}

	for _, v := range vols {
		if CivoCmdOptions.Nuke {
			res, err := c.Client.DeleteVolume(v.ID)
			if err != nil {
				fmt.Println("err getting volumes", err)
			}
			fmt.Println("success delete: ", res.Result)
		} else {
			fmt.Printf("nuke set to %t, not removing volume %s\n", CivoCmdOptions.Nuke, v.ID)
		}
	}
	return nil
}
