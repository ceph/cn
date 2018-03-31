package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
)

var (
	// PrivilegedContainer whether or not the container should run Privileged
	PrivilegedContainer bool
)

// CliClusterStart is the Cobra CLI call
func CliClusterStart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   startNano,
		Example: "cn start \n" +
			"cn start --work-dir /tmp \n" +
			"cn start --image ceph/daemon:tag-stable-3.0-luminous-ubuntu-16.04",
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().StringVarP(&WorkingDirectory, "work-dir", "d", "/usr/share/ceph-nano", "Directory to work from")
	cmd.Flags().StringVarP(&ImageName, "image", "i", "ceph/daemon", "USE AT YOUR OWN RISK. Ceph container image to use, format is 'username/image:tag'.")
	cmd.Flags().BoolVar(&PrivilegedContainer, "privileged", false, "Starts the container in privileged mode")
	cmd.Flags().BoolVar(&Help, "help", false, "help for start")

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
	ContainerName := ContainerNamePrefix + args[0]
	ContainerNameToShow := ContainerName[len(ContainerNamePrefix):]

	if status := containerStatus(ContainerName, true, "created"); status {
		removeContainer(ContainerName)
	}

	if status := containerStatus(ContainerName, false, "running"); status {
		fmt.Println("Cluster " + ContainerNameToShow + " is already running!")
	} else if status := containerStatus(ContainerName, true, "exited"); status {
		fmt.Println("Starting cluster " + ContainerNameToShow + "...")
		startContainer(ContainerName)
	} else {
		pullImage()
		fmt.Println("Running cluster " + ContainerNameToShow + "...")
		runContainer(cmd, args)
	}
	echoInfo(ContainerName)
}

// runContainer creates a new container when nothing exists
func runContainer(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]
	RgwPort := generateRGWPortToUse()
	if RgwPort == "notfound" {
		log.Fatal("Unable to find a port between 8000 and 8100.")
	}
	RgwNatPort := RgwPort + "/tcp"

	exposedPorts := nat.PortSet{
		nat.Port(RgwNatPort): {},
	}

	portBindings := nat.PortMap{
		nat.Port(RgwNatPort): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: RgwPort,
			},
		},
	}

	envs := []string{
		"RGW_CIVETWEB_PORT=" + RgwPort, // DON'T TOUCH MY POSITION IN THE SLICE OR YOU WILL BREAK dockerInspect()
		"DEBUG=verbose",
		"CEPH_DEMO_UID=" + CephNanoUID,
		"NETWORK_AUTO_DETECT=4",
		"CEPH_DAEMON=demo",
		"DEMO_DAEMONS=mon,mgr,osd,rgw"}

	ressources := container.Resources{
		Memory:   536870912, // 512MB
		NanoCPUs: 1,
	}

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
		Binds:        []string{WorkingDirectory + ":" + TempPath},
		Resources:    ressources,
		Privileged:   PrivilegedContainer,
	}

	resp, err := getDocker().ContainerCreate(ctx, config, hostConfig, nil, ContainerName)
	if err != nil {
		log.Fatal(err)
	}

	err = getDocker().ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
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
			log.Fatal(err)
		}
	}
}

// startContainer starts a container that is stopped
func startContainer(ContainerName string) {
	if err := getDocker().ContainerStart(ctx, ContainerName, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}
}
