package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdCp is the Cobra CLI call
func cliS3CmdCp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cp CLUSTER BUCKET1/OBJECT1 BUCKET2/OBJECT2",
		Short: "Copy object",
		Args:  cobra.ExactArgs(3),
		Run:   S3CmdCp,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// S3CmdCp wraps s3cmd command in the container
func S3CmdCp(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)

	command := []string{"s3cmd", "cp", "s3://" + args[1], "s3://" + args[2]}
	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + args[0]
	fmt.Println(output)
}
