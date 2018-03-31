package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// CliS3CmdMv is the Cobra CLI call
func CliS3CmdMv() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mv CLUSTER BUCKET1/OBJECT1 BUCKET2/OBJECT2",
		Short: "Move object",
		Args:  cobra.ExactArgs(3),
		Run:   S3CmdMv,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// S3CmdMv wraps s3cmd command in the container
func S3CmdMv(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]

	notExistCheck(ContainerName)
	notRunningCheck(ContainerName)
	command := []string{"s3cmd", "mv", "s3://" + args[1], "s3://" + args[2]}
	output := strings.TrimSuffix(string(execContainer(ContainerName, command)), "\n") + " on cluster " + ContainerName
	fmt.Println(output)
}
