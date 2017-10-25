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

func IsAwsCliInstalled() bool {
	_, err := executeCommand([]string{"help"})
	return err == nil
}

func IsAwsCredentialsConfigured() bool {
	_, err := executeCommand([]string{"s3", "ls"})
	return err == nil
}

func executeCommand(args []string) (string, error) {
	if viper.GetBool("verbose") {
		fmt.Println("Command executed:", "aws", args)
	}
	byteOutput, err := exec.Command("aws", args...).Output()
	return string(byteOutput), err
}
