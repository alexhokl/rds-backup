package client_test

import (
	"testing"

	"github.com/alexhokl/rds-backup/client"
)

type testGetStatus struct {
	params   client.DatabaseParameters
	taskID   string
	commands []string
}

var isEnvironmentSatisfiedTests = []test{
	{[]string{"sqlcmd -?"}},
}

// var getStatusTests = []testGetStatus{
// 	{
// 		client.DatabaseParameters{
// 			Server:       "myserver.com",
// 			Username:     "me",
// 			Password:     "password",
// 			DatabaseName: "dbname",
// 		},
// 		"",
// 		[]string{`sqlcmd -S myserver.com -d dbname -U me -P password -Q SET NOCOUNT ON

// 	DECLARE @s TABLE (
// 	task_id INT,
// 	task_type VARCHAR(20),
// 	database_name VARCHAR(20),
// 	complete INT,
// 	duration INT,
// 	lifecycle VARCHAR(20),
// 	task_info VARCHAR(MAX),
// 	last_updated DATETIME,
// 	created_at DATETIME,
// 	S3_object_arn VARCHAR(MAX),
// 	overwrite_S3_backup_file BIT,
// 	KMS_master_key_arn VARCHAR(100)
// )

// 	INSERT INTO @s
// 	exec msdb.dbo.rds_task_status @db_name='dbname'

// 	SELECT TOP 1 lifecycle FROM @s

// 	SET NOCOUNT OFF`
// 		},
// 	},
// }

func TestIsEnvironmentSatisfied(t *testing.T) {
	for i, test := range isEnvironmentSatisfiedTests {
		cmdLine := &MockCommandLine{}
		c := &client.NativeClient{}
		c.IsEnvironmentSatisfied(cmdLine)
		testCommands(t, i, cmdLine.Commands, test.commands)
	}
}

// func TestGetStatus(t *testing.T) {
// 	for i, test := range getStatusTests {
// 		cmdLine := &MockCommandLine{}
// 		c := &client.NativeClient{}
// 		c.GetStatus(cmdLine, &test.params, test.taskID)
// 		testCommands(t, i, cmdLine.Commands, test.commands)
// 	}
// }
