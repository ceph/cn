package cmd

import (
	"encoding/json"
	"strconv"

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
		Short: "List ceph/daemon tags (default print the first 10 tags)",
		Args:  cobra.NoArgs,
		Run:   listImageTags,
	}
	cmd.Flags().BoolVarP(&ListAllTags, "all", "a", false, "List all the tags for ceph/daemon (can be verbose)")

	return cmd
}

// listImageTags lists container image tags
func listImageTags(cmd *cobra.Command, args []string) {
	var numPage int
	var url string

	// Creating the maps for JSON
	m := map[string]interface{}{}

	numPage = 1
	if ListAllTags {
		numPage = pageCount()
	}

	for i := 1; i <= numPage; i++ {
		// convert numPage into a string for concatenation, (see the end)
		url = "https://registry.hub.docker.com/v2/repositories/ceph/daemon/tags/?page=" + strconv.Itoa(i)
		output := curlURL(url)

		// Parsing/Unmarshalling JSON encoding/json
		err := json.Unmarshal([]byte(output), &m)
		if err != nil {
			panic(err)
		}
		parseMap(m, "name")
	}
}
