package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliS3CmdInfo is the Cobra CLI call
func CliS3CmdInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info BUCKET/OBJECT",
		Short: "Get various information about Buckets or Files",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdInfo,
	}
	return cmd
}

// S3CmdInfo wraps s3cmd command in the container
func S3CmdInfo(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	command := []string{"s3cmd", "info", "s3://" + args[0]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
