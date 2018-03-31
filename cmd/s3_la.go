package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// CliS3CmdLa is the Cobra CLI call
func CliS3CmdLa() *cobra.Command {
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
	ContainerName := ContainerNamePrefix + args[0]

	notExistCheck(ContainerName)
	notRunningCheck(ContainerName)
	command := []string{"s3cmd", "la"}
	output := strings.TrimSuffix(string(execContainer(ContainerName, command)), "\n") + " on cluster " + ContainerName
	if len(output) == 1 {
		command := []string{"s3cmd", "ls"}
		output := strings.TrimSuffix(string(execContainer(ContainerName, command)), "\n") + " on cluster " + ContainerName
		fmt.Println(output)
	} else {
		fmt.Println(output)
	}
}
