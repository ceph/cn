package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

var (
	// IamSure means whether or not the user wants to purge
	IamSure bool

	// Help shows a customer help
	Help bool

	// DeleteAll also deletes the container image
	DeleteAll bool
)

// CliClusterPurge is the Cobra CLI call
func CliClusterPurge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Purge object storage server. DANGEROUS!",
		Args:  cobra.ExactArgs(1),
		Run:   purgeNano,
		DisableFlagsInUseLine: false,
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVar(&IamSure, "yes-i-am-sure", false, "YES I know what I'm doing and I want to purge")
	cmd.Flags().BoolVar(&DeleteAll, "all", false, "This also deletes the container image")
	cmd.Flags().BoolVar(&Help, "help", false, "help for purge")

	return cmd
}

// purgeNano purges Ceph Nano.
func purgeNano(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]
	ContainerNameToShow := ContainerName[len(ContainerNamePrefix):]

	if !IamSure {
		fmt.Printf("Purge option is too dangerous please set the right flag. \n \n")
		cmd.Help()
		os.Exit(1)
	}
	notExistCheck(ContainerName)
	log.Println("Purging cluster " + ContainerNameToShow + "...")
	removeContainer(ContainerName)
}

func removeContainer(ContainerName string) {
	Data := dockerInspect(ContainerName, "BindsData")

	if DeleteAll {
		ImageName = dockerInspect(ContainerName, "image")
	}
	options := types.ContainerRemoveOptions{
		RemoveLinks:   false,
		RemoveVolumes: true,
		Force:         true,
	}
	// we don't necessarily want to catch errors here
	// it's not an issue if the container does not exist
	getDocker().ContainerRemove(ctx, ContainerName, options)

	if Data != "noDataDir" && Data != "/dev" {
		testDev, err := GetFileType(Data)
		if err != nil {
			log.Fatal(err)
		}
		if testDev == "directory" {
			err := os.RemoveAll(Data)
			if err != nil {
				log.Println("Something went wrong while removing " + Data + ".\n" +
					"You need to purge the directory manually, next time run me as 'root' to avoid that.")
				log.Fatal(err)
			}
		}
	}

	if DeleteAll {
		options := types.ImageRemoveOptions{
			Force:         true,
			PruneChildren: true,
		}
		log.Println("Removing container image" + ImageName + "...")
		getDocker().ImageRemove(ctx, ImageName, options)
	}
}
