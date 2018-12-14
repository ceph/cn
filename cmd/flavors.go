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

	"github.com/apcera/termtables"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cmdFlavors = &cobra.Command{
		Use:   "flavors [command]",
		Short: "Interact with flavors",
		Args:  cobra.NoArgs,
	}
)

func init() {
	cmdFlavors.AddCommand(
		cliFlavorsList(),
		cliFlavorsShow(),
	)
}

// cliFlavorsList is the Cobra CLI call
func cliFlavorsList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "Print the list of flavors",
		Args:  cobra.NoArgs,
		Run:   listFlavors,
	}
	return cmd
}

// cliFlavorsShow is the Cobra CLI call
func cliFlavorsShow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show a flavor",
		Args:  cobra.ExactArgs(1),
		Run:   showFlavors,
	}
	return cmd
}

func listFlavors(cmd *cobra.Command, args []string) {
	table := termtables.CreateTable()
	table.AddHeaders("NAME", "MEMORY_SIZE", "CPU_COUNT")
	for flavor := range getItemsFromGroup(FLAVORS) {
		table.AddRow(flavor, getMemorySize(flavor), getCPUCount(flavor))
	}
	fmt.Println(table.Render())
}

func showFlavors(cmd *cobra.Command, args []string) {
	flavorName := args[0]
	flavor := FLAVORS + "." + flavorName
	if isEntryExists(FLAVORS, flavorName) {
		fmt.Println("\nDetails of the flavor " + flavorName + ":")
		if flavorName == "default" {
			PrettyPrint(getDefaultParameters())
		} else {
			//PrettyPrint(viper.AllKeys())
			PrettyPrint(viper.Get(flavor))
		}
	} else {
		// The flavor doesn't exist, let's report an empty structure
		fmt.Println("{}")
	}
}
