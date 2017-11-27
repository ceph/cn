package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliS3CmdLa is the Cobra CLI call
func CliS3CmdLa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "la",
		Short: "List all object in all buckets",
		Args:  cobra.NoArgs,
		Run:   S3CmdLa,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdLa wraps s3cmd command in the container
func S3CmdLa(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	command := []string{"s3cmd", "la"}
	output := execContainer(ContainerName, command)
	if len(output) == 1 {
		command := []string{"s3cmd", "ls"}
		o := execContainer(ContainerName, command)
		fmt.Printf("%s", o)
	} else {
		fmt.Printf("%s", output)
	}
}
