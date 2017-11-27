package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// CliStartNano is the Cobra CLI call
func CliStartNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start object storage server",
		Args:  cobra.NoArgs,
		Run:   startNano,
		Example: "cn start --work-dir /tmp \n" +
			"cn start",
	}
	cmd.Flags().StringVarP(&WorkingDirectory, "work-dir", "d", "/usr/share/ceph-nano", "Directory to work from")
	return cmd
}

// startNano starts Ceph Nano
func startNano(cmd *cobra.Command, args []string) {
	// Test for a leftover container
	// Usually happens when someone fails to run the container on an exposed directory
	// Typical error on Docker For Mac you will see:
	// panic: Error response from daemon: Mounts denied:
	// The path /usr/share/ceph-nano is not shared from OS X and is not known to Docker.
	// You can configure shared paths from Docker -> Preferences... -> File Sharing.
	if status := containerStatus(true, "created"); status {
		removeContainer(ContainerName)
	}

	if status := containerStatus(false, "running"); status {
		fmt.Println("ceph-nano is already running!")
	} else if status := containerStatus(true, "exited"); status {
		fmt.Println("Starting ceph-nano...")
		startContainer()
	} else {
		fmt.Println("Running ceph-nano...")
		runContainer(cmd, args)
	}
	echoInfo()
}

// runContainer creates a new container when nothing exists
func runContainer(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	pullImage()

	exposedPorts, portBindings, _ := nat.ParsePortSpecs([]string{":8000:8000"})
	envs := []string{
		"DEBUG=verbose",
		"CEPH_DEMO_UID=" + CephNanoUID,
		"NETWORK_AUTO_DETECT=4",
		"MON_IP=127.0.0.1",
		"RGW_CIVETWEB_PORT=" + RgwPort,
		"CEPH_DAEMON=demo"}

	config := &container.Config{
		Image:        ImageName,
		Hostname:     ContainerName + "-faa32aebf00b",
		ExposedPorts: exposedPorts,
		Env:          envs,
		Volumes: map[string]struct{}{
			"/etc/ceph":     struct{}{},
			"/var/lib/ceph": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        []string{WorkingDirectory + ":/tmp"},
	}

	// TODO --memory 512m

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, ContainerName)
	if err != nil {
		panic(err)
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	// The if removes the error:
	//panic: runtime error: invalid memory address or nil pointer dereference
	//[signal SIGSEGV: segmentation violation code=0x1 addr=0x20 pc=0x137a2b4]
	if err != nil {
		if strings.Contains(err.Error(), "Mounts denied") {
			fmt.Println("ERROR: It looks like you need to use the --work-dir option. \n" +
				"This typically happens when Docker is not running natively (e.g: Docker for Mac/Windows). \n" +
				"The path /usr/share/ceph-nano is not shared from OS X / Windows and is not known to Docker. \n" +
				"You can configure shared paths from Docker -> Preferences... -> File Sharing.) \n" +
				"Alternatively, you can simply use the --work-dir option to point to an already shared directory. \n" +
				"On Docker for Mac / Windows, shared directories can be found in the settings.")
			cmd.Help()
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}

// startContainer starts a container that is stopped
func startContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, ContainerName, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}
