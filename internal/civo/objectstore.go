package civo

import (
	"errors"
	"fmt"

	"github.com/civo/civogo"
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

			c.logger.Infof("finding object store credential for object store %q", objStore.Name)

			credentialDetails, err := c.client.FindObjectStoreCredential(objStore.Name)
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					c.logger.Infof("no object store credentials for object store %q", objStore.Name)
				} else {
					return fmt.Errorf("unable to find object store credential %q: %w", objStore.Name, err)
				}
			} else {
				c.logger.Infof("found object store credential %q - ID: %q", credentialDetails.Name, credentialDetails.ID)
			}

			if c.nuke {
				c.logger.Infof("deleting object store %q", objStore.Name)

				if _, err := c.client.DeleteObjectStore(objStore.ID); err != nil {
					return fmt.Errorf("unable to delete object store %q: %w", objStore.Name, err)
				}

				outputwriter.WriteStdoutf("deleted object store %q", objStore.Name)

				if credentialDetails != nil {
					c.logger.Infof("deleting object store credential %q with ID %q", credentialDetails.Name, credentialDetails.ID)

					if _, err = c.client.DeleteObjectStoreCredential(credentialDetails.ID); err != nil {
						return fmt.Errorf("unable to delete object store credential %q (ID: %q): %w", credentialDetails.Name, credentialDetails.ID, err)
					}

					outputwriter.WriteStdoutf("deleted object store credential %q with ID %q", credentialDetails.Name, credentialDetails.ID)
				}
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
				c.logger.Infof("deleting object store credential %q", objStoreCred.Name)

				_, err := c.client.DeleteObjectStoreCredential(objStoreCred.ID)
				if err != nil {
					return fmt.Errorf("unable to delete object store credential %q: %w", objStoreCred.Name, err)
				}

				outputwriter.WriteStdoutf("deleted object store credential %q", objStoreCred.Name)
			} else {
				c.logger.Warnf("refusing to delete object store credential %q: nuke is not enabled", objStoreCred.Name)
			}
		}
	}

	return nil
}
