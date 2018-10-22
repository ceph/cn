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

// cliClusterLogs is the Cobra CLI call
func cliClusterLogs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [cluster]",
		Short: "Print an object storage server logs",
		Args:  cobra.ExactArgs(1),
		Run:   logsNano,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// logsNano prints rgw logs
func logsNano(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]
	showS3Logs(containerName)
}

func showS3Logs(containerName string) {
	notExistCheck(containerName)
	c := []string{"cat", "/var/log/ceph/client.rgw." + containerName + "-faa32aebf00b.log"}
	output := execContainer(containerName, c)
	if strings.Contains("No such file or directory", output) {
		return
	}
	fmt.Printf("%s", output)
}
