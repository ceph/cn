package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliS3CmdRb is the Cobra CLI call
func CliS3CmdRb() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rb BUCKET",
		Short: "Remove bucket",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdRb,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdRb wraps s3cmd command in the container
func S3CmdRb(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	command := []string{"s3cmd", "rb", "s3://" + args[0]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
