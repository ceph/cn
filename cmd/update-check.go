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
	"strings"

	"github.com/elgs/gojq"
	"github.com/spf13/cobra"
)

// cliUpdateCheckNano is the Cobra CLI call
func cliUpdateCheckNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-check",
		Short: "Print cn current and latest version number",
		Args:  cobra.NoArgs,
		Run:   updateCheckNano,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// updateCheckNano print Ceph Nano version
func updateCheckNano(cmd *cobra.Command, args []string) {
	url := githubCNReleasesURL
	output := curlURL(url)

	parser, err := gojq.NewStringQuery(string(output))
	if err != nil {
		log.Fatal(err)
	}

	message, err := parser.Query("message")
	// if a message exists in the answer, let's print it and return
	if err == nil {
		fmt.Println(message)
		return
	}

	latestTag, err := parser.Query("[0].tag_name")
	if err != nil {
		log.Fatal(err)
	}

	cnVersionSplit := strings.Fields(cnVersion)
	cnVersionNum := cnVersionSplit[0]

	fmt.Println("Current version:", cnVersionNum)
	fmt.Println("Latest version:", latestTag)

	if latestTag != cnVersionNum {
		latestTagURL, err := parser.Query("[0].html_url")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("There is a newer version of cn available. Download it here:", latestTagURL)
	}
}
