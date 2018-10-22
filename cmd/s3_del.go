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
	"strings"

	"github.com/spf13/cobra"
)

var (
	// S3CmdRec is the option to apply when trying to delete content
	S3CmdRec bool
)

// cliS3CmdDel is the Cobra CLI call
func cliS3CmdDel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del [CLUSTER] [BUCKET/OBJECT]",
		Short: "Delete file from bucket",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdDel,
	}
	cmd.Flags().BoolVarP(&debugS3, "debug", "d", false, "Run S3 commands in debug mode")
	//cmd.Flags().BoolVarP(&S3CmdRec, "recursive", "r", false, "Recursive removal.")
	//cmd.Flags().BoolVarP(&S3CmdForce, "force", "f", false, "Force removal.")

	return cmd
}

// S3CmdDel wraps s3cmd command in the container
func S3CmdDel(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	notRunningCheck(containerName)

	/*
		S3CmdOpt = "--verbose"
			if S3CmdRec {
				S3CmdOpt = S3CmdOpt + " --recursive"
			}
			if S3CmdForce {
				S3CmdOpt = S3CmdOpt + " --force"
			}
			command := []string{"s3cmd", "del", S3CmdOpt, "s3://" + args[0]}
	*/
	command := []string{"s3cmd", "del", "s3://" + args[1]}
	if debugS3 {
		command = append(command, "--debug")
	}

	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerNameToShow
	fmt.Println(output)
}
