package civo

import (
	"fmt"
)

func (c *CivoConfiguration) NukeKubernetesClusters(CivoCmdOptions *CivoCmdOptions) {
	fmt.Println("deleting all kubernetes clusters")
	fmt.Println(CivoCmdOptions.Region)
	fmt.Println(CivoCmdOptions.Nuke)

	// clusters, err := client.ListKubernetesClusters()
	// if err != nil {
	// 	fmt.Println("err getting clusters", err)
	// }

	// for _, cl := range clusters.Items {
	// 	fmt.Println("cluster name: ", cl.Name)
	// 	fmt.Println("cluster id: ", cl.ID)
	// 	clusterVolumes, err := client.ListVolumesForCluster(cl.ID)
	// 	if err != nil {
	// 		fmt.Println("err getting cluster volumes", err)
	// 	}
	// 	for _, v := range clusterVolumes {
	// 		res, err := client.DeleteVolume(v.ID)
	// 		if err != nil {
	// 			fmt.Println("err deleting cluster volumes", err)
	// 		}
	// 		fmt.Println("cluster vol delete http code ", res.ErrorCode)
	// 	}

	// 	res, err := client.DeleteKubernetesCluster(cl.ID)
	// 	if err != nil {
	// 		fmt.Println("err deleting cluster", err)
	// 	}
	// 	fmt.Println("cluster delete http code ", res.ErrorCode)
	// 	network, err := client.FindNetwork(cl.Name)
	// 	if err != nil {
	// 		fmt.Println("err finding network", err)
	// 	}
	// 	res, err = client.DeleteNetwork(network.ID)
	// 	fmt.Println("delete network http code ", res.ErrorCode)
	// }

}
