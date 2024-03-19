package civo

import (
	"fmt"
	"github.com/civo/civogo"
)

func (c *CivoConfiguration) NukeVolumes(client *civogo.Client) {
	vols, err := client.ListVolumes()
	if err != nil {
		fmt.Println("err getting volumes", err)
	}

	for _, v := range vols {
		res, err := client.DeleteVolume(v.ID)
		if err != nil {
			fmt.Println("err getting volumes", err)
		}
		fmt.Println("success delete: ", &res.ErrorCode)
	}
}
