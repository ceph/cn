package cmd

import (
	"github.com/spf13/cobra"
)

func cliEnterNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enter [cluster]",
		Short: "Connect inside a given cluster",
		Args:  cobra.ExactArgs(1),
		Run:   enterNano,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

func enterNano(cmd *cobra.Command, args []string) {
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow

	notExistCheck(containerName)
	enterContainer(containerName)
}
