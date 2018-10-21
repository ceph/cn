package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cliS3CmdLa is the Cobra CLI call
func cliS3CmdLa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "la CLUSTER",
		Short: "List all object in all buckets",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdLa,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// S3CmdLa wraps s3cmd command in the container
func S3CmdLa(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	notRunningCheck(containerName)

	command := []string{"s3cmd", "la"}
	output := string(execContainer(containerName, command))

	if len(output) == 1 {
		command = []string{"s3cmd", "ls"}
		output = execContainer(containerName, command)
	}
	fmt.Println(output)
}
