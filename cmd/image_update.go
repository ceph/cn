package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

// CliImageUpdate is the Cobra CLI call
func CliImageUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update IMAGE",
		Short: "Update a given container image (makes sense when running on a 'latest')",
		Args:  cobra.ExactArgs(1),
		Run:   updateNano,
		Long:  "IMPORTANT: if cn was run with --image option make sure to use the same image if you're expecting to update that image",
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// updateNano updates the container image
func updateNano(cmd *cobra.Command, args []string) {
	imageName := args[0]

	if !pullImage() {
		events, err := getDocker().ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			log.Fatal(err)
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
				log.Fatal(err)
			}
		}

		if event != nil {
			if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", imageName)) {
				log.Println("New image " + imageName + " downloaded.")
			}

			if strings.Contains(event.Status, fmt.Sprintf("Image is up to date for %s", imageName)) {
				log.Println("Image " + imageName + " is up to date.")
			}
		}
	}
}
