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

	"github.com/spf13/cobra"
)

// cliS3CmdLs is the Cobra CLI call
func cliS3CmdLs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls [CLUSTER] [BUCKET]",
		Short: "List objects or buckets",
		Args:  cobra.RangeArgs(1, 2),
		Run:   S3CmdLs,
	}
	cmd.Flags().BoolVarP(&debugS3, "debug", "d", false, "Run S3 commands in debug mode")

	return cmd
}

// S3CmdLs wraps s3cmd command in the container
func S3CmdLs(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	notRunningCheck(containerName)

	var command []string

	if len(args) == 1 {
		command = []string{"s3cmd", "ls"}
	} else {
		command = []string{"s3cmd", "ls", "s3://" + args[1]}
	}

	if debugS3 {
		command = append(command, "--debug")
	}

	fmt.Println(execContainer(containerName, command))
}
