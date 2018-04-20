package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const (
	cliName        = "cn"
	cliDescription = `Ceph Nano - One step S3 in container with Ceph.

                  *(((((((((((((
                (((((((((((((((((((
              ((((((((*     ,(((((((*
             ((((((             ((((((
            *((((,               ,((((/
            ((((,     ((((((/     *((((
            ((((     (((((((((     ((((
            /(((     (((((((((     ((((
             (((.     (((((((     /(((/
              (((                *((((
              .(((              (((((
         ,(((((((*             /(((
          .(((((  ((( (/  //   (((
                 /(((.  /(((((  /(((((
                        .((((/ (/
`
	cephNanoUID         = "nano"       // cephNanoUID is the uid of the S3 user
	containerNamePrefix = "ceph-nano-" // containerNamePrefix is name of the container
	tempPath            = "/tmp/"      // tempPath is the temporary path inside the container
)

var (
	// Version is the Ceph Nano version
	cnVersion = "undefined"

	// imageName is the name of the container image
	imageName = "ceph/daemon"

	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"cn"},
		//Long:
	}

	// dockerCli initializes the client connection
	dockerCli *client.Client

	// ctx opens context
	ctx = context.Background()
)

func getDocker() *client.Client {
	// If the connection with docker is not yet established
	if dockerCli == nil {
		cli, err := client.NewEnvClient()
		if err != nil {
			log.Fatal(err)
		}

		// Let's make a first Docker command to check if the protocol is consistent
		var apiVersion string
		_, err = cli.Info(ctx)
		if err != nil {
			// Oops, unable to handle server's protocol
			serverVersion := fmt.Sprint(err)
			if strings.Contains(serverVersion, "is too new") {
				ss := strings.SplitAfter(serverVersion, "Maximum supported API version is ")
				apiVersion = ss[1]
			} else if strings.Contains(serverVersion, "client is newer than server") {
				ss := strings.SplitAfter(serverVersion, "server API version: ")
				// trim last character since this 'ss[1]' is '1.24.'
				apiVersion = ss[1][:len(ss[1])-1]
			} else {
				// That's an error we don't know, let's stop here
				log.Fatal(err)
			}

			// The client version shall be degraded as it's greater than the server's one
			if len(apiVersion) > 0 {
				os.Setenv("DOCKER_API_VERSION", apiVersion)
				log.Println("Warning: degrading Docker client API version to " + apiVersion + " to match server's version.")
				// As the DOCKER_API_VERSION variable is updated, we have to restart the communication to get it
				return getDocker()
			}
		}
		// Ok, the Docker connection is valid & functional, let's return that context
		dockerCli = cli
	}
	return dockerCli
}

// Main is the main function calling the whole program
func Main(version string) {
	cnVersion = version

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		cmdCluster,
		cmdS3,
		cmdImage,
		cliVersionNano(),
	)
}

func init() {
	cobra.EnableCommandSorting = false
}
