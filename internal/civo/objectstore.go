package civo

import (
	"fmt"

	"github.com/civo/civogo"
)

func (c *CivoConfiguration) NukeObjectStores(CivoCmdOptions *CivoCmdOptions) error {

	pageOfItems, err := c.Client.ListObjectStores()
	if err != nil {
		fmt.Println("err getting pg items")
	}
	fmt.Println(pageOfItems.Pages) // 4

	i := 0
	for i < pageOfItems.Pages {
		deletePageOfObjectStores(c.Client, CivoCmdOptions)
		deletePageOfObjectStoreCredentials(c.Client, CivoCmdOptions)
		i++
	}
	return nil
}

func (c *CivoConfiguration) NukeObjectStoreCredentials(CivoCmdOptions *CivoCmdOptions) error {
	pageOfCredItems, err := c.Client.ListObjectStoreCredentials()
	if err != nil {
		fmt.Println("err getting pg items")
	}
	j := 0
	for j < pageOfCredItems.Pages {
		deletePageOfObjectStoreCredentials(c.Client, CivoCmdOptions)
		j++
	}
	return nil
}

func deletePageOfObjectStores(client *civogo.Client, CivoCmdOptions *CivoCmdOptions) {

	pageOfItems, err := client.ListObjectStores()
	if err != nil {
		fmt.Println("err")
	}

	for _, os := range pageOfItems.Items {
		if CivoCmdOptions.Nuke {
			_, err := client.DeleteObjectStore(os.ID)
			if err != nil {
				fmt.Println("err")
			}
			fmt.Println("deleted object store ", os.Name)

		} else {
			fmt.Printf("nuke set to %t, not removing objectstore %s\n", CivoCmdOptions.Nuke, os.ID)
		}
	}
}

func deletePageOfObjectStoreCredentials(client *civogo.Client, CivoCmdOptions *CivoCmdOptions) {

	pageOfItems, err := client.ListObjectStoreCredentials()
	if err != nil {
		fmt.Println("err", err)
	}

	for _, os := range pageOfItems.Items {
		if CivoCmdOptions.Nuke {
			_, err := client.DeleteObjectStoreCredential(os.ID)
			if err != nil {
				fmt.Println("err", err)
			}
			fmt.Println("deleted object store credential", os.Name)

		} else {
			fmt.Printf("nuke set to %t, not removing objectstore credential %s\n", CivoCmdOptions.Nuke, os.ID)
		}
	}
}
