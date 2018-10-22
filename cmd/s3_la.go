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

// cliS3CmdLa is the Cobra CLI call
func cliS3CmdLa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "la [CLUSTER]",
		Short: "List all object in all buckets",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdLa,
	}
	cmd.Flags().BoolVarP(&debugS3, "debug", "d", false, "Run S3 commands in debug mode")

	return cmd
}

// S3CmdLa wraps s3cmd command in the container
func S3CmdLa(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	notRunningCheck(containerName)

	command := []string{"s3cmd", "la"}
	if debugS3 {
		command = append(command, "--debug")
	}

	output := string(execContainer(containerName, command))

	if len(output) == 1 {
		command = []string{"s3cmd", "ls"}
		output = execContainer(containerName, command)
		if debugS3 {
			command = append(command, "--debug")
		}
	}
	fmt.Println(output)
}
