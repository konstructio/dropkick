package cmd

import (
	"context"
	"os"

	"github.com/konstructio/dropkick/internal/civo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var CivoCmdOptions *civo.CivoCmdOptions = &civo.CivoCmdOptions{}

var civoCmd = &cobra.Command{
	Use:   "civo",
	Short: "clean civo resources",
	Long:  `clean civo resources`,
	Run: func(cmd *cobra.Command, args []string) {

		if os.Getenv("CIVO_TOKEN") == "" {
			log.Fatal("no civoCmd token present")
		}

		civoConf := civo.CivoConfiguration{
			Client:  civo.NewClient(os.Getenv("CIVO_TOKEN"), CivoCmdOptions.Region),
			Context: context.Background(),
		}

		err := civoConf.NukeKubernetesClusters(CivoCmdOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = civoConf.NukeObjectStores(CivoCmdOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = civoConf.NukeObjectStoreCredentials(CivoCmdOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = civoConf.NukeVolumes(CivoCmdOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = civoConf.NukeNetworks(CivoCmdOptions)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(civoCmd)
	civoCmd.Flags().BoolVar(&CivoCmdOptions.Nuke, "nuke", CivoCmdOptions.Nuke, "required to confirm deletion of resources")

	civoCmd.Flags().StringVar(&CivoCmdOptions.Region, "region", CivoCmdOptions.Region, "the civo region to clean")
	err := civoCmd.MarkFlagRequired("region")
	if err != nil {
		log.Fatal(err)
	}
}
