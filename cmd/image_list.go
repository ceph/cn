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
	"os"

	"github.com/spf13/cobra"
)

var (
	// ListAllTags whether or not to list all the image tags
	ListAllTags bool
)

// CliImageList is the Cobra CLI call
func CliImageList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List container image tags (default print the first 100 tags)",
		Args:  cobra.NoArgs,
		Run:   listImageTags,
	}
	cmd.Flags().BoolVarP(&ListAllTags, "all", "a", false, "List all the tags of the container image (can be verbose)")

	return cmd
}

// listImageTags lists container image tags
func listImageTags(cmd *cobra.Command, args []string) {
	if os.Getenv("CN_REGISTRY") == "redhat" {
		listRedHatRegistryImageTags()
	} else {
		listDockerRegistryImageTags()
	}
	listQuayRegistryImageTags()
}
