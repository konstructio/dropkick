package civo

import (
	"fmt"
)

func (c *CivoConfiguration) NukeNetworks(CivoCmdOptions *CivoCmdOptions) error {
	networks, err := c.Client.ListNetworks()
	if err != nil {
		fmt.Println("err getting networks", err)
	}

	for _, n := range networks {
		if CivoCmdOptions.Nuke {
			res, err := c.Client.DeleteNetwork(n.ID)
			if err != nil {
				fmt.Println("err getting networks", err)
			}
			fmt.Println("success delete: ", res.Result)
		} else {
			fmt.Printf("nuke set to %t, not removing network %s\n", CivoCmdOptions.Nuke, n.ID)
		}
	}
	return nil
}
