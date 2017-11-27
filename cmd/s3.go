package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cmdS3 = &cobra.Command{
		Use:   "s3 [command] [arg]",
		Short: "Interact with S3 object server",
		Args:  cobra.NoArgs,
	}
	// S3CmdForce means force operation
	S3CmdForce bool
)

func init() {
	cmdS3.AddCommand(
		CliS3CmdMb(),
		CliS3CmdRb(),
		CliS3CmdLs(),
		CliS3CmdLa(),
		CliS3CmdPut(),
		CliS3CmdGet(),
		CliS3CmdDel(),
		CliS3CmdDu(),
		CliS3CmdInfo(),
		CliS3CmdCp(),
		CliS3CmdMv(),
		CliS3CmdSync())
}
