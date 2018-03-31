package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliClusterLogs is the Cobra CLI call
func CliClusterLogs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Print object storage server logs",
		Args:  cobra.ExactArgs(1),
		Run:   logsNano,
	}

	return cmd
}

// logsNano prints rgw logs
func logsNano(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]
	showS3Logs(ContainerName)
}

func showS3Logs(ContainerName string) {
	notExistCheck(ContainerName)
	c := []string{"cat", "/var/log/ceph/client.rgw." + ContainerName + "-faa32aebf00b.log"}
	output := execContainer(ContainerName, c)
	fmt.Printf("%s", output)
}
