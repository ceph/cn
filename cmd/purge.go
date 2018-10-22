/*
 * Ceph Nano (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Below main package has canonical imports for 'go get' and 'go build'
 * to work with all other clones of github.com/ceph/cn repository. For
 * more information refer https://golang.org/doc/go1.4#canonicalimports
 */

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

// cliClusterPurge is the Cobra CLI call
func cliClusterPurge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge [cluster]",
		Short: "Purge an object storage server. DANGEROUS!",
		Args:  cobra.ExactArgs(1),
		Run:   purgeNano,
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVar(&IamSure, "yes-i-am-sure", false, "YES I know what I'm doing and I want to purge")
	cmd.Flags().BoolVar(&DeleteAll, "all", false, "This also deletes the container image")
	cmd.Flags().BoolVar(&Help, "help", false, "help for purge")

	return cmd
}

// purgeNano purges Ceph Nano.
func purgeNano(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	if !IamSure {
		fmt.Printf("Purge option is too dangerous please set the right flag. \n \n")
		cmd.Help()
		os.Exit(1)
	}
	notExistCheck(containerName)
	log.Println("Purging cluster " + containerNameToShow + "...")
	removeContainer(containerName)
}

func removeContainer(containerName string) {
	dataOsd := dockerInspect(containerName, "BindsData")

	if DeleteAll {
		imageName = dockerInspect(containerName, "image")
	}
	options := types.ContainerRemoveOptions{
		RemoveLinks:   false,
		RemoveVolumes: true,
		Force:         true,
	}
	// we don't necessarily want to catch errors here
	// it's not an issue if the container does not exist
	getDocker().ContainerRemove(ctx, containerName, options)

	if dataOsd != "noDataDir" && dataOsd != "/dev" {
		testDev, err := getFileType(dataOsd)
		if err != nil {
			log.Fatal(err)
		}
		if testDev == "directory" {
			err := os.RemoveAll(dataOsd)
			if err != nil {
				log.Println("Something went wrong while removing " + dataOsd + ".\n" +
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
		log.Println("Removing container image " + imageName + "...")
		getDocker().ImageRemove(ctx, imageName, options)
	}
}
