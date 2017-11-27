package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// CliS3CmdPut is the Cobra CLI call
func CliS3CmdPut() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put FILE BUCKET",
		Short: "Put file into bucket",
		Long: `When you put a file into a bucket, NEVER use the fullpath of the file because Ceph Nano already hardcodes the path of your work directory.
This means you should be working for within your working directory.
Then if you want to put an object simply do 'cn s3 put FILE' instead of cn s3 put <working-directory/FILE.
If you wan to upload a file that is in a directory withing your working directory simply do 'cn s3 put DIRECTORY/FILE.'`,
		Args: cobra.ExactArgs(2),
		Run:  S3CmdPut,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdPut wraps s3cmd command in the container
func S3CmdPut(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	dir := dockerInspect()
	fileName := args[0]
	bucketName := args[1]
	fileNameBase := path.Base(fileName)

	if _, err := os.Stat(dir + "/" + fileNameBase); os.IsNotExist(err) {
		_, err := copyFile(fileName, dir+"/"+fileNameBase)
		if err != nil {
			log.Fatal(err)
		}
	}

	command := []string{"s3cmd", "put", "/tmp/" + fileNameBase, "s3://" + bucketName}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
