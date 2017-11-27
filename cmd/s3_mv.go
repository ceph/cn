package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliS3CmdMv is the Cobra CLI call
func CliS3CmdMv() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mv BUCKET1/OBJECT1 BUCKET2/OBJECT2",
		Short: "Move object",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdMv,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdMv wraps s3cmd command in the container
func S3CmdMv(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	command := []string{"s3cmd", "mv", "s3://" + args[0], "s3://" + args[1]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
