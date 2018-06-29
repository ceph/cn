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
	cmd.Flags().BoolVarP(&debugS3, "debug", "d", false, "Run S3 commands in debug mode")

	return cmd
}

// S3CmdCp wraps s3cmd command in the container
func S3CmdCp(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	notRunningCheck(containerName)

	command := []string{"s3cmd", "cp", "s3://" + args[1], "s3://" + args[2]}
	if debugS3 {
		command = append(command, "--debug")
	}

	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerNameToShow
	fmt.Println(output)
}
