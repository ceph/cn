package cmd

import (
	"log"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

// cliClusterStatus is the Cobra CLI call
func cliClusterStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [cluster]",
		Short: "Stat an object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   statusNano,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// statusNano shows Ceph Nano status
func statusNano(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)
	echoInfo(containerName)
}

// containerStatus checks container status
// the parameter corresponds to the type listOptions and its entry all
func containerStatus(containerName string, allList bool, containerState string) bool {
	listOptions := types.ContainerListOptions{
		All:   allList,
		Quiet: true,
	}
	containers, err := getDocker().ContainerList(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}

	// run the loop on both indexes, it's fine they have the same length
	for _, container := range containers {
		for i := range container.Names {
			if container.Names[i] == "/"+containerName && container.State == containerState {
				return true
			}
		}
	}
	return false
}
