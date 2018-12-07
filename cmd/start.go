/*
 * Ceph Nano (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Below main package has canonical imports for 'go get' and 'go build'
 * to work with all other clones of github.com/ceph/cn repository. For
 * more information refer https://golang.org/doc/go1.4#canonicalimports
 */

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
	// privilegedContainer whether or not the container should run Privileged
	privilegedContainer bool

	// dataOsd points to either the directory or drive to use to store Ceph's data
	dataOsd string

	// workingDirectory is the working directory where objects can be put inside S3
	workingDirectory string

	// sizeBluestoreBlock is the size of BLUESTORE_BLOCK_SIZE
	sizeBluestoreBlock string
)

// cliClusterStart is the Cobra CLI call
func cliClusterStart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [cluster]",
		Short: "Start an object storage server",
		Args:  cobra.ExactArgs(1),
		Run:   startNano,
		Example: "cn cluster start mycluster \n" +
			"cn cluster start mycluster --work-dir /tmp \n" +
			"cn cluster start mycluster --image ceph/daemon:latest-luminous \n" +
			"cn cluster start mycluster -b /dev/sdb \n" +
			"cn cluster start mycluster -b /srv/nano -s 20GB \n" +
			"cn cluster start mycluster --privileged \n",
	}
	cmd.Flags().SortFlags = false
	cmd.Flags().StringVarP(&workingDirectory, "work-dir", "d", "/usr/share/ceph-nano", "Directory to work from")
	cmd.Flags().StringVarP(&imageName, "image", "i", "ceph/daemon", "USE AT YOUR OWN RISK. Ceph container image to use, format is 'registry/username/image:tag'. The registry is optional")
	cmd.Flags().StringVarP(&dataOsd, "data", "b", "", "Configure Ceph Nano underlying storage with a specific directory or physical block device. Block device support only works on Linux running under 'root', only also directory might need running as 'root' if SeLinux is enabled.")
	cmd.Flags().StringVarP(&sizeBluestoreBlock, "size", "s", "", "Configure Ceph Nano underlying storage size when using a specific directory")
	cmd.Flags().BoolVar(&privilegedContainer, "privileged", false, "Starts the container in privileged mode")
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
	containerNameToShow := args[0]
	containerName := containerNamePrefix + containerNameToShow
	if len(workingDirectory) == 0 {
		workingDirectory = workingDirectory + "-" + containerNameToShow
	}

	if status := containerStatus(containerName, true, "created"); status {
		removeContainer(containerName)
	}

	if status := containerStatus(containerName, false, "running"); status {
		log.Println("Cluster " + containerNameToShow + " is already running!")
	} else if status := containerStatus(containerName, true, "exited"); status {
		log.Println("Starting cluster " + containerNameToShow + "...")
		startContainer(containerName)
	} else {
		pullImage()
		runContainer(cmd, args)
	}
	echoInfo(containerName)
}

// runContainer creates a new container when nothing exists
func runContainer(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]
	containerNameToShow := args[0]
	rgwPort := generateRGWPortToUse()
	if rgwPort == "notfound" {
		log.Fatal("Unable to find a port between 8000 and 8100 for the S3 endpoint.")
	}
	cnBrowserPort := generateBrowserPortToUse()
	if cnBrowserPort == "notfound" {
		log.Fatal("Unable to find a port between 5000 and 5100 for the UI endpoint.")
	}
	rgwNatPort := rgwPort + "/tcp"
	cnBrowserNatPort := cnBrowserPort + "/tcp"

	exposedPorts := nat.PortSet{
		nat.Port(rgwNatPort):       {},
		nat.Port(cnBrowserNatPort): {},
	}

	portBindings := nat.PortMap{
		nat.Port(rgwNatPort): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: rgwPort,
			},
		},
		nat.Port(cnBrowserNatPort): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: cnBrowserPort,
			},
		},
	}

	ips, _ := getInterfaceIPv4s()

	envs := []string{
		"RGW_CIVETWEB_PORT=" + rgwPort, // DON'T TOUCH MY POSITION IN THE SLICE OR YOU WILL BREAK dockerInspect()
		"SREE_PORT=" + cnBrowserPort,   // DON'T TOUCH MY POSITION IN THE SLICE OR YOU WILL BREAK dockerInspect()
		"EXPOSED_IP=" + ips[0].String(),
		"DEBUG=verbose",
		"CEPH_DEMO_UID=" + cephNanoUID,
		"MON_IP=127.0.0.1",
		"CEPH_PUBLIC_NETWORK=0.0.0.0/0",
		"CEPH_DAEMON=demo",
		"DEMO_DAEMONS=mon,mgr,osd,rgw",
		"SREE_VERSION=v0.1", // keep this for backward compatiblity, the option is gone since https://github.com/ceph/ceph-container/pull/1232
	}

	volumeBindings := []string{
		workingDirectory + ":" + tempPath,
	}

	ressources := container.Resources{
		Memory:   0, // To be adjusted by config
		NanoCPUs: 1,
	}

	volumes := map[string]struct{}{
		"/etc/ceph":     struct{}{},
		"/var/lib/ceph": struct{}{},
	}

	if len(dataOsd) != 0 {
		testDev, err := getFileType(dataOsd)
		if err != nil {
			log.Fatal(err)
		}
		if testDev != "blockdev" && testDev != "directory" {
			log.Fatalf("We only accept a directory or a block device, however the specified file type is a %s", testDev)
		}
		if testDev == "directory" {
			testEmptyDir := isEmpty(dataOsd)
			if !testEmptyDir {
				log.Fatal(dataOsd + " is not empty, doing nothing.")
			}
			if runtime.GOOS == "linux" {
				applySeLinuxLabel(dataOsd)
			}
			envs = append(envs, "OSD_PATH="+dataOsd)
			volumeBindings = append(volumeBindings, dataOsd+":"+dataOsd)

			// Did someone specify a particular size for cn data store in this directory?
			if len(sizeBluestoreBlock) != 0 {
				sizeBluestoreBlockToBytes := toBytes(sizeBluestoreBlock)
				if sizeBluestoreBlockToBytes == 0 {
					log.Fatal("Wrong unit passed: ", sizeBluestoreBlock, ". Please refer to https://en.wikipedia.org/wiki/Byte.")
				}
				envs = append(envs, "BLUESTORE_BLOCK_SIZE="+string(sizeBluestoreBlockToBytes))
			}
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
			testDevOpen, _ := exclusiveOpenFailsOnDevice(dataOsd)
			if testDevOpen {
				log.Fatal(dataOsd + " is accessed by another process, doing nothing.")
			}

			// Second test: search for filesystem and partition table
			diskFormat := getDiskFormat(dataOsd)
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
					log.Fatal(dataOsd + " has a filesystem: " + fstype + ", doing nothing.")
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
					num := getDiskPartitions(dataOsd)
					partCount := len(num) - partedUselessLinesCount
					if partCount != 0 {
						log.Fatal(dataOsd + " has a partition table type " + pttype + " and " + strconv.Itoa(partCount) + " partition(s) doing nothing.")
					}
				}
			}
			// If we arrive here, it should be safe to use the device.
			envs = append(envs, "OSD_DEVICE="+dataOsd)
			privilegedContainer = true
			volumeBindings = append(volumeBindings, "/dev:/dev")
			// place holder once 'demo' will use ceph-volume
			// volumeBindings = append(volumeBindings, "/run/lvm/lvmetad.socket:/run/lvm/lvmetad.socket")
		}
	}

	config := &container.Config{
		Image:        imageName,
		Hostname:     containerName + "-faa32aebf00b",
		ExposedPorts: exposedPorts,
		Env:          envs,
		Volumes:      volumes,
	}

	ressources.Memory = getMemorySizeInBytes(containerNameToShow)
	ressources.NanoCPUs = getCPUCount(containerNameToShow)

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        volumeBindings,
		Resources:    ressources,
		Privileged:   privilegedContainer,
	}

	log.Printf("Running cluster %s (%s Memory / %d CPU) ...", containerNameToShow, getMemorySize(containerNameToShow), ressources.NanoCPUs)

	resp, err := getDocker().ContainerCreate(ctx, config, hostConfig, nil, containerName)
	if err != nil {
		log.Fatal(err)
	}

	err = getDocker().ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	// The if removes the error:
	//panic: runtime error: invalid memory address or nil pointer dereference
	//[signal SIGSEGV: segmentation violation code=0x1 addr=0x20 pc=0x137a2b4]
	if err != nil {
		log.Fatal(err)
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
func startContainer(containerName string) {
	if err := getDocker().ContainerStart(ctx, containerName, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}
}
