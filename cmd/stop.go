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
	"os"
	"time"

	"github.com/spf13/cobra"
)

// cliClusterStop is the Cobra CLI call
func cliClusterStop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [cluster]",
		Short: "Stop an object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   stopNano,
		DisableFlagsInUseLine: true,
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
