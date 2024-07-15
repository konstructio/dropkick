/*
Copyright (C) 2021-2024, konstruct

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package cmd

import (
	"fmt"

	"github.com/konstructio/dropkick/pkg/configs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version for dropkick",
	Long:  `All software has versions. This is dropkick's`,
	Run: func(cmd *cobra.Command, args []string) {
		versionMsg := `
##
### dropkick golang utility version:` + fmt.Sprintf("`%s`", configs.Version)

		fmt.Println(versionMsg)
	},
}
