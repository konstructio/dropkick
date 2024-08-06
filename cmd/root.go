/*
Copyright © 2024 konstruct info@konstructio.io
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:           "dropkick",
		Short:         "A brief description of your application",
		Long:          ``,
		SilenceUsage:  true, // prevents printing usage when an error occurs
		SilenceErrors: true, // we want to print the error ourselves
	}

	// Add subcommands
	rootCmd.AddCommand(getCivoCommand())
	rootCmd.AddCommand(getDigitalOceanCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
