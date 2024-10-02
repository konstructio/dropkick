package cmd

import (
	"context"
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
			return runCivo(cmd.Context(), cmd.OutOrStderr(), opts, os.Getenv("CIVO_TOKEN"))
		},
	}

	civoCmd.Flags().BoolVar(&opts.nuke, "nuke", false, "required to confirm deletion of resources")
	civoCmd.Flags().StringVar(&opts.region, "region", "", "the civo region to clean")
	civoCmd.Flags().StringVar(&opts.nameFilter, "name-contains", "", "if set, only resources with a name containing this string will be selected")
	civoCmd.Flags().BoolVar(&opts.onlyOrphans, "orphans-only", false, "only delete orphaned resources (only load balancers, volumes, object store credentials, SSH keys, networks and firewalls)")

	if err := civoCmd.MarkFlagRequired("region"); err != nil {
		log.Fatal(err)
	}

	return civoCmd
}

func runCivo(ctx context.Context, output io.Writer, opts civoOptions, token string) error {
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
		civo.WithContext(ctx),
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
		if err := client.NukeOrphanedResources(ctx); err != nil {
			return fmt.Errorf("unable to nuke orphaned resources: %w", err)
		}
		return nil
	}

	if err := client.NukeEverything(ctx); err != nil {
		if opts.nuke {
			return fmt.Errorf("unable to nuke resources: %w", err)
		}

		return fmt.Errorf("unable to process resources: %w", err)
	}

	return nil
}
