package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// CliS3CmdLs is the Cobra CLI call
func CliS3CmdLs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls CLUSTER BUCKET",
		Short: "List objects or buckets",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdLs,
	}

	return cmd
}

// S3CmdLs wraps s3cmd command in the container
func S3CmdLs(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]

	notExistCheck(ContainerName)
	notRunningCheck(ContainerName)
	command := []string{"s3cmd", "ls", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(ContainerName, command)), "\n") + " on cluster " + ContainerName
	fmt.Println(output)
}
