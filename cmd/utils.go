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
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
)

// validateEnv verifies the ability to run the program
func validateEnv() {
	seLinux()
}

// seLinux checks if Selinux is installed and set to Enforcing,
// we relabel our WorkingDirectory to allow the container to access files in this directory
func seLinux() {
	if _, err := os.Stat("/sbin/getenforce"); !os.IsNotExist(err) {
		out, err := exec.Command("getenforce").Output()
		if err != nil {
			log.Fatal(err)
		}
		if string(out) == "Enforcing" {
			if _, err := os.Stat(WorkingDirectory); os.IsNotExist(err) {
				os.Mkdir(WorkingDirectory, 0755)
			}
			exec.Command("sudo chcon -Rt svirt_sandbox_file_t %s", WorkingDirectory)
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
// prefer CIDRs with higher octects, this in-turn skips the
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
				return nil, fmt.Errorf("Unable to parse addrss %s, error %s", addr, err)
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
func execContainer(ContainerName string, cmd []string) []byte {
	optionsCreate := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	response, err := getDocker().ContainerExecCreate(ctx, ContainerName, optionsCreate)
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

// grepForSuccess searchs for the word 'SUCCESS' inside the container logs
func grepForSuccess() bool {
	out, err := getDocker().ContainerLogs(ctx, ContainerName, types.ContainerLogsOptions{ShowStdout: true})
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
func cephNanoHealth() {
	// setting timeout values
	timeout := 60
	poll := 0

	// wait for 60sec to validate that the container started properly
	for poll < timeout {
		if grepForSuccess() {
			return
		}
		time.Sleep(time.Second * 1)
		poll++
	}

	// if we reach here, something is broken in the container
	fmt.Print("The container from Ceph Nano never reached a clean state. Show the container logs:")
	// ideally we would return the second value of GrepForSuccess when it's false
	// this would mean having 2 return values for GrepForSuccess
	out, err := getDocker().ContainerLogs(ctx, ContainerName, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	newStr := buf.String()
	fmt.Println(newStr)
	log.Fatal("Please open an issue at: https://github.com/ceph/cn.")
}

// curlS3 queries S3 URL
func curlS3() bool {
	ips, _ := getInterfaceIPv4s()
	// Taking the first IP is probably not ideal
	// IMHO, using the interface with most of the traffic is better
	var url string
	url = "http://" + ips[0].String() + ":" + RgwPort

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

// CephNanoS3Health loops for 30 seconds while testing Ceph RGW heatlh
func cephNanoS3Health() {
	// setting timeout
	timeout := 30
	poll := 0

	for poll < timeout {
		if curlS3() {
			return
		}
		time.Sleep(time.Second * 1)
		poll++
	}
	fmt.Println("S3 gateway is not responding. Showing S3 logs:")
	showS3Logs()
	log.Fatal("Please open an issue at: https://github.com/ceph/cn.")
}

// echoInfo prints useful information about Ceph Nano
func echoInfo() {
	// Always wait the container to be ready
	cephNanoHealth()
	cephNanoS3Health()

	// Fetch Amazon Keys
	CephNanoAccessKey, CephNanoSecretKey := getAwsKey()

	// Get Ceph health
	cmd := []string{"ceph", "health"}
	c := execContainer(ContainerName, cmd)

	// Get IPs, later using the first IP of the list is not ideal
	// However, Docker binds RGW port on 0.0.0.0 so any address will work
	ips, _ := getInterfaceIPv4s()

	// Get the working directory
	dir := dockerInspect("bind")

	InfoLine :=
		"\n" + strings.TrimSpace(string(c)) + " is the Ceph status \n" +
			"S3 object server address is: http://" + ips[0].String() + ":" + RgwPort + "\n" +
			"S3 user is: nano \n" +
			"S3 access key is: " + CephNanoAccessKey + "\n" +
			"S3 secret key is: " + CephNanoSecretKey + "\n" +
			"Your working directory is: " + dir + "\n"
	fmt.Println(InfoLine)
}

// getAwsKey gets AWS keys from inside the container
func getAwsKey() (string, string) {
	cmd := []string{"cat", "/nano_user_details"}

	output := execContainer(ContainerName, cmd)

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

	CephNanoAccessKey := parsedMap.Keys[0].AccessKey
	CephNanoSecretKey := parsedMap.Keys[0].SecretKey
	return CephNanoAccessKey, CephNanoSecretKey
}

// dockerInspect inspect the container Binds
func dockerInspect(pattern string) string {
	inspect, err := getDocker().ContainerInspect(ctx, ContainerName)
	if err != nil {
		log.Fatal(err)
	}

	if pattern == "bind" {
		parts := strings.Split(inspect.HostConfig.Binds[0], ":")
		return parts[0]
	}
	// this assumes a default that we are looking for the image name
	parts := inspect.Config.Image
	return parts
}

// inspectImage inspect a given image
func inspectImage() map[string]string {
	ImageName = dockerInspect("image")
	i, _, err := getDocker().ImageInspectWithRaw(ctx, ImageName)
	if err != nil {
		var m map[string]string
		m = make(map[string]string)
		m["head"] = "unknown"
		return m
	}
	return i.Config.Labels
}

// pullImage downloads the container image
func pullImage() bool {
	_, _, err := getDocker().ImageInspectWithRaw(ctx, ImageName)
	if err != nil {
		fmt.Print("The container image is not present, pulling it. \n" +
			"This operation can take a few minutes.")

		out, err := getDocker().ImagePull(ctx, ImageName, types.ImagePullOptions{})
		if err != nil {
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

func notExistCheck() {
	if (!containerStatus(false, "running")) && (!containerStatus(false, "exited")) {
		fmt.Println("ceph-nano does not exist yet.")
		os.Exit(0)
	}
}

func notRunningCheck() {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running.")
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
