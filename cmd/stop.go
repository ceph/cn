package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
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
	timeout := 5 * time.Second
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is already stopped.")
		os.Exit(0)
	} else if status := containerStatus(false, "running"); !status {
		fmt.Println("ceph-nano does not exist yet.")
		os.Exit(0)
	} else {
		fmt.Println("Stopping ceph-nano... ")
		if err := getDocker().ContainerStop(ctx, ContainerName, &timeout); err != nil {
			log.Fatal(err)
		}
	}
}
