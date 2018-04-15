package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdMb is the Cobra CLI call
func cliS3CmdMb() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mb CLUSTER BUCKET",
		Short: "Make bucket",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdMb,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// S3CmdMb wraps s3cmd command in the container
func S3CmdMb(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)
	command := []string{"s3cmd", "mb", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerName
	fmt.Println(output)
}
