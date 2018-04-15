package cmd

import (
	"fmt"
	"strings"

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
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)
	command := []string{"s3cmd", "la"}
	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerName
	if len(output) == 1 {
		command := []string{"s3cmd", "ls"}
		output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerName
		fmt.Println(output)
	} else {
		fmt.Println(output)
	}
}
