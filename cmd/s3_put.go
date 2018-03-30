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
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdPut,
		DisableFlagsInUseLine: true,
	}

	return cmd
}

// S3CmdPut wraps s3cmd command in the container
func S3CmdPut(cmd *cobra.Command, args []string) {
	notExistCheck()
	notRunningCheck()
	dir := dockerInspect("bind")
	fileName := args[0]
	bucketName := args[1]
	fileNameBase := path.Base(fileName)

	if _, err := os.Stat(dir + "/" + fileNameBase); os.IsNotExist(err) {
		_, err := copyFile(fileName, dir+"/"+fileNameBase)
		if err != nil {
			log.Fatal(err)
		}
	}

	command := []string{"s3cmd", "put", TempPath + fileNameBase, "s3://" + bucketName}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
