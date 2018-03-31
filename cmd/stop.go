package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// CliClusterStop is the Cobra CLI call
func CliClusterStop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   stopNano,
	}

	return cmd
}

// stopNano stops Ceph Nano
func stopNano(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]
	ContainerNameToShow := ContainerName[len(ContainerNamePrefix):]
	timeout := 5 * time.Second

	if status := containerStatus(ContainerName, true, "exited"); status {
		fmt.Println("Cluster " + ContainerNameToShow + " is already stopped.")
		os.Exit(0)
	} else if status := containerStatus(ContainerName, false, "running"); !status {
		fmt.Println("Cluster " + ContainerNameToShow + " does not exist yet.")
		os.Exit(0)
	} else {
		fmt.Println("Stopping cluster " + ContainerNameToShow + "...")
		if err := getDocker().ContainerStop(ctx, ContainerName, &timeout); err != nil {
			log.Fatal(err)
		}
	}
}
