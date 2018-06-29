package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// cliS3CmdSync is the Cobra CLI call
func cliS3CmdSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync CLUSTER LOCAL_DIR BUCKET",
		Short: "Synchronize a directory tree to S3",
		Args:  cobra.ExactArgs(3),
		Run:   S3CmdSync,
		DisableFlagsInUseLine: true,
	}
	cmd.Flags().BoolVarP(&debugS3, "debug", "d", false, "Run S3 commands in debug mode")

	return cmd
}

// S3CmdSync wraps s3cmd command in the container
func S3CmdSync(cmd *cobra.Command, args []string) {
	containerName := containerNamePrefix + args[0]
	containerNameToShow := containerName[len(containerNamePrefix):]

	notExistCheck(containerName)
	notRunningCheck(containerName)
	localDir := args[1]
	bucketName := args[2]
	dir := dockerInspect(containerName, "Binds")
	destDir := tempPath

	if localDir != dir {
		destDir = dir + "/" + localDir
		err := copyDir(localDir, destDir)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Syncing directory '%s' in the '%s' bucket. \n"+
		"It might take some time depending on the amount of data. \n"+
		"Do not expect any output until the upload is finished. \n \n", localDir, bucketName)

	command := []string{"s3cmd", "sync", destDir, "s3://" + bucketName}
	if debugS3 {
		command = append(command, "--debug")
	}

	output := strings.TrimSuffix(string(execContainer(containerName, command)), "\n") + " on cluster " + containerNameToShow
	fmt.Println(output)
}
