package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
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
	notExistCheck()
	fmt.Println("Restarting ceph-nano...")
	if err := getDocker().ContainerRestart(ctx, ContainerName, nil); err != nil {
		log.Fatal(err)
	}
	echoInfo()
}
