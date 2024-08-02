package civo

import (
	"fmt"
)

// NukeObjectStores deletes all object stores associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeObjectStores() error {
	totalPages := 1
	for page := 1; page <= totalPages; page++ {
		items, err := c.client.ListObjectStores()
		if err != nil {
			return fmt.Errorf("unable to list object stores: %w", err)
		}

		totalPages = items.Pages

		for _, objStore := range items.Items {
			if c.nuke {
				status, err := c.client.DeleteObjectStore(objStore.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store %q: %w", objStore.ID, err)
				}
				fmt.Printf("deleted object store %q\n", objStore.ID)

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %s when deleting object store %q: %s", status.ErrorCode, objStore.ID, status.ErrorDetails)
				}

				status, err = c.client.DeleteObjectStoreCredential(objStore.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store credential %q: %w", objStore.ID, err)
				}
				fmt.Printf("deleted object store credential %q\n", objStore.ID)

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %s when deleting object store credential %q: %s", status.ErrorCode, objStore.ID, status.ErrorDetails)
				}
			} else {
				fmt.Printf("refusing to delete object store %q: nuke is not enabled\n", objStore.ID)
			}
		}
	}

	return nil
}

// NukeObjectStoreCredentials deletes all object store credentials associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeObjectStoreCredentials() error {
	totalPages := 1
	for page := 1; page <= totalPages; page++ {
		items, err := c.client.ListObjectStoreCredentials()
		if err != nil {
			return fmt.Errorf("unable to list object store credentials: %w", err)
		}

		totalPages = items.Pages

		for _, objStoreCred := range items.Items {
			if c.nuke {
				status, err := c.client.DeleteObjectStoreCredential(objStoreCred.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store credential %q: %w", objStoreCred.ID, err)
				}
				fmt.Printf("deleted object store credential %q\n", objStoreCred.ID)

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %s when deleting object store credential %q: %s", status.ErrorCode, objStoreCred.ID, status.ErrorDetails)
				}

			} else {
				fmt.Printf("refusing to delete object store credential %q: nuke is not enabled\n", objStoreCred.ID)
			}
		}
	}

	return nil
}
