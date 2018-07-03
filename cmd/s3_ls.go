package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cliS3CmdLs is the Cobra CLI call
func cliS3CmdLs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls CLUSTER [BUCKET]",
		Short: "List objects or buckets",
		Args:  cobra.RangeArgs(1, 2),
		Run:   S3CmdLs,
	}

	return cmd
}

// S3CmdLs wraps s3cmd command in the container
func S3CmdLs(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)

	var command []string

	if len(args) == 1 {
		command = []string{"s3cmd", "ls"}
	} else {
		command = []string{"s3cmd", "ls", "s3://" + args[1]}
	}
	fmt.Println(execContainer(containerName, command))
}
