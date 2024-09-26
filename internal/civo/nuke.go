package civo

import (
	"context"
	"fmt"

	"github.com/konstructio/dropkick/internal/civo/sdk"
	"github.com/konstructio/dropkick/internal/compare"
	"github.com/konstructio/dropkick/internal/outputwriter"
)

func (c *Civo) NukeEverything() error {
	// The order in which these resources are deleted matter. We start with those
	// resources that have dependencies. In Civo, certain resources won't delete
	// their dependencies (for example, deleting an Instance that is on a Network
	// won't delete the Network because it could be shared with other Instances).
	// So we delete first those resources that could cascade delete other resources.

	// We start by deleting all Load Balancers, which tie to Firewalls and Kubernetes Clusters.
	if err := nukeAllResources(context.Background(), c, sdk.LoadBalancer{}); err != nil {
		return err
	}

	// Kubernetes clusters depend on volumes (PVCs), but we will delete the PVCs
	// after the clusters to clean them all.
	if err := nukeAllResources(context.Background(), c, sdk.KubernetesCluster{}); err != nil {
		return err
	}

	// Then we delete the instances, which might also have volumes attached.
	if err := nukeAllResources(context.Background(), c, sdk.Instance{}); err != nil {
		return err
	}

	// Now we delete the volumes.
	if err := nukeAllResources(context.Background(), c, sdk.Volume{}); err != nil {
		return err
	}

	// And we also delete the SSH keys, including those now orphaned by instances.
	if err := nukeAllResources(context.Background(), c, sdk.SSHKey{}); err != nil {
		return err
	}

	// Then we delete object stores, which will leave their credentials orphaned.
	if err := nukeAllResources(context.Background(), c, sdk.ObjectStore{}); err != nil {
		return err
	}

	// Now we delete the object store credentials.
	if err := nukeAllResources(context.Background(), c, sdk.ObjectStoreCredential{}); err != nil {
		return err
	}

	// Firewalls are deleted next, since we need them to be gone before deleting networks.
	if err := nukeAllResources(context.Background(), c, sdk.Firewall{}); err != nil {
		return err
	}

	// And finally, we delete the networks.
	if err := nukeAllResources(context.Background(), c, sdk.Network{}); err != nil {
		return err
	}

	return nil
}

// nukeAllResources deletes all resources of a given type associated with the Civo client.
func nukeAllResources[T sdk.Resource](ctx context.Context, civo *Civo, resource T) error {
	civo.logger.Infof("listing %ss", resource.GetResourceType())

	resources, err := sdk.GetAll[T](ctx, civo.client)
	if err != nil {
		return fmt.Errorf("unable to list %s: %w", resource.GetResourceType(), err)
	}

	civo.logger.Infof("found %d %ss", len(resources), resource.GetResourceType())

	for _, r := range resources {
		err := nukeIndividual(ctx, civo, r)
		if err != nil {
			return err
		}
	}

	return nil
}

// nukeSlice receives a slice of resources and deletes them all.
func nukeSlice[T sdk.Resource](ctx context.Context, civo *Civo, resources []T) error {
	for _, resource := range resources {
		err := nukeIndividual(ctx, civo, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

// nukeIndividual deletes a single resource.
func nukeIndividual[T sdk.Resource](ctx context.Context, civo *Civo, resource T) error {
	if resource.GetID() == "" {
		return fmt.Errorf("the \"ID\" field in the resource %s is empty", resource.GetResourceType())
	}

	if resource.GetName() == "" {
		return fmt.Errorf("the \"name\" field in the resource %s is empty", resource.GetResourceType())
	}

	civo.logger.Infof("found %s: name: %q - ID: %q", resource.GetResourceType(), resource.GetName(), resource.GetID())

	if civo.nameFilter != "" && !compare.ContainsIgnoreCase(resource.GetName(), civo.nameFilter) {
		civo.logger.Warnf("skipping %s %q: name does not match filter", resource.GetResourceType(), resource.GetName())
		return nil
	}

	if !civo.nuke {
		civo.logger.Warnf("refusing to delete %s %q: nuke is not enabled", resource.GetResourceType(), resource.GetName())
		return nil
	}

	civo.logger.Infof("deleting %s %q", resource.GetResourceType(), resource.GetName())

	err := sdk.Delete(ctx, civo.client, resource)
	if err != nil {
		return fmt.Errorf("unable to delete %s %q: %w", resource.GetResourceType(), resource.GetName(), err)
	}

	outputwriter.WriteStdoutf("deleted %s %q", resource.GetResourceType(), resource.GetName())
	return nil
}
