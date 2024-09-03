package civo

import (
	"fmt"

	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeObjectStores deletes all object stores associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeObjectStores() error {
	totalPages := 1
	for page := 1; page <= totalPages; page++ {
		c.logger.Infof("listing object stores for page %d", page)

		items, err := c.client.ListObjectStores()
		if err != nil {
			return fmt.Errorf("unable to list object stores: %w", err)
		}

		c.logger.Infof("found %d object stores and %d pages", len(items.Items), items.Pages)

		totalPages = items.Pages

		for _, objStore := range items.Items {
			c.logger.Infof("found object store %q", objStore.Name)

			if c.nameFilter != "" && !compare.ContainsIgnoreCase(objStore.Name, c.nameFilter) {
				c.logger.Warnf("skipping object store %q: name does not match filter", objStore.Name)
				continue
			}

			if c.nuke {
				c.logger.Infof("deleting object store %q", objStore.Name)
				status, err := c.client.DeleteObjectStore(objStore.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store %q: %w", objStore.Name, err)
				}

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %q when deleting object store %q: %s", status.ErrorCode, objStore.Name, status.ErrorDetails)
				}
				outputwriter.WriteStdoutf("deleted object store %q", objStore.Name)

				c.logger.Infof("deleting object store credential %q", objStore.Name)
				status, err = c.client.DeleteObjectStoreCredential(objStore.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store credential %q: %w", objStore.Name, err)
				}

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %q when deleting object store credential %q: %s", status.ErrorCode, objStore.Name, status.ErrorDetails)
				}
				outputwriter.WriteStdoutf("deleted object store credential %q", objStore.Name)
			} else {
				c.logger.Warnf("refusing to delete object store %q: nuke is not enabled", objStore.Name)
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
		c.logger.Infof("listing object store credentials for page %d", page)

		items, err := c.client.ListObjectStoreCredentials()
		if err != nil {
			return fmt.Errorf("unable to list object store credentials: %w", err)
		}

		totalPages = items.Pages
		c.logger.Infof("found %d object store credentials on page %d (%d pages total)", len(items.Items), page, items.Pages)

		for _, objStoreCred := range items.Items {
			c.logger.Infof("found object store credential %q - ID: %q", objStoreCred.Name, objStoreCred.ID)

			if c.nameFilter != "" && !compare.ContainsIgnoreCase(objStoreCred.Name, c.nameFilter) {
				c.logger.Warnf("skipping object store credential %q: name does not match filter", objStoreCred.Name)
				continue
			}

			if c.nuke {
				status, err := c.client.DeleteObjectStoreCredential(objStoreCred.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store credential %q: %w", objStoreCred.Name, err)
				}

				if status.ErrorCode != "200" {
					return fmt.Errorf("Civo returned an error code %q when deleting object store credential %q: %s", status.ErrorCode, objStoreCred.Name, status.ErrorDetails)
				}
				outputwriter.WriteStdoutf("deleted object store credential %q", objStoreCred.Name)
			} else {
				c.logger.Warnf("refusing to delete object store credential %q: nuke is not enabled", objStoreCred.Name)
			}
		}
	}

	return nil
}
