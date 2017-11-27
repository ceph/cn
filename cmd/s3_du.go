package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliS3CmdDu is the Cobra CLI call
func CliS3CmdDu() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "du BUCKET/PREFIX",
		Short: "Disk usage by buckets",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdDu,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdDu wraps s3cmd command in the container
func S3CmdDu(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	command := []string{"s3cmd", "du", "s3://" + args[0]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
