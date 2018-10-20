package cmd

import (
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
	listDockerRegistryImageTags()
}
