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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

// CliImageUpdate is the Cobra CLI call
func CliImageUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update IMAGE",
		Short: "Update a given container image (makes sense when running on a 'latest')",
		Args:  cobra.ExactArgs(1),
		Run:   updateNano,
		Long:  "IMPORTANT: if cn was run with --image option make sure to use the same image if you're expecting to update that image",
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// updateNano updates the container image
func updateNano(cmd *cobra.Command, args []string) {
	imageName := args[0]

	if !pullImage() {
		events, err := getDocker().ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			log.Fatal(err)
		}

		d := json.NewDecoder(events)

		type Event struct {
			Status         string `json:"status"`
			Error          string `json:"error"`
			Progress       string `json:"progress"`
			ProgressDetail struct {
				Current int `json:"current"`
				Total   int `json:"total"`
			} `json:"progressDetail"`
		}

		var event *Event
		for {
			if err := d.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
		}

		if event != nil {
			if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", imageName)) {
				log.Println("New image " + imageName + " downloaded.")
			}

			if strings.Contains(event.Status, fmt.Sprintf("Image is up to date for %s", imageName)) {
				log.Println("Image " + imageName + " is up to date.")
			}
		}
	}
}
