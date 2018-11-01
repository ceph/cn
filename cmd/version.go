package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cliVersionNano is the Cobra CLI call
func cliVersionNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of cn",
		Args:  cobra.NoArgs,
		Run:   versionNano,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// versionNano print Ceph Nano version
func versionNano(cmd *cobra.Command, args []string) {
	fmt.Println("ceph-nano version " + cnVersion)
}
