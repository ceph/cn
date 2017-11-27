package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliLogsNano is the Cobra CLI call
func CliLogsNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Print object storage server logs",
		Args:  cobra.NoArgs,
		Run:   logsNano,
	}
	return cmd
}

// logsNano prints rgw logs
func logsNano(cmd *cobra.Command, args []string) {
	showS3Logs()
}

func showS3Logs() {
	notExistCheck()
	c := []string{"cat", "/var/log/ceph/client.rgw.ceph-nano-faa32aebf00b.log"}
	output := execContainer(ContainerName, c)
	fmt.Printf("%s", output)
}
