package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliS3CmdMb is the Cobra CLI call
func CliS3CmdMb() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mb BUCKET",
		Short: "Make bucket",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdMb,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdMb wraps s3cmd command in the container
func S3CmdMb(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	command := []string{"s3cmd", "mb", "s3://" + args[0]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
