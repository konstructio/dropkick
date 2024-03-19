package civo

import (
	"fmt"
)

func (c *CivoConfiguration) NukeNetworks() {
	networks, err := c.Client.ListNetworks()
	if err != nil {
		fmt.Println("err getting networks", err)
	}

	for _, n := range networks {
		res, err := c.Client.DeleteNetwork(n.ID)
		if err != nil {
			fmt.Println("err getting networks", err)
		}
		fmt.Println("success delete: ", &res.ID)
	}
}
