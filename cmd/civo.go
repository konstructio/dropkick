package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/konstructio/dropkick/internal/civo"
	"github.com/konstructio/dropkick/internal/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type civoOptions struct {
	nuke        bool
	region      string
	nameFilter  string
	quiet       bool
	onlyOrphans bool
}

func getCivoCommand() *cobra.Command {
	var opts civoOptions

	civoCmd := &cobra.Command{
		Use:   "civo",
		Short: "clean civo resources",
		Long:  `clean civo resources`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.quiet = cmd.Flags().Lookup("quiet").Value.String() == "true"
			return runCivo(cmd.OutOrStderr(), opts, os.Getenv("CIVO_TOKEN"))
		},
	}

	civoCmd.Flags().BoolVar(&opts.nuke, "nuke", false, "required to confirm deletion of resources")
	civoCmd.Flags().StringVar(&opts.region, "region", "", "the civo region to clean")
	civoCmd.Flags().StringVar(&opts.nameFilter, "name-contains", "", "if set, only resources with a name containing this string will be selected")
	civoCmd.Flags().BoolVar(&opts.onlyOrphans, "orphans-only", false, "only delete orphaned resources (only volumes, SSH keys and networks)")

	// On orphaned resources, we don't want to filter by name since the
	// filter is already that's just for the resources we want to delete
	civoCmd.MarkFlagsMutuallyExclusive("name-contains", "orphans-only")

	if err := civoCmd.MarkFlagRequired("region"); err != nil {
		log.Fatal(err)
	}

	return civoCmd
}

func runCivo(output io.Writer, opts civoOptions, token string) error {
	if token == "" {
		return errors.New("required environment variable $CIVO_TOKEN not found: get one at https://dashboard.civo.com/security")
	}

	// Create a logger and make it quiet
	var log *logger.Logger
	if opts.quiet {
		log = logger.New(io.Discard)
	} else {
		log = logger.New(output)
	}

	client, err := civo.New(
		civo.WithToken(token),
		civo.WithRegion(opts.region),
		civo.WithNameFilter(opts.nameFilter),
		civo.WithNuke(opts.nuke),
		civo.WithLogger(log),
	)
	if err != nil {
		return fmt.Errorf("unable to create new client: %w", err)
	}

	if opts.onlyOrphans {
		if err := client.NukeOrphanedResources(); err != nil {
			return fmt.Errorf("unable to nuke orphaned resources: %w", err)
		}
		return nil
	}

	if err := client.NukeKubernetesClusters(); err != nil {
		return fmt.Errorf("unable to cleanup Kubernetes clusters: %w", err)
	}

	if err := client.NukeObjectStores(); err != nil {
		return fmt.Errorf("unable to cleanup object stores: %w", err)
	}

	if err := client.NukeObjectStoreCredentials(); err != nil {
		return fmt.Errorf("unable to cleanup object store credentials: %w", err)
	}

	if err := client.NukeVolumes(); err != nil {
		return fmt.Errorf("unable to cleanup volumes: %w", err)
	}

	if err := client.NukeNetworks(); err != nil {
		return fmt.Errorf("unable to cleanup networks: %w", err)
	}

	if err := client.NukeSSHKeys(); err != nil {
		return fmt.Errorf("unable to cleanup SSH keys: %w", err)
	}

	if err := client.NukeFirewalls(); err != nil {
		return fmt.Errorf("unable to cleanup firewalls: %w", err)
	}

	return nil
}
