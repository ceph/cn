package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliS3CmdCp is the Cobra CLI call
func CliS3CmdCp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cp BUCKET1/OBJECT1 BUCKET2/OBJECT2",
		Short: "Copy object",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdCp,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdCp wraps s3cmd command in the container
func S3CmdCp(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()

	command := []string{"s3cmd", "cp", "s3://" + args[0], "s3://" + args[1]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
