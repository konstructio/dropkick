package civo

import (
	"fmt"

	"github.com/civo/civogo"
)

func (c *CivoConfiguration) NukeObjectStores(client *civogo.Client) {

	pageOfItems, err := client.ListObjectStores()
	if err != nil {
		fmt.Println("err getting pg items")
	}
	fmt.Println(pageOfItems.Pages) // 4

	i := 0
	for i < pageOfItems.Pages {
		deletePageOfObjectStores(c.Client)
		deletePageOfObjectStoreCredentials(c.Client)
		i++
	}
}

func (c *CivoConfiguration) NukeObjectStoreCredentials(client *civogo.Client) {
	pageOfCredItems, err := client.ListObjectStoreCredentials()
	if err != nil {
		fmt.Println("err getting pg items")
	}
	j := 0
	for j < pageOfCredItems.Pages {
		deletePageOfObjectStoreCredentials(c.Client)
		j++
	}
}

func deletePageOfObjectStores(client *civogo.Client) {

	pageOfItems, err := client.ListObjectStores()
	if err != nil {
		fmt.Println("err")
	}

	for _, os := range pageOfItems.Items {
		_, err := client.DeleteObjectStore(os.ID)
		if err != nil {
			fmt.Println("err")
		}
		fmt.Println("deleted object store ", os.Name)
	}
}

func deletePageOfObjectStoreCredentials(client *civogo.Client) {

	pageOfItems, err := client.ListObjectStoreCredentials()
	if err != nil {
		fmt.Println("err", err)
	}

	for _, os := range pageOfItems.Items {
		_, err := client.DeleteObjectStoreCredential(os.ID)
		if err != nil {
			fmt.Println("err", err)
		}
		fmt.Println("deleted object store credential", os.Name)
	}
}
