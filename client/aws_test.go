package client_test

import (
	"testing"

	"github.com/alexhokl/rds-backup/client"
)

type downloadTest struct {
	bucketName        string
	filename          string
	downloadDirectory string
	commands          []string
}

var isAwsCliInstalledTests = []test{
	{[]string{"aws help"}},
}

var isAwsCredentialsConfiguredTests = []test{
	{[]string{"aws s3 ls"}},
}

var downloadBackupTests = []downloadTest{
	{"my-bucket", "backup.1.bak", "/home/user/backups", []string{"aws s3 cp s3://my-bucket/backup.1.bak /home/user/backups/backup.1.bak"}},
}

func TestIsAwsCliInstalled(t *testing.T) {
	for i, test := range isAwsCliInstalledTests {
		cmdLine := &MockCommandLine{}
		client.IsAwsCliInstalled(cmdLine)
		testCommands(t, i, cmdLine.Commands, test.commands)
	}
}

func TestIsAwsCredentialsConfigured(t *testing.T) {
	for i, test := range isAwsCredentialsConfiguredTests {
		cmdLine := &MockCommandLine{}
		client.IsAwsCredentialsConfigured(cmdLine)
		testCommands(t, i, cmdLine.Commands, test.commands)
	}
}

func TestDownloadBackup(t *testing.T) {
	for i, test := range downloadBackupTests {
		cmdLine := &MockCommandLine{}
		client.DownloadBackup(cmdLine, test.bucketName, test.filename, test.downloadDirectory)
		testCommands(t, i, cmdLine.Commands, test.commands)
	}
}
