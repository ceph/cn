package cmd

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// cliClusterStop is the Cobra CLI call
func cliClusterStop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop an object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   stopNano,
	}

	return cmd
}

// stopNano stops Ceph Nano
func stopNano(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow
	timeout := 5 * time.Second

	if status := containerStatus(containerName, true, "exited"); status {
		log.Println("Cluster " + containerNameToShow + " is already stopped.")
		os.Exit(0)
	} else if status := containerStatus(containerName, false, "running"); !status {
		log.Println("Cluster " + containerNameToShow + " does not exist yet.")
		os.Exit(0)
	} else {
		log.Println("Stopping cluster " + containerNameToShow + "...")
		if err := getDocker().ContainerStop(ctx, containerName, &timeout); err != nil {
			log.Fatal(err)
		}
	}
}
