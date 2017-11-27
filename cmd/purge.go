package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// IamSure means whenver or not the user wants to purge
	IamSure bool

	// Help shows a customer help
	Help bool
)

// CliPurgeNano is the Cobra CLI call
func CliPurgeNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Purge object storage server. DANGEROUS!",
		Args:  cobra.NoArgs,
		Run:   purgeNano,
		DisableFlagsInUseLine: false,
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVar(&IamSure, "yes-i-am-sure", false, "YES I know what I'm doing and I want to purge")
	cmd.Flags().BoolVar(&Help, "help", false, "help for purge")
	return cmd
}

// purgeNano purges Ceph Nano.
func purgeNano(cmd *cobra.Command, args []string) {
	if !IamSure {
		fmt.Printf("Purge option is too dangerous please set the right flag. \n \n")
		cmd.Help()
		os.Exit(1)
	}
	notExistCheck()
	fmt.Println("Purging ceph-nano... ")
	removeContainer(ContainerName)
}

func removeContainer(name string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := types.ContainerRemoveOptions{
		RemoveLinks:   false,
		RemoveVolumes: true,
		Force:         true,
	}
	// we don't necessarily want to catch errors here
	// it's not an issue if the container does not exist
	cli.ContainerRemove(ctx, name, options)
}
