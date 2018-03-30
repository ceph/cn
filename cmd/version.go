package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliVersionNano is the Cobra CLI call
func CliVersionNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Ceph Nano",
		Args:  cobra.NoArgs,
		Run:   versionNano,
	}

	return cmd
}

// versionNano print Ceph Nano version
func versionNano(cmd *cobra.Command, args []string) {
	fmt.Println("ceph-nano version " + cnVersion)
	if status := containerStatus(true, "exited"); status {
		os.Exit(0)
	}
	if status := containerStatus(false, "running"); !status {
		os.Exit(0)
	}
	if ii := inspectImage(); ii["head"] == "unknown" {
		fmt.Println("ceph-nano container image version is unknown (no image pulled yet)")
	} else if len(ii["head"]) == 0 {
		fmt.Println("ceph-nano container image version is unknown (you're likely running an old image, one that doesn't have a commit label)")
	} else {
		fmt.Println("ceph-nano container image version " + ii["head"])
	}

}
