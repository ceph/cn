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
	"github.com/spf13/cobra"
)

var (
	cmdS3 = &cobra.Command{
		Use:   "s3 [command] [arg]",
		Short: "Interact with a particular S3 object server",
		Args:  cobra.NoArgs,
	}
	// S3CmdForce means force operation
	S3CmdForce bool

	// debugS3 means use the '--debug' flag in the s3cmd command
	debugS3 bool
)

func init() {
	cmdS3.AddCommand(
		cliS3CmdMb(),
		cliS3CmdRb(),
		cliS3CmdLs(),
		cliS3CmdLa(),
		cliS3CmdPut(),
		cliS3CmdGet(),
		cliS3CmdDel(),
		cliS3CmdDu(),
		cliS3CmdInfo(),
		cliS3CmdCp(),
		cliS3CmdMv(),
		cliS3CmdSync())
}
