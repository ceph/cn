package cmd

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// CliRestartNano is the Cobra CLI call
func CliRestartNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart object storage server",
		Args:  cobra.NoArgs,
		Run:   restartNano,
	}
	return cmd
}

// restartNano restarts Ceph Nano
func restartNano(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	notExistCheck()
	fmt.Println("Restarting ceph-nano...")
	if err := cli.ContainerRestart(ctx, ContainerName, nil); err != nil {
		panic(err)
	}
	echoInfo()
}
