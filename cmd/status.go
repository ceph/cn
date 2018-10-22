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
