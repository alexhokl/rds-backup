package client

import (
	"fmt"
	"os/exec"

	"github.com/spf13/viper"
)

func DownloadBackup(bucketName string, filename string) error {
	args := []string{
		"s3",
		"cp",
		fmt.Sprintf("s3://%s/%s", bucketName, filename),
		filename,
	}

	_, err := executeCommand(args)

	return err
}

func executeCommand(args []string) (string, error) {
	byteOutput, err := exec.Command("aws", args...).Output()

	if viper.GetBool("verbose") {
		fmt.Println("Command executed:", "aws", args)
	}

	return string(byteOutput), err
}
