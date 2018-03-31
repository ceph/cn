package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// S3CmdSkip means do not do anything when object exists
	S3CmdSkip bool

	// S3CmdContinue means
	S3CmdContinue bool

	// S3CmdOpt is the option to apply when trying to get a file and the destination already exists
	S3CmdOpt string
)

// CliS3CmdGet is the Cobra CLI call
func CliS3CmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get CLUSTER BUCKET/OBJECT [LOCAL_FILE]",
		Short: "Get file out of a bucket",
		Args:  cobra.RangeArgs(2, 3),
		Run:   S3CmdGet,
		DisableFlagsInUseLine: false,
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().BoolVarP(&S3CmdSkip, "skip", "s", true, "Skip over files that exist at the destination")
	cmd.Flags().BoolVarP(&S3CmdForce, "force", "f", false, "Force overwrite files that exist at the destination")
	cmd.Flags().BoolVarP(&S3CmdContinue, "continue", "c", false, "Continue getting a partially downloaded file")

	return cmd
}

// S3CmdGet wraps s3cmd command in the container
func S3CmdGet(cmd *cobra.Command, args []string) {
	ContainerName := ContainerNamePrefix + args[0]

	notExistCheck(ContainerName)
	notRunningCheck(ContainerName)
	BucketObjectName := args[1]
	var fileName string

	if len(args) > 1 {
		fileName = args[2]
	} else {
		fileName = BucketObjectName
	}

	BucketObjectNameBase := path.Base(BucketObjectName)

	if S3CmdForce {
		S3CmdOpt = "--force"
	} else if S3CmdSkip {
		S3CmdOpt = "--skip-existing"
	} else if S3CmdContinue {
		S3CmdOpt = "--continue"
	}
	// if args
	command := []string{"s3cmd", "get", S3CmdOpt, "s3://" + BucketObjectName, TempPath}
	output := strings.TrimSuffix(string(execContainer(ContainerName, command)), "\n") + " on cluster " + ContainerName

	dir := dockerInspect(ContainerName, "Binds")
	if fileName != BucketObjectName {
		//if _, err := os.Stat(fileName); os.Stat.Mode.IsDir(err) {
		if info, err := os.Stat(fileName); err == nil && info.IsDir() {
			fileName = fileName + "/" + BucketObjectNameBase
		}
		_, err := copyFile(dir+"/"+BucketObjectNameBase, fileName)
		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println(output)
}
