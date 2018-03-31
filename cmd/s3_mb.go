package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// CliS3CmdMb is the Cobra CLI call
func CliS3CmdMb() *cobra.Command {
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
	ContainerName := ContainerNamePrefix + args[0]

	notExistCheck(ContainerName)
	notRunningCheck(ContainerName)
	command := []string{"s3cmd", "mb", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(ContainerName, command)), "\n") + " on cluster " + ContainerName
	fmt.Println(output)
}
