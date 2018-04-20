package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// cliClusterRestart is the Cobra CLI call
func cliClusterRestart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart an object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   restartNano,
	}

	return cmd
}

// restartNano restarts Ceph Nano
func restartNano(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]
	containerNameToShow := containerName[len(containerNamePrefix):]

	notExistCheck(containerName)
	log.Println("Restarting cluster " + containerNameToShow + "...")
	if err := getDocker().ContainerRestart(ctx, containerName, nil); err != nil {
		log.Fatal(err)
	}
	echoInfo(containerName)
}
