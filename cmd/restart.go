package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// CliClusterRestart is the Cobra CLI call
func CliClusterRestart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   restartNano,
	}

	return cmd
}

// restartNano restarts Ceph Nano
func restartNano(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]
	ContainerNameToShow := ContainerName[len(ContainerNamePrefix):]

	notExistCheck(ContainerName)
	log.Println("Restarting cluster " + ContainerNameToShow + "...")
	if err := getDocker().ContainerRestart(ctx, ContainerName, nil); err != nil {
		log.Fatal(err)
	}
	echoInfo(ContainerName)
}
