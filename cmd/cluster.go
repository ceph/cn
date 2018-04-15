package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cmdCluster = &cobra.Command{
		Use:   "cluster [command] [arg]",
		Short: "Interact with a particular Ceph cluster",
		Args:  cobra.NoArgs,
	}
)

func init() {
	cmdCluster.AddCommand(
		cliClusterList(),
		cliClusterStart(),
		cliClusterStatus(),
		cliClusterStop(),
		cliClusterRestart(),
		cliClusterLogs(),
		cliClusterPurge(),
	)
}
