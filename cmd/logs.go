package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cliClusterLogs is the Cobra CLI call
func cliClusterLogs() *cobra.Command {
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
	containerName := containerNamePrefix + args[0]
	showS3Logs(containerName)
}

func showS3Logs(containerName string) {
	notExistCheck(containerName)
	c := []string{"cat", "/var/log/ceph/client.rgw." + containerName + "-faa32aebf00b.log"}
	output := execContainer(containerName, c)
	fmt.Printf("%s", output)
}
