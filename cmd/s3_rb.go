package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdRb is the Cobra CLI call
func cliS3CmdRb() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rb CLUSTER BUCKET",
		Short: "Remove bucket",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdRb,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// S3CmdRb wraps s3cmd command in the container
func S3CmdRb(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)
	command := []string{"s3cmd", "rb", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + args[0]
	fmt.Println(output)
}
