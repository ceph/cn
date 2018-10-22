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
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// S3CmdSkip means do not do anything when object exists
	S3CmdSkip bool

	// S3CmdContinue means
	S3CmdContinue bool

	// S3CmdOpt is the option to apply when trying to get a file and the destination already exists
	S3CmdOpt string
)

// cliS3CmdGet is the Cobra CLI call
func cliS3CmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [CLUSTER] [BUCKET/OBJECT] [LOCAL_FILE]",
		Short: "Get file out of a bucket",
		Args:  cobra.RangeArgs(2, 3),
		Run:   S3CmdGet,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVarP(&S3CmdSkip, "skip", "s", true, "Skip over files that exist at the destination")
	cmd.Flags().BoolVarP(&S3CmdForce, "force", "f", false, "Force overwrite files that exist at the destination")
	cmd.Flags().BoolVarP(&S3CmdContinue, "continue", "c", false, "Continue getting a partially downloaded file")
	cmd.Flags().BoolVarP(&debugS3, "debug", "d", false, "Run S3 commands in debug mode")

	return cmd
}

// S3CmdGet wraps s3cmd command in the container
func S3CmdGet(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	notRunningCheck(containerName)
	BucketObjectName := args[1]
	var fileName string

	if len(args) > 1 {
		fileName = args[2]
	} else {
		fileName = BucketObjectName
	}

	BucketObjectNameBase := path.Base(BucketObjectName)

	if S3CmdForce {
		S3CmdOpt = "--force"
	} else if S3CmdSkip {
		S3CmdOpt = "--skip-existing"
	} else if S3CmdContinue {
		S3CmdOpt = "--continue"
	}
	// if args
	command := []string{"s3cmd", "get", S3CmdOpt, "s3://" + BucketObjectName, tempPath}
	if debugS3 {
		command = append(command, "--debug")
	}

	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerNameToShow

	dir := dockerInspect(containerName, "Binds")
	if fileName != BucketObjectName {
		//if _, err := os.Stat(fileName); os.Stat.Mode.IsDir(err) {
		if info, err := os.Stat(fileName); err == nil && info.IsDir() {
			fileName = fileName + "/" + BucketObjectNameBase
		}
		_, err := copyFile(dir+"/"+BucketObjectNameBase, fileName)
		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println(output)
}
