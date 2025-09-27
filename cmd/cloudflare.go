package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/konstructio/dropkick/internal/cloudflare"
	"github.com/konstructio/dropkick/internal/logger"
	"github.com/spf13/cobra"
)

type cloudflareOptions struct {
	nuke      bool
	domain    string
	subdomain string
	quiet     bool
}

func getCloudflareCommand() *cobra.Command {
	var opts cloudflareOptions

	cloudflareCmd := &cobra.Command{
		Use:   "cloudflare",
		Short: "clean cloudflare dns resources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts.quiet = cmd.Flags().Lookup("quiet").Value.String() == "true"
			return runCloudflare(cmd.Context(), cmd.OutOrStderr(), opts, os.Getenv("CLOUDFLARE_API_TOKEN"))
		},
	}

	cloudflareCmd.Flags().BoolVar(&opts.nuke, "nuke", false, "required to confirm deletion of resources")
	cloudflareCmd.Flags().StringVar(&opts.domain, "domain", "", "the cloudflare apex domain to clean")
	cloudflareCmd.Flags().StringVar(&opts.subdomain, "subdomain", "", "the subdomain to clean")

	cloudflareCmd.MarkFlagRequired("domain")

	return cloudflareCmd
}

func runCloudflare(ctx context.Context, output io.Writer, opts cloudflareOptions, token string) error {
	if token == "" {
		return errors.New("CLOUDFLARE_API_TOKEN environment variable not found: get an api key here https://dash.cloudflare.com/profile/api-tokens")
	}

	// Create a logger and make it quiet
	var log *logger.Logger
	if opts.quiet {
		log = logger.New(io.Discard)
	} else {
		log = logger.New(output)
	}

	client, err := cloudflare.New(
		cloudflare.WithToken(token),
		cloudflare.WithZoneName(opts.domain),
		cloudflare.WithSubdomain(opts.subdomain),
		cloudflare.WithNuke(opts.nuke),
		cloudflare.WithLogger(log),
	)
	if err != nil {
		return fmt.Errorf("unable to create new cloudflare client: %w", err)
	}

	if err := client.NukeDNSRecords(ctx); err != nil {
		if opts.nuke {
			return fmt.Errorf("unable to nuke resources: %w", err)
		}

		return fmt.Errorf("unable to process resources: %w", err)
	}

	return nil
}
