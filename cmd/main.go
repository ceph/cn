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
	cliDescription = "Ceph Nano - One step S3 in container with Ceph."
)

var (
	// Version is the Ceph Nano version
	cnVersion = "undefined"

	// WorkingDirectory is the working directory where objects can be put inside S3
	WorkingDirectory = "/usr/share/ceph-nano"

	// CephNanoUID is the uid of the S3 user
	CephNanoUID = "nano"

	// RgwPort is the rgw listenning port
	RgwPort = "8000"

	// ContainerName is name of the container
	ContainerName = "ceph-nano"

	// ImageName is the name of the container image
	ImageName = "ceph/daemon"

	// TempPath is the temporary path inside the container
	TempPath = "/tmp/"

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
			// Oups, unable to handle server's protocol
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
				fmt.Println("Warning: Degrading Docker client API version to " + apiVersion + " to match server's version")
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
	validateEnv()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		CliStartNano(),
		CliStopNano(),
		CliStatusNano(),
		CliPurgeNano(),
		CliLogsNano(),
		CliRestartNano(),
		cmdS3,
		CliUpdateNano(),
		CliVersionNano(),
	)
}

func init() {
	cobra.EnableCommandSorting = false
}
