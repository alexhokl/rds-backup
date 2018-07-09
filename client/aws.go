package client

import (
	"fmt"
	"os"
	"path/filepath"
)

// DownloadBackup downloads a SQL backup from a S3 bucket
func DownloadBackup(c Command, bucketName string, filename string, downloadDirectory string) error {
	currentDirectory, _ := os.Getwd()
	pathToBak := filepath.Join(currentDirectory, filename)
	if downloadDirectory != "" {
		pathToBak = filepath.Join(downloadDirectory, filename)
	}
	args := []string{
		"s3",
		"cp",
		fmt.Sprintf("s3://%s/%s", bucketName, filename),
		pathToBak,
	}

	fmt.Printf("Download of backup from AWS S3 (s3://%s/%s) started...\n", bucketName, filename)
	_, err := executeCommand(c, args)

	if err == nil {
		fmt.Printf("Download of the backup has been completed (%s)\n", pathToBak)
	}

	return err
}

// IsAwsCliInstalled returns if AWS CLI has been installed
func IsAwsCliInstalled(c Command) bool {
	_, err := executeCommand(c, []string{"help"})
	return err == nil
}

// IsAwsCredentialsConfigured returns if AWS CLI credentials has been configured
func IsAwsCredentialsConfigured(c Command) bool {
	_, err := executeCommand(c, []string{"s3", "ls"})
	return err == nil
}

func executeCommand(c Command, args []string) (string, error) {
	return c.Execute("aws", args)
}
