package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/civo/sdk"
	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

// deleteIterator returns a function that can be used to iterate over resources.
func (c *Civo) deleteIterator(ctx context.Context) func(sdk.APIResource) error {
	return func(resource sdk.APIResource) error {
		c.logger.Infof("found %s: name: %q - ID: %q", resource.GetResourceType(), resource.GetName(), resource.GetID())

		if c.nameFilter != "" && !compare.ContainsIgnoreCase(resource.GetName(), c.nameFilter) {
			c.logger.Warnf("skipping %s %q: name does not match filter", resource.GetResourceType(), resource.GetName())
			return nil
		}

		if !c.nuke {
			c.logger.Warnf("refusing to delete %s %q: nuke is not enabled", resource.GetResourceType(), resource.GetName())
			return nil
		}

		c.logger.Infof("deleting %s %q", resource.GetResourceType(), resource.GetName())

		err := c.client.Delete(ctx, resource)
		if err != nil {
			return fmt.Errorf("unable to delete %s %q: %w", resource.GetResourceType(), resource.GetName(), err)
		}

		outputwriter.WriteStdoutf("deleted %s %q", resource.GetResourceType(), resource.GetName())
		return nil
	}
}
