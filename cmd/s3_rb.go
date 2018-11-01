package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdRb is the Cobra CLI call
func cliS3CmdRb() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rb [CLUSTER] [BUCKET]",
		Short: "Remove bucket",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdRb,
		DisableFlagsInUseLine: true,
	}
	cmd.Flags().BoolVarP(&debugS3, "debug", "d", false, "Run S3 commands in debug mode")

	return cmd
}

// S3CmdRb wraps s3cmd command in the container
func S3CmdRb(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	notRunningCheck(containerName)
	command := []string{"s3cmd", "rb", "s3://" + args[1]}
	if debugS3 {
		command = append(command, "--debug")
	}

	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerNameToShow
	fmt.Println(output)
}
