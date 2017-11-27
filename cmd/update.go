package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// CliUpdateNano is the Cobra CLI call
func CliUpdateNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the container image",
		Args:  cobra.NoArgs,
		Run:   updateNano,
	}
	return cmd
}

// updateNano updates the container image
func updateNano(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	if !pullImage() {
		events, err := cli.ImagePull(ctx, ImageName, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}

		d := json.NewDecoder(events)

		type Event struct {
			Status         string `json:"status"`
			Error          string `json:"error"`
			Progress       string `json:"progress"`
			ProgressDetail struct {
				Current int `json:"current"`
				Total   int `json:"total"`
			} `json:"progressDetail"`
		}

		var event *Event
		for {
			if err := d.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
		}

		if event != nil {
			if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", ImageName)) {
				fmt.Println("New image downloaded.")
			}

			if strings.Contains(event.Status, fmt.Sprintf("Image is up to date for %s", ImageName)) {
				fmt.Println("Image is up to date.")
			}
		}
	}
}
