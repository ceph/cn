package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// CliS3CmdInfo is the Cobra CLI call
func CliS3CmdInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info CLUSTER BUCKET/OBJECT",
		Short: "Get various information about Buckets or Files",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdInfo,
	}

	return cmd
}

// S3CmdInfo wraps s3cmd command in the container
func S3CmdInfo(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]

	notExistCheck(ContainerName)
	notRunningCheck(ContainerName)
	command := []string{"s3cmd", "info", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(ContainerName, command)), "\n") + " on cluster " + ContainerName
	fmt.Println(output)
}
