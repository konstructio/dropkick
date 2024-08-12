package civo

import (
	"fmt"

	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeObjectStores deletes all object stores associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeObjectStores() error {
	totalPages := 1
	for page := 1; page <= totalPages; page++ {
		c.logger.Printf("listing object stores for page %d", page)

		items, err := c.client.ListObjectStores()
		if err != nil {
			return fmt.Errorf("unable to list object stores: %w", err)
		}

		c.logger.Printf("found %d object stores and %d pages", len(items.Items), items.Pages)

		totalPages = items.Pages

		for _, objStore := range items.Items {
			c.logger.Printf("found object store %q", objStore.ID)

			if c.nuke {
				c.logger.Printf("deleting object store %q", objStore.ID)
				status, err := c.client.DeleteObjectStore(objStore.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store %q: %w", objStore.ID, err)
				}

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %s when deleting object store %q: %s", status.ErrorCode, objStore.ID, status.ErrorDetails)
				}
				outputwriter.WriteStdoutf("deleted object store %q", objStore.ID)

				c.logger.Printf("deleting object store credential %q", objStore.ID)
				status, err = c.client.DeleteObjectStoreCredential(objStore.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store credential %q: %w", objStore.ID, err)
				}

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %s when deleting object store credential %q: %s", status.ErrorCode, objStore.ID, status.ErrorDetails)
				}
				outputwriter.WriteStdoutf("deleted object store credential %q", objStore.ID)
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

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %s when deleting object store credential %q: %s", status.ErrorCode, objStoreCred.ID, status.ErrorDetails)
				}
				outputwriter.WriteStdoutf("deleted object store credential %q", objStoreCred.ID)
			} else {
				fmt.Printf("refusing to delete object store credential %q: nuke is not enabled\n", objStoreCred.ID)
			}
		}
	}

	return nil
}
