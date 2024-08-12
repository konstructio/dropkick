package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/konstructio/dropkick/internal/digitalocean"
	"github.com/konstructio/dropkick/internal/logger"
	"github.com/konstructio/dropkick/pkg/env"
	"github.com/spf13/cobra"
)

type doOptions struct {
	nuke            bool
	token           string
	spacesAccessKey string
	spacesSecretKey string
	spacesRegion    string
}

func getDigitalOceanCommand() *cobra.Command {
	var opts doOptions

	cmd := &cobra.Command{
		Use:   "digitalocean",
		Short: "clean digitalocean resources",
		Long:  `clean digitalocean resources`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.token = env.GetFirstNotEmpty("DIGITALOCEAN_TOKEN")
			opts.spacesAccessKey = env.GetFirstNotEmpty("DIGITALOCEAN_SPACES_ACCESS_KEY", "SPACES_KEY")
			opts.spacesSecretKey = env.GetFirstNotEmpty("DIGITALOCEAN_SPACES_SECRET_KEY", "SPACES_SECRET")
			opts.spacesRegion = env.GetFirstNotEmpty("DIGITALOCEAN_SPACES_REGION", "SPACES_REGION")
			quiet := cmd.Flags().Lookup("quiet").Value.String() == "true"
			return runDigitalOcean(cmd.OutOrStderr(), opts, quiet)
		},
	}

	cmd.Flags().BoolVar(&opts.nuke, "nuke", false, "required to confirm deletion of resources")
	return cmd
}

func runDigitalOcean(output io.Writer, opts doOptions, quiet bool) error {
	// Check token
	if opts.token == "" {
		return errors.New("required environment variable $DIGITALOCEAN_TOKEN not set")
	}

	// Check spaces credentials
	if opts.spacesAccessKey == "" {
		return errors.New("required environment variable $DIGITALOCEAN_SPACES_ACCESS_KEY or $SPACES_KEY not set")
	}
	if opts.spacesSecretKey == "" {
		return errors.New("required environment variable $DIGITALOCEAN_SPACES_SECRET_KEY or $SPACES_SECRET not set")
	}
	if opts.spacesRegion == "" {
		return errors.New("required environment variable $DIGITALOCEAN_SPACES_REGION or $SPACES_REGION not set")
	}

	// Create a logger and make it quiet
	var log *logger.Logger
	if quiet {
		log = logger.New(io.Discard)
	} else {
		log = logger.New(output)
	}

	// Create DigitalOcean client
	client, err := digitalocean.New(
		digitalocean.WithToken(opts.token),
		digitalocean.WithS3Storage(opts.spacesAccessKey, opts.spacesSecretKey, opts.spacesRegion),
		digitalocean.WithNuke(opts.nuke),
		digitalocean.WithContext(context.Background()),
		digitalocean.WithLogger(log),
	)
	if err != nil {
		return fmt.Errorf("unable to create new client: %w", err)
	}

	// Cleanup resources
	if err := client.NukeKubernetesClusters(); err != nil {
		return fmt.Errorf("unable to cleanup Kubernetes clusters: %w", err)
	}

	if err := client.NukeS3Storage(); err != nil {
		return fmt.Errorf("unable to cleanup spaces storage: %w", err)
	}

	if err := client.NukeVolumes(); err != nil {
		return fmt.Errorf("unable to cleanup volumes: %w", err)
	}

	return nil
}
