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

func getCivoCommand() *cobra.Command {
	var (
		nuke       bool
		region     string
		nameFilter string
	)

	civoCmd := &cobra.Command{
		Use:   "civo",
		Short: "clean civo resources",
		Long:  `clean civo resources`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			quiet := cmd.Flags().Lookup("quiet").Value.String() == "true"
			return runCivo(cmd.OutOrStderr(), region, os.Getenv("CIVO_TOKEN"), nameFilter, nuke, quiet)
		},
	}

	civoCmd.Flags().BoolVar(&nuke, "nuke", false, "required to confirm deletion of resources")
	civoCmd.Flags().StringVar(&region, "region", "", "the civo region to clean")
	civoCmd.Flags().StringVar(&nameFilter, "name-contains", "", "if set, only resources with a name containing this string will be selected")
	err := civoCmd.MarkFlagRequired("region")
	if err != nil {
		log.Fatal(err)
	}

	return civoCmd
}

func runCivo(output io.Writer, region, token, nameFilter string, nuke, quiet bool) error {
	if token == "" {
		return errors.New("required environment variable $CIVO_TOKEN not found: get one at https://dashboard.civo.com/security")
	}

	// Create a logger and make it quiet
	var log *logger.Logger
	if quiet {
		log = logger.New(io.Discard)
	} else {
		log = logger.New(output)
	}

	client, err := civo.New(
		civo.WithToken(token),
		civo.WithRegion(region),
		civo.WithNameFilter(nameFilter),
		civo.WithNuke(nuke),
		civo.WithLogger(log),
	)
	if err != nil {
		return fmt.Errorf("unable to create new client: %w", err)
	}

	if err := client.NukeKubernetesClusters(); err != nil {
		return fmt.Errorf("unable to cleanup Kubernetes clusters: %w", err)
	}

	if err := client.NukeObjectStores(); err != nil {
		return fmt.Errorf("unable to cleanup object stores: %w", err)
	}

	err = client.NukeObjectStoreCredentials()
	if err != nil {
		return fmt.Errorf("unable to cleanup object store credentials: %w", err)
	}

	err = client.NukeVolumes()
	if err != nil {
		return fmt.Errorf("unable to cleanup volumes: %w", err)
	}

	err = client.NukeNetworks()
	if err != nil {
		return fmt.Errorf("unable to cleanup networks: %w", err)
	}

	return nil
}
