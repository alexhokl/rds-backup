package client

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// DatabaseParameters contains the database information
type DatabaseParameters struct {
	Server       string
	Username     string
	Password     string
	DatabaseName string
}

// BackupParameters contains database and destination bucket information
type BackupParameters struct {
	DatabaseParameters
	BucketName string
	Filename   string
}

const statusTableDeclaration = `DECLARE @s TABLE (
	task_id INT,
	task_type VARCHAR(20),
	database_name VARCHAR(20),
	complete INT,
	duration INT,
	lifecycle VARCHAR(20),
	task_info VARCHAR(MAX),
	last_updated DATETIME,
	created_at DATETIME,
	S3_object_arn VARCHAR(MAX),
	overwrite_S3_backup_file BIT,
	KMS_master_key_arn VARCHAR(100)
)`

const createTableDeclaration = `DECLARE @s TABLE (
	task_id INT,
	task_type VARCHAR(20),
	lifecycle VARCHAR(20),
	created_at DATETIME,
	last_updated DATETIME,
	database_name VARCHAR(20),
	S3_object_arn VARCHAR(MAX),
	overwrite_S3_backup_file BIT,
	KMS_master_key_arn VARCHAR(100),
	task_progress INT,
	task_info VARCHAR(MAX)
)`

// GetStatus returns the status of the latest backup
func GetStatus(params *DatabaseParameters, taskID string) (string, error) {
	query := "SELECT TOP 1 lifecycle FROM @s"
	if taskID != "" {
		query = fmt.Sprintf("SELECT lifecycle FROM @s WHERE task_id = %s", taskID)
	}

	statement := fmt.Sprintf(`SET NOCOUNT ON

	%s

	INSERT INTO @s
	exec msdb.dbo.rds_task_status @db_name='%s'

	%s

	SET NOCOUNT OFF`, statusTableDeclaration, params.DatabaseName, query)

	args := getCommandArgs(params, statement)
	output, err := execute(args)
	if err != nil {
		return "", err
	}
	return getSQLOutput(output), nil
}

// GetCompletionPercentage returns the percentage of completion of the latest backup
func GetCompletionPercentage(params *DatabaseParameters) (string, error) {
	statement := fmt.Sprintf(`SET NOCOUNT ON

	%s

	INSERT INTO @s
	exec msdb.dbo.rds_task_status @db_name='%s'

	SELECT TOP 1 complete FROM @s

	SET NOCOUNT OFF`, statusTableDeclaration, params.DatabaseName)

	args := getCommandArgs(params, statement)
	output, err := execute(args)
	if err != nil {
		return "", err
	}
	return getSQLOutput(output), nil
}

// GetTaskMessage returns the message of the latest backup task
func GetTaskMessage(params *DatabaseParameters) (string, error) {
	statement := fmt.Sprintf(`SET NOCOUNT ON

	%s

	INSERT INTO @s
	exec msdb.dbo.rds_task_status @db_name='%s'

	SELECT TOP 1 task_info FROM @s

	SET NOCOUNT OFF`, statusTableDeclaration, params.DatabaseName)

	args := getCommandArgs(params, statement)
	output, err := execute(args)
	if err != nil {
		return "", err
	}
	return getSQLOutput(output), nil
}

// StartBackup creates a new backup
func StartBackup(params *BackupParameters) (string, error) {
	statement := fmt.Sprintf(`SET NOCOUNT ON

		%s

		INSERT INTO @s
		exec msdb.dbo.rds_backup_database
			@source_db_name='%s',
			@s3_arn_to_backup_to='arn:aws:s3:::%s/%s.bak',
			@overwrite_S3_backup_file=1;

		SELECT TOP 1 task_id FROM @s

		SET NOCOUNT OFF`,
		createTableDeclaration,
		params.DatabaseName,
		params.BucketName,
		params.Filename)

	args := getCommandArgs(&params.DatabaseParameters, statement)
	output, err := execute(args)
	if err != nil {
		return "", err
	}
	lines := strings.Split(output, "\n")
	if len(lines) < 4 {
		return "", errors.New(output)
	}
	return strings.TrimSpace(lines[3]), nil
}

func execute(args []string) (string, error) {
	byteOutput, err := exec.Command("docker", args...).Output()
	return string(byteOutput), err
}

func getSQLOutput(rawOutput string) string {
	lines := strings.Split(rawOutput, "\n")
	return strings.TrimSpace(lines[2])
}

func getCommandArgs(params *DatabaseParameters, statement string) []string {
	return []string{
		"exec",
		"-t",
		"mssql",
		"/opt/mssql-tools/bin/sqlcmd",
		"-S",
		params.Server,
		"-d",
		params.DatabaseName,
		"-U",
		params.Username,
		"-P",
		params.Password,
		"-Q",
		statement,
	}
}
