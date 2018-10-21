package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdDu is the Cobra CLI call
func cliS3CmdDu() *cobra.Command {
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
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)
	command := []string{"s3cmd", "du", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + args[0]
	fmt.Println(output)
}
