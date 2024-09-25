package civo

import (
	"context"
	"errors"
	"fmt"

	"github.com/konstructio/dropkick/internal/civov2"
	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// NukeObjectStores deletes all object stores associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeObjectStores() error {
	c.logger.Infof("listing object stores")

	items, err := c.client.GetAllObjectStores(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list object stores: %w", err)
	}

	c.logger.Infof("found %d object stores", len(items))

	for _, objStore := range items {
		c.logger.Infof("found object store %q", objStore.Name)

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(objStore.Name, c.nameFilter) {
			c.logger.Warnf("skipping object store %q: name does not match filter", objStore.Name)
			continue
		}

		c.logger.Infof("finding object store credential for object store %q", objStore.Name)

		credentialDetails, err := c.client.GetObjectStoreCredential(context.Background(), objStore.ID)
		if err != nil {
			if errors.Is(err, civov2.ErrNotFound) {
				c.logger.Infof("no object store credentials for object store %q", objStore.Name)
			} else {
				return fmt.Errorf("unable to find object store credential %q: %w", objStore.Name, err)
			}
		} else {
			c.logger.Infof("found object store credential %q - ID: %q", credentialDetails.Name, credentialDetails.ID)
		}

		if c.nuke {
			c.logger.Infof("deleting object store %q", objStore.Name)

			if err := c.client.DeleteObjectStore(context.Background(), objStore.ID); err != nil {
				return fmt.Errorf("unable to delete object store %q: %w", objStore.Name, err)
			}

			outputwriter.WriteStdoutf("deleted object store %q", objStore.Name)

			if credentialDetails != nil {
				c.logger.Infof("deleting object store credential %q with ID %q", credentialDetails.Name, credentialDetails.ID)

				if err = c.client.DeleteObjectStoreCredential(context.Background(), credentialDetails.ID); err != nil {
					return fmt.Errorf("unable to delete object store credential %q (ID: %q): %w", credentialDetails.Name, credentialDetails.ID, err)
				}

				outputwriter.WriteStdoutf("deleted object store credential %q with ID %q", credentialDetails.Name, credentialDetails.ID)
			}
		} else {
			c.logger.Warnf("refusing to delete object store %q: nuke is not enabled", objStore.Name)
		}
	}

	return nil
}

// NukeObjectStoreCredentials deletes all object store credentials associated with the Civo client.
// It returns an error if the deletion process encounters any issues.
func (c *Civo) NukeObjectStoreCredentials() error {
	c.logger.Infof("listing object store credentials")

	items, err := c.client.GetAllObjectStoreCredentials(context.Background())
	if err != nil {
		return fmt.Errorf("unable to list object store credentials: %w", err)
	}

	c.logger.Infof("found %d object store credentials", len(items))

	for _, objStoreCred := range items {
		c.logger.Infof("found object store credential %q - ID: %q", objStoreCred.Name, objStoreCred.ID)

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(objStoreCred.Name, c.nameFilter) {
			c.logger.Warnf("skipping object store credential %q: name does not match filter", objStoreCred.Name)
			continue
		}

		if c.nuke {
			c.logger.Infof("deleting object store credential %q", objStoreCred.Name)

			err := c.client.DeleteObjectStoreCredential(context.Background(), objStoreCred.ID)
			if err != nil {
				return fmt.Errorf("unable to delete object store credential %q: %w", objStoreCred.Name, err)
			}

			outputwriter.WriteStdoutf("deleted object store credential %q", objStoreCred.Name)
		} else {
			c.logger.Warnf("refusing to delete object store credential %q: nuke is not enabled", objStoreCred.Name)
		}
	}

	return nil
}
