package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// CliStopNano is the Cobra CLI call
func CliStopNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop object storage server",
		Args:  cobra.NoArgs,
		Run:   stopNano,
	}
	return cmd
}

// stopNano stops Ceph Nano
func stopNano(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	timeout := 5 * time.Second
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is already stopped.")
		os.Exit(0)
	} else if status := containerStatus(false, "running"); !status {
		fmt.Println("ceph-nano does not exist yet.")
		os.Exit(0)
	} else {
		fmt.Println("Stopping ceph-nano... ")
		if err := cli.ContainerStop(ctx, ContainerName, &timeout); err != nil {
			panic(err)
		}
	}
}
