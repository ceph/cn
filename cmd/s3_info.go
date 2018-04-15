package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdInfo is the Cobra CLI call
func cliS3CmdInfo() *cobra.Command {
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
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)
	command := []string{"s3cmd", "info", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerName
	fmt.Println(output)
}
