package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// CliS3CmdSync is the Cobra CLI call
func CliS3CmdSync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync LOCAL_DIR BUCKET",
		Short: "Synchronize a directory tree to S3",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdSync,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdSync wraps s3cmd command in the container
func S3CmdSync(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	localDir := args[0]
	bucketName := args[1]
	dir := dockerInspect()
	destDir := "/tmp"

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
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
