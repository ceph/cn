package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cmdImage = &cobra.Command{
		Use:   "image [command] [arg]",
		Short: "Interact with cn's container image(s)",
		Args:  cobra.NoArgs,
	}
)

func init() {
	cmdImage.AddCommand(
		CliImageUpdate(),
		CliImageList(),
	)
}
