package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdPut is the Cobra CLI call
func cliS3CmdPut() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put CLUSTER FILE BUCKET",
		Short: "Put file into bucket",
		Args:  cobra.ExactArgs(3),
		Run:   S3CmdPut,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// S3CmdPut wraps s3cmd command in the container
func S3CmdPut(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]

	notExistCheck(containerName)
	notRunningCheck(containerName)
	dir := dockerInspect(containerName, "Binds")
	fileName := args[1]
	bucketName := args[2]
	fileNameBase := path.Base(fileName)

	if _, err := os.Stat(dir + "/" + fileNameBase); os.IsNotExist(err) {
		_, err := copyFile(fileName, dir+"/"+fileNameBase)
		if err != nil {
			log.Fatal(err)
		}
	}

	command := []string{"s3cmd", "put", tempPath + fileNameBase, "s3://" + bucketName}
	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + args[0]
	fmt.Println(output)
}
