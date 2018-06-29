package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// S3CmdRec is the option to apply when trying to delete content
	S3CmdRec bool
)

// cliS3CmdDel is the Cobra CLI call
func cliS3CmdDel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del CLUSTER BUCKET/OBJECT",
		Short: "Delete file from bucket",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdDel,
		DisableFlagsInUseLine: true,
	}
	//cmd.Flags().BoolVarP(&S3CmdRec, "recursive", "r", false, "Recursive removal.")
	//cmd.Flags().BoolVarP(&S3CmdForce, "force", "f", false, "Force removal.")

	return cmd
}

// S3CmdDel wraps s3cmd command in the container
func S3CmdDel(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]
	containerNameToShow := containerName[len(containerNamePrefix):]

	notExistCheck(containerName)
	notRunningCheck(containerName)

	/*
		S3CmdOpt = "--verbose"
			if S3CmdRec {
				S3CmdOpt = S3CmdOpt + " --recursive"
			}
			if S3CmdForce {
				S3CmdOpt = S3CmdOpt + " --force"
			}
			command := []string{"s3cmd", "del", S3CmdOpt, "s3://" + args[0]}
	*/
	command := []string{"s3cmd", "del", "s3://" + args[1]}
	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerNameToShow
	fmt.Println(output)
}
