package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cmdS3 = &cobra.Command{
		Use:   "s3 [command] [arg]",
		Short: "Interact with a particular S3 object server",
		Args:  cobra.NoArgs,
	}
	// S3CmdForce means force operation
	S3CmdForce bool

	// debugS3 means use the '--debug' flag in the s3cmd command
	debugS3 bool
)

func init() {
	cmdS3.AddCommand(
		cliS3CmdMb(),
		cliS3CmdRb(),
		cliS3CmdLs(),
		cliS3CmdLa(),
		cliS3CmdPut(),
		cliS3CmdGet(),
		cliS3CmdDel(),
		cliS3CmdDu(),
		cliS3CmdInfo(),
		cliS3CmdCp(),
		cliS3CmdMv(),
		cliS3CmdSync())
}
