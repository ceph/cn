package cmd

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// CliStatusNano is the Cobra CLI call
func CliStatusNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Stat object storage server",
		Args:  cobra.NoArgs,
		Run:   statusNano,
	}
	return cmd
}

// statusNano shows Ceph Nano status
func statusNano(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	echoInfo()
}

// containerStatus checks container status
// the parameter corresponds to the type listOptions and its entry all
func containerStatus(allList bool, containerState string) bool {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	listOptions := types.ContainerListOptions{
		All:   allList,
		Quiet: true,
	}
	containers, err := cli.ContainerList(context.Background(), listOptions)
	if err != nil {
		panic(err)
	}

	// run the loop on both indexes, it's fine they have the same length
	for _, container := range containers {
		for i := range container.Names {
			if container.Names[i] == "/ceph-nano" && container.State == containerState {
				return true
			}
		}
	}
	return false
}
