package client

import (
	"fmt"
	"os/exec"
)

func DownloadBackup(bucketName string, filename string) error {
	args := []string{
		"s3",
		"cp",
		fmt.Sprintf("s3://%s/%s.bak", bucketName, filename),
		fmt.Sprintf("%s.bak", filename),
	}

	_, err := executeCommand(args)

	return err
}

func executeCommand(args []string) (string, error) {
	byteOutput, err := exec.Command("aws", args...).Output()
	return string(byteOutput), err
}
