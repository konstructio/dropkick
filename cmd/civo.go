package cmd

import (
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
		nuke   bool
		region string
	)

	civoCmd := &cobra.Command{
		Use:   "civo",
		Short: "clean civo resources",
		Long:  `clean civo resources`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCivo(cmd.OutOrStdout(), region, nuke)
		},
	}

	civoCmd.Flags().BoolVar(&nuke, "nuke", false, "required to confirm deletion of resources")
	civoCmd.Flags().StringVar(&region, "region", "", "the civo region to clean")
	err := civoCmd.MarkFlagRequired("region")
	if err != nil {
		log.Fatal(err)
	}

	return civoCmd
}

func runCivo(output io.Writer, region string, nuke bool) error {
	token := os.Getenv("CIVO_TOKEN")
	if token == "" {
		return fmt.Errorf("required environment variable $CIVO_TOKEN not found")
	}

	client, err := civo.New(
		civo.WithToken(token),
		civo.WithRegion(region),
		civo.WithNuke(nuke),
		civo.WithLogger(logger.New(output)),
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
