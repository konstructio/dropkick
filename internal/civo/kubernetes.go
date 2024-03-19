package civo

import (
	"fmt"
)

func (c *CivoConfiguration) NukeKubernetesClusters(CivoCmdOptions *CivoCmdOptions) error {
	fmt.Printf("deleting all kubernetes clusters in %s with --nuke set to %t \n", CivoCmdOptions.Region, CivoCmdOptions.Nuke)

	clusters, err := c.Client.ListKubernetesClusters()
	if err != nil {
		fmt.Println("err getting clusters", err)
	}

	for _, cl := range clusters.Items {
		fmt.Println("cluster name: ", cl.Name)
		fmt.Println("cluster id: ", cl.ID)
		clusterVolumes, err := c.Client.ListVolumesForCluster(cl.ID)
		if err != nil {
			fmt.Println("err getting cluster volumes", err)
		}
		for _, v := range clusterVolumes {
			if CivoCmdOptions.Nuke {
				res, err := c.Client.DeleteVolume(v.ID)
				if err != nil {
					fmt.Println("err deleting cluster volumes", err)
				}
				fmt.Println("cluster vol delete http code ", res.ErrorCode)
			} else {
				fmt.Printf("nuke set to %t, not removing volume %s\n", CivoCmdOptions.Nuke, v.ID)
			}
		}
		if CivoCmdOptions.Nuke {
			res, err := c.Client.DeleteKubernetesCluster(cl.ID)
			if err != nil {
				fmt.Println("err deleting cluster", err)
			}
			fmt.Println("cluster delete http code ", res.ErrorCode)
		} else {
			fmt.Printf("nuke set to %t, not cluster %s\n", CivoCmdOptions.Nuke, cl.ID)
		}

		network, err := c.Client.FindNetwork(cl.Name)
		if err != nil {
			fmt.Println("err finding network", err)
		}

		if CivoCmdOptions.Nuke {
			res, err := c.Client.DeleteNetwork(network.ID)
			if err != nil {
				fmt.Println("err deleting cluster network", err)
			}
			fmt.Println("delete network http code ", res.ErrorCode)
		} else {
			fmt.Printf("nuke set to %t, not removing network %s\n", CivoCmdOptions.Nuke, cl.ID)
		}
	}
	return nil
}
