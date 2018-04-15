package cmd

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
)

var (
	// PrivilegedContainer whether or not the container should run Privileged
	PrivilegedContainer bool

	// Data points to either the directory or drive to use to store Ceph's data
	Data string
)

// CliClusterStart is the Cobra CLI call
func CliClusterStart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   startNano,
		Example: "cn cluster start mycluster \n" +
			"cn cluster start mycluster --work-dir /tmp \n" +
			"cn cluster start mycluster --image ceph/daemon:latest-luminous \n" +
			"cn cluster start mycluster -b /dev/sdb \n" +
			"cn cluster start mycluster -b /srv/nano \n" +
			"cn cluster start mycluster --privileged \n",
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().StringVarP(&WorkingDirectory, "work-dir", "d", "/usr/share/ceph-nano", "Directory to work from")
	cmd.Flags().StringVarP(&ImageName, "image", "i", "ceph/daemon", "USE AT YOUR OWN RISK. Ceph container image to use, format is 'username/image:tag'.")
	cmd.Flags().StringVarP(&Data, "data", "b", "", "Configure Ceph Nano underlying storage with a specific directory or physical block device. Block device support only works on Linux running under 'root', only also directory might need running as 'root' if SeLinux is enabled.")
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
		log.Println("Cluster " + ContainerNameToShow + " is already running!")
	} else if status := containerStatus(ContainerName, true, "exited"); status {
		log.Println("Starting cluster " + ContainerNameToShow + "...")
		startContainer(ContainerName)
	} else {
		pullImage()
		runContainer(cmd, args)
	}
	echoInfo(ContainerName)
}

// runContainer creates a new container when nothing exists
func runContainer(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]
	ContainerNameToShow := ContainerName[len(ContainerNamePrefix):]
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
		"DEMO_DAEMONS=mon,mgr,osd,rgw",
	}

	volumeBindings := []string{
		WorkingDirectory + ":" + TempPath,
	}

	ressources := container.Resources{
		Memory:   536870912, // 512MB
		NanoCPUs: 1,
	}

	volumes := map[string]struct{}{
		"/etc/ceph":     struct{}{},
		"/var/lib/ceph": struct{}{},
	}

	if len(Data) != 0 {
		testDev, err := GetFileType(Data)
		if err != nil {
			log.Fatal(err)
		}
		if testDev != "blockdev" && testDev != "directory" {
			log.Fatalf("We only accept a directory or a block device, however the specified file type is a %s", testDev)
		}
		if testDev == "directory" {
			testEmptyDir := IsEmpty(Data)
			if !testEmptyDir {
				log.Fatal(Data + " is not empty, doing nothing.")
			}
			if runtime.GOOS == "linux" {
				ApplySeLinuxLabel(Data)
			}
			envs = append(envs, "OSD_PATH="+Data)
			volumeBindings = append(volumeBindings, Data+":"+Data)
		}
		if testDev == "blockdev" {
			meUserName, meID := whoAmI()
			if meID != "0" {
				log.Fatal("Hey " + meUserName + "! Run me as 'root' when using a block device.")
			}
			// We don't have the logic to do the introspection without using blkid.
			// Unfortunately, blkid is not available on macOS or Windows
			if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
				log.Fatal("Operating system: " + runtime.GOOS + " is not supported in the scenario")
			}
			// We run a couple of test here to ensure the device can be used:
			// 1. test of the device is accessed by a process (open it with O_EXCL)
			// 2. test if the device has a partition table and/or a filesystem
			// 3. test if they are partitions in that partition table (you can have a partition table with 0 partitions)

			// First test: is the device opened by a process?!«
			testDevOpen, _ := ExclusiveOpenFailsOnDevice(Data)
			if testDevOpen {
				log.Fatal(Data + " is accessed by another process, doing nothing.")
			}

			// Second test: search for filesystem and partition table
			diskFormat := GetDiskFormat(Data)
			lines := strings.Split(diskFormat, "\n")
			var fstype, pttype string
			for _, l := range lines {
				if len(l) <= 0 {
					// Ignore empty line.
					continue
				}
				cs := strings.Split(l, "=")
				if len(cs) != 2 {
					log.Fatal("blkid returns invalid output: " + diskFormat + ". This potentially means no partition label on your disk.")
				}
				// TYPE is filesystem type, and PTTYPE is partition table type, according
				// to https://www.kernel.org/pub/linux/utils/util-linux/v2.21/libblkid-docs/.
				if cs[0] == "TYPE" {
					fstype = cs[1]
					log.Fatal(Data + " has a filesystem: " + fstype + ", doing nothing.")
				} else if cs[0] == "PTTYPE" {
					// Third test: number of partitions
					pttype = cs[1]
					// Now we test if the disk has partition(s)
					// We know parted will return 2 lines if there is a partition table:
					//
					// BYT;
					// /dev/sdc:100GB:scsi:512:512:gpt:HP LOGICAL VOLUME:;
					//
					// So we remove the first 2 lines of the output
					// The third one is always the partition number
					partedUselessLinesCount := 2
					num := GetDiskPartitions(Data)
					partCount := len(num) - partedUselessLinesCount
					if partCount != 0 {
						log.Fatal(Data + " has a partition table type " + pttype + " and " + strconv.Itoa(partCount) + " partition(s) doing nothing.")
					}
				}
			}
			// If we arrive here, it should be safe to use the device.
			envs = append(envs, "OSD_DEVICE="+Data)
			PrivilegedContainer = true
			volumeBindings = append(volumeBindings, "/dev:/dev")
			// place holder once 'demo' will use ceph-volume
			// volumeBindings = append(volumeBindings, "/run/lvm/lvmetad.socket:/run/lvm/lvmetad.socket")
		}
	}

	config := &container.Config{
		Image:        ImageName,
		Hostname:     ContainerName + "-faa32aebf00b",
		ExposedPorts: exposedPorts,
		Env:          envs,
		Volumes:      volumes,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        volumeBindings,
		Resources:    ressources,
		Privileged:   PrivilegedContainer,
	}

	log.Println("Running cluster " + ContainerNameToShow + "...")

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
			log.Println("ERROR: It looks like you need to use the --work-dir option. \n" +
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
