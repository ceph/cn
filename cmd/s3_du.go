package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// CliS3CmdDu is the Cobra CLI call
func CliS3CmdDu() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "du CLUSTER BUCKET/PREFIX",
		Short: "Disk usage by buckets",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdDu,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// S3CmdDu wraps s3cmd command in the container
func S3CmdDu(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]

	notExistCheck(ContainerName)
	notRunningCheck(ContainerName)
	command := []string{"s3cmd", "du", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(ContainerName, command)), "\n") + " on cluster " + ContainerName
	fmt.Println(output)
}
