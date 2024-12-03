package cmd

import (
	"io"
	"os"

	"github.com/konstructio/dropkick/configs"
	"github.com/konstructio/dropkick/internal/logger"
	"github.com/spf13/cobra"
)

func getVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print the version for dropkick cli",
		Long:  `print the version for dropkick cli`,
		RunE: func(_ *cobra.Command, _ []string) error {

			log := logger.New(io.Writer(os.Stdout))

			log.Infof("dropkick cli version: %q", configs.Version)
			return nil
		},
	}
	return cmd
}
