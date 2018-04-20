package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/jmoiron/jsonq"
	"golang.org/x/sys/unix"
)

// getSeLinuxStatus gets SeLinux status
func getSeLinuxStatus() string {
	testBinaryExist("getenforce")

	out, err := exec.Command("getenforce").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

// applySeLinuxLabel checks if SeLinux is installed and set to Enforcing,
// we relabel our workingDirectory to allow the container to access files in this directory
func applySeLinuxLabel(dir string) {
	testBinaryExist("getenforce")

	selinuxStatus := getSeLinuxStatus()
	lines := strings.Split(selinuxStatus, "\n")
	for _, l := range lines {
		if len(l) <= 0 {
			// Ignore empty line.
			continue
		}
		if l == "Enforcing" {
			meUserName, meID := whoAmI()
			if meID != "0" {
				log.Fatal("Hey " + meUserName + "! Run me as 'root' so I can apply the right SeLinux label on " + dir)
			}
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.Mkdir(dir, 0755)
			}
			testBinaryExist("chcon")
			cmd := "chcon " + " -Rt" + " svirt_sandbox_file_t " + dir
			_, err := exec.Command("chcon", "-Rt", "svirt_sandbox_file_t", dir).Output()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Executing: " + cmd)
		}
	}
}

// byLastOctetValue implements sort.Interface used in sorting a list
// of ip address by their last octet value.
type byLastOctetValue []net.IP

func (n byLastOctetValue) Len() int      { return len(n) }
func (n byLastOctetValue) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n byLastOctetValue) Less(i, j int) bool {
	return []byte(n[i].To4())[3] < []byte(n[j].To4())[3]
}

// getInterfaceIPv4s is synonymous to net.InterfaceAddrs()
// returns net.IP IPv4 only representation of the net.Addr.
// Additionally the returned list is sorted by their last
// octet value.
//
// [The logic to sort by last octet is implemented to
// prefer CIDRs with higher octets, this in-turn skips the
// localhost/loopback address to be not preferred as the
// first ip on the list. Subsequently this list helps us print
// a user friendly message with appropriate values].
func getInterfaceIPv4s() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("Unable to determine network interface address. %s", err)
	}
	// Go through each return network address and collate IPv4 addresses.
	var nips []net.IP
	for _, addr := range addrs {
		if addr.Network() == "ip+net" {
			var nip net.IP
			// Attempt to parse the addr through CIDR.
			nip, _, err = net.ParseCIDR(addr.String())
			if err != nil {
				return nil, fmt.Errorf("Unable to parse address %s, error %s", addr, err)
			}
			// Collect only IPv4 addrs.
			if nip.To4() != nil {
				nips = append(nips, nip)
			}
		}
	}
	// Sort the list of IPs by their last octet value.
	sort.Sort(sort.Reverse(byLastOctetValue(nips)))
	return nips, nil
}

// execContainer execs a given command inside the container
func execContainer(containerName string, cmd []string) []byte {
	optionsCreate := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	response, err := getDocker().ContainerExecCreate(ctx, containerName, optionsCreate)
	if err != nil {
		log.Fatal(err)
	}

	optionsAttach := types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	}
	connection, err := getDocker().ContainerExecAttach(ctx, response.ID, optionsAttach)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()
	output, err := ioutil.ReadAll(connection.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Remove 8 first characters to get a readable content
	// Sometimes the command returns nothing, without the following if the program fails without
	// runtime error: slice bounds out of range
	if len(output) > 0 {
		return output[8:]
	}
	return nil
}

// grepForSuccess searches for the word 'SUCCESS' inside the container logs
func grepForSuccess(containerName string) bool {
	out, err := getDocker().ContainerLogs(ctx, containerName, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	newStr := buf.String()

	if strings.Contains(newStr, "SUCCESS") {
		return true
	}
	return false
}

// cephNanoHealth loops on grepForSuccess for 30 seconds, fails after.
func cephNanoHealth(containerName string) {
	// setting timeout values
	timeout := 60
	poll := 0

	// wait for 60sec to validate that the container started properly
	for poll < timeout {
		if grepForSuccess(containerName) {
			return
		}
		time.Sleep(time.Second * 1)
		poll++
	}

	// if we reach here, something is broken in the container
	log.Println("The container " + containerName + " never reached a clean state. Showing the container logs now:")
	// ideally we would return the second value of GrepForSuccess when it's false
	// this would mean having 2 return values for GrepForSuccess
	out, err := getDocker().ContainerLogs(ctx, containerName, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	newStr := buf.String()
	fmt.Println(newStr)
	log.Fatal("Please open an issue at: https://github.com/ceph/cn with the logs above.")
}

// curlTestURL tests a given URL
func curlTestURL(url string) bool {
	response, err := http.Get(url)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	if _, err := ioutil.ReadAll(response.Body); err != nil {
		return false
	}
	return true
}

// curlURL queries a given URL and returns its content
func curlURL(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Println("URL " + url + " is unreachable.")
		log.Fatal(err)
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

// countTags queries the number of tags
func countTags() int {
	var url string
	data := map[string]interface{}{}
	url = "https://registry.hub.docker.com/v2/repositories/ceph/daemon/tags/"
	output := curlURL(url)
	dec := json.NewDecoder(strings.NewReader(string(output)))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	tagCount, _ := jq.Int("count")
	return tagCount
}

func pageCount() int {
	tagCount := countTags()
	pageCount := tagCount / 10
	return int(pageCount)
}

// parseMap parses a json element
// re-adapted code from:
// https://stackoverflow.com/questions/29366038/looping-iterate-over-the-second-level-nested-json-in-go-lang
func parseMap(aMap map[string]interface{}, keyType string) {
	for key, val := range aMap {
		switch concreteVal := val.(type) {
		case []interface{}:
			parseArray(val.([]interface{}), keyType)
		default:
			if key == keyType {
				fmt.Println(concreteVal)
			}
		}
	}
}

// parseArray parses json array
// re-adapted code from:
// https://stackoverflow.com/questions/29366038/looping-iterate-over-the-second-level-nested-json-in-go-lang
func parseArray(anArray []interface{}, keyType string) {
	for _, val := range anArray {
		switch concreteVal := val.(type) {
		case map[string]interface{}:
			parseMap(val.(map[string]interface{}), keyType)
		default:
			fmt.Println(concreteVal)
		}
	}
}

// CephNanoS3Health loops for 30 seconds while testing Ceph RGW health
func cephNanoS3Health(containerName string, rgwPort string) {
	// setting timeout
	timeout := 30
	poll := 0
	ips, _ := getInterfaceIPv4s()
	// Taking the first IP is probably not ideal
	// IMHO, using the interface with most of the traffic is better
	url := "http://" + ips[0].String() + ":" + rgwPort

	for poll < timeout {
		if curlTestURL(url) {
			return
		}
		time.Sleep(time.Second * 1)
		poll++
	}
	log.Println("Timeout while trying to reach: " + url)
	log.Println("S3 gateway for cluster " + containerName + " is not responding. Showing S3 logs:")
	showS3Logs(containerName)
	log.Fatal("Please open an issue at: https://github.com/ceph/cn.")
}

// echoInfo prints useful information about Ceph Nano
func echoInfo(containerName string) {
	// Get listening port
	rgwPort := dockerInspect(containerName, "PortBindings")

	// Always wait the container to be ready
	cephNanoHealth(containerName)
	cephNanoS3Health(containerName, rgwPort)

	// Fetch Amazon Keys
	cephNanoAccessKey, cephNanoSecretKey := getAwsKey(containerName)

	// Get Ceph health
	cmd := []string{"ceph", "health"}
	c := execContainer(containerName, cmd)

	// Get IPs, later using the first IP of the list is not ideal
	// However, Docker binds RGW port on 0.0.0.0 so any address will work
	ips, _ := getInterfaceIPv4s()

	// Get the working directory
	dir := dockerInspect(containerName, "Binds")

	infoLine :=
		"\n" + strings.TrimSpace(string(c)) + " is the Ceph status \n" +
			"S3 object server address is: http://" + ips[0].String() + ":" + rgwPort + "\n" +
			"S3 user is: nano \n" +
			"S3 access key is: " + cephNanoAccessKey + "\n" +
			"S3 secret key is: " + cephNanoSecretKey + "\n" +
			"Your working directory is: " + dir + "\n"
	fmt.Println(infoLine)
}

// getAwsKey gets AWS keys from inside the container
func getAwsKey(containerName string) (string, string) {
	cmd := []string{"cat", "/nano_user_details"}

	output := execContainer(containerName, cmd)

	// declare structures for json
	type s3Details []struct {
		AccessKey string `json:"Access_key"`
		SecretKey string `json:"Secret_key"`
	}
	type jason struct {
		Keys s3Details
	}
	// assign variable to our json struct
	var parsedMap jason

	json.Unmarshal(output, &parsedMap)

	cephNanoAccessKey := parsedMap.Keys[0].AccessKey
	cephNanoSecretKey := parsedMap.Keys[0].SecretKey
	return cephNanoAccessKey, cephNanoSecretKey
}

// dockerInspect inspects the container Binds
func dockerInspect(containerName string, pattern string) string {
	inspect, err := getDocker().ContainerInspect(ctx, containerName)
	if err != nil {
		log.Fatal(err)
	}

	if pattern == "Binds" {
		parts := strings.Split(inspect.HostConfig.Binds[0], ":")
		return parts[0]
	}

	if pattern == "PortBindings" {
		parts := strings.Split(inspect.Config.Env[0], "=")
		return parts[1]
	}

	// The part is helpful when passing a dedicated directory to store Ceph's data
	// We look for bindmounts, if we find more than 1 (the first one is the work-dir)
	// then this means we passed a dedicated directory, this is used by the purge function
	// to remove the OSD data content once we purge the cluster
	if pattern == "BindsData" {
		parts := inspect.HostConfig.Binds
		if len(inspect.HostConfig.Binds) >= 2 {
			parts = strings.Split(parts[1], ":")
		} else {
			return "noDataDir"
		}
		return parts[1]
	}

	// this assumes a default that we are looking for the image name
	parts := inspect.Config.Image
	return parts
}

// inspectImage inspects a given image
func inspectImage(ImageID string, dataType string) string {
	i, _, err := getDocker().ImageInspectWithRaw(ctx, ImageID)
	if err != nil {
		// sometimes the image does not exist anymore, we want to report that
		return "image is not present, did you remove it?"
	}
	if dataType == "tag" {
		// If the tag disappeared, probably because a newer tag with a same appeared
		// Let's return RepoDigests
		if len(i.RepoTags) == 0 {
			return strings.Join(i.RepoDigests, "")
		}
		return i.RepoTags[0]
	}
	if dataType == "created" {
		return i.Created
	}
	if len(i.ContainerConfig.Labels["RELEASE"]) == 0 {
		return "unknown image release, are you running an official image?"
	}
	return i.ContainerConfig.Labels["RELEASE"]
}

// pullImage downloads the container image
func pullImage() bool {
	_, _, err := getDocker().ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		fmt.Print("The container image is not present, pulling it. \n" +
			"This operation can take a few minutes.")

		out, err := getDocker().ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			// the error message will appear on a new line after the info above
			log.Println()
			log.Fatal(err)
		}

		reader := bufio.NewReader(out)
		defer out.Close() // pullResp is io.ReadCloser
		var respo bytes.Buffer
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				// it could be EOF or read error
				break
			}
			respo.Write(line)
			respo.WriteByte('\n')
			fmt.Print(".")
		}
		fmt.Println("")
		return true
	}
	return false
}

func notExistCheck(containerName string) {
	containerNameToShow := containerName[len(containerNamePrefix):]

	if (!containerStatus(containerName, false, "running")) && (!containerStatus(containerName, false, "exited")) {
		log.Println("Cluster " + containerNameToShow + " does not exist yet.")
		os.Exit(0)
	}
}

func notRunningCheck(containerName string) {
	containerNameToShow := containerName[len(containerNamePrefix):]

	if status := containerStatus(containerName, true, "exited"); status {
		log.Println("Cluster " + containerNameToShow + " is not running.")
		os.Exit(0)
	}
}

func copyFile(srcName, dstName string) (int64, error) {
	src, e := os.Open(srcName)
	if e != nil {
		return 0, errors.New("Error while opening file for reading. Caused by: " + e.Error())
	}

	dst, e := os.Create(dstName)
	if e != nil {
		src.Close()
		return 0, errors.New("Error while opening file for writing. Caused by: " + e.Error())
	}

	numBytesWritten, e := io.Copy(dst, src)
	if e != nil {
		dst.Close()
		src.Close()
		return 0, errors.New("Error while copying. Caused by: " + e.Error())
	}

	e = dst.Close()
	if e != nil {
		src.Close()
		return numBytesWritten, errors.New("Error while closing. Caused by: " + e.Error())
	}

	e = src.Close()
	if e != nil {
		return numBytesWritten, errors.New("Error while closing. Caused by: " + e.Error())
	}

	return numBytesWritten, nil
}

// copyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return errors.New("Error can not stat source. Caused by: " + err.Error())
	}
	if !si.IsDir() {
		return errors.New("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return errors.New("Error can not stat destination. Caused by: " + err.Error())
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return errors.New("Error can not create directories. Caused by: " + err.Error())
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.New("Error can not read directories. Caused by: " + err.Error())
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return errors.New("Error copying directory. Caused by: " + err.Error())
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			_, err = copyFile(srcPath, dstPath)
			if err != nil {
				return errors.New("Error copying file. Caused by: " + err.Error())
			}
		}
	}

	return nil
}

// checkPortInUsed checks if a port is in-used
func checkPortInUsed(portNum string) bool {
	hostName := "0.0.0.0"
	seconds := 1
	timeOut := time.Duration(seconds) * time.Second

	_, err := net.DialTimeout("tcp", net.JoinHostPort(hostName, portNum), timeOut)

	// if there is an error this means the port is not used
	// and the connection can not be established
	if err != nil {
		return true
	}
	return false
}

// generateRGWPortToUse generates the binding port for Ceph Rados Gateway
func generateRGWPortToUse() string {
	maxPort := 8100
	for i := 8000; i <= maxPort; i++ {
		portNumStr := fmt.Sprint(i)
		status := checkPortInUsed(portNumStr)
		if status {
			return portNumStr
		}
	}
	return "notfound"
}

// getFileType checks wether a specified data is directory, a block device or something else
// function borrowed from https://github.com/andrewsykim/kubernetes/blob/2deb7af9b248a7ddc00e61fcd08aa9ea8d2d09cc/pkg/util/mount/mount_linux.go#L416
func getFileType(pathname string) (string, error) {
	finfo, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		return "notfound", fmt.Errorf("path %q does not exist", pathname)
	}
	// err in call to os.Stat
	if err != nil {
		return "error", err
	}

	mode := finfo.Sys().(*syscall.Stat_t).Mode
	switch mode & syscall.S_IFMT {
	case syscall.S_IFSOCK:
		return "socket", nil
	case syscall.S_IFBLK:
		return "blockdev", nil
	case syscall.S_IFCHR:
		return "chardev", nil
	case syscall.S_IFDIR:
		return "directory", nil
	case syscall.S_IFREG:
		return "file", nil
	}

	return "error", fmt.Errorf("only recognize file, directory, socket, block device and character device")
}

// testBinaryExist tests if a binary is present on the system
func testBinaryExist(binary string) bool {
	binary, err := exec.LookPath(binary)
	if err != nil {
		log.Fatal(binary + " is not installed!")
	}
	return true
}

// getDiskFormat returns information about a disk such as filesystem and partition table
func getDiskFormat(disk string) string {
	testBinaryExist("blkid")

	out, err := exec.Command("blkid", "-p", "-s", "TYPE", "-s", "PTTYPE", "-o", "export", disk).Output()
	if err != nil {
		// Disk device does not have any label or is a partition
		// For `blkid`, if the specified token (TYPE/PTTYPE, etc) was
		// not found, or no (specified) devices could be identified, an
		// exit code of 2 is returned.
		// We are not 100% that this is the problem, but the code called before this already made sure
		// that the device was a block special device, so if blkid fails that's probably because there is
		// nothing to see. My first assumption is that the device does not have a partition label OR
		// is a partition
		log.Fatal("I suspect either the disk" + disk + " has no partition label or is a partition.\n" +
			"If you gave me a whole device, make sure it has a partition table (e.g: gpt). \n" +
			"If you gave me a partition, I don't support partitions yet, give me a whole device.\n" +
			"As an alternative, you can create a filesystem on this partition and give the mountpoint to me.\n" +
			"\nAlso if the disk was an OSD you need to zap it (e.g: with 'ceph-disk zap').")
	}
	return string(out)
}

// getDiskPartitions returns the list of partitions on a disk
func getDiskPartitions(disk string) []string {
	testBinaryExist("parted")

	out, err := exec.Command("parted", "-s", "-m", disk, "print").Output()
	if err != nil {
		log.Fatal(err)
	}
	var partitions []string
	lines := strings.Split(string(out), "\n")
	// build a slice
	for _, l := range lines {
		if len(l) <= 0 {
			// Ignore empty line.
			continue
		}
		partitions = append(partitions, l)
	}
	return partitions
}

// exclusiveOpenFailsOnDevice tries to open a device with O_EXCL flag
// stolen and re-adapted from https://github.com/kubernetes/kubernetes/blob/77d18dbad9d2abea7d7b7b22be02fb422e03f0a9/pkg/util/mount/mount_linux.go#L306
func exclusiveOpenFailsOnDevice(pathname string) (bool, error) {
	fd, errno := unix.Open(pathname, unix.O_RDONLY|unix.O_EXCL, 0)
	// If the device is in use, open will return an invalid fd.
	// When this happens, it is expected that Close will fail and throw an error.
	defer unix.Close(fd)
	if errno == nil {
		// device not in use
		return false, nil
	} else if errno == unix.EBUSY {
		// device is in use
		return true, nil
	}
	// error during call to Open
	return false, errno
}

// isEmpty tests if a directory is empty or not
func isEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
}

// whoAmI returns the id of the user running cn
func whoAmI() (string, string) {
	me, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return me.Username, me.Uid
}
