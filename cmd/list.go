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
	"regexp"

	"github.com/apcera/termtables"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

// cliClusterList is the Cobra CLI call
func cliClusterList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "Print the list of object storage server(s)",
		Args:  cobra.NoArgs,
		Run:   listNano,
	}
	return cmd
}

// listNano prints running Ceph cluster(s)
func listNano(cmd *cobra.Command, args []string) {
	showNanoClusters()
}

func showNanoClusters() {
	listOptions := types.ContainerListOptions{
		All:   true,
		Quiet: true,
	}
	containers, err := getDocker().ContainerList(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}

	table := termtables.CreateTable()
	table.AddHeaders("NAME", "STATUS", "IMAGE", "IMAGE RELEASE", "IMAGE CREATION TIME")

	// run the loop on both indexes, it's fine they have the same length
	for _, container := range containers {
		for i := range container.Names {
			match, _ := regexp.MatchString(containerNamePrefix, container.Names[i])
			if match {
				// remove 7 first char since container.ImageID is in the form of sha256:<ID>
				containerImgTag := inspectImage(container.ImageID[7:], "tag")
				containerImgCreated := inspectImage(container.ImageID[7:], "created")
				containerImgRelease := inspectImage(container.ImageID[7:], "release")
				containerNameToShow := container.Names[i][len(containerNamePrefix):]
				// We trim again so we can remove the '/' since container name returned is /ceph-nano
				table.AddRow(containerNameToShow[1:], container.State, containerImgTag, containerImgRelease, containerImgCreated)
			}
		}
	}
	fmt.Println(table.Render())
}
