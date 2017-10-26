package client

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
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

type DockerSqlClient struct{}

func (c DockerSqlClient) IsEnvironmentSatisfied() bool {
	if !isDockerInstalled() {
		return false
	}
	if !isDockerContentTrustDisabled() {
		fmt.Println("Docker Content Trust is not disabled yet. Please run 'export DOCKER_CONTENT_TRUST=0'")
		return false
	}

	serverContainerName := getSqlServerContainerName()
	if serverContainerName == "" {
		return false
	}

	return true
}

// GetStatus returns the status of the latest backup
func (c DockerSqlClient) GetStatus(params *DatabaseParameters, taskID string) (string, error) {
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
func (c DockerSqlClient) GetCompletionPercentage(params *DatabaseParameters) (string, error) {
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
func (c DockerSqlClient) GetTaskMessage(params *DatabaseParameters) (string, error) {
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
func (c DockerSqlClient) StartBackup(params *BackupParameters) (string, error) {
	statement := fmt.Sprintf(`SET NOCOUNT ON

		%s

		INSERT INTO @s
		exec msdb.dbo.rds_backup_database
			@source_db_name='%s',
			@s3_arn_to_backup_to='arn:aws:s3:::%s/%s',
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

func Restore(filename string, containerName string, password string, databaseName string, dataName string, logName string) error {
	_, errFile := os.Stat(filename)
	if errFile != nil {
		return errFile
	}

	currentDirectory, _ := os.Getwd()
	directoryToMount := filepath.Dir(filepath.Join(currentDirectory, filename))

	createArgs := []string{
		"run",
		"--name",
		containerName,
		"-p",
		"1433:1433",
		"-v",
		fmt.Sprintf("%s/:/var/backups/", directoryToMount),
		"-e",
		fmt.Sprintf("SA_PASSWORD=%s", password),
		"-e",
		"ACCEPT_EULA=Y",
		"-d",
		"microsoft/mssql-server-linux",
	}

	_, errCreate := execute(createArgs)
	if errCreate != nil {
		return errCreate
	}

	fmt.Printf("MSSQL container %s is created. Waiting for SQL server to complete initialisation...\n", containerName)

	time.Sleep(90 * time.Second)

	fmt.Println("Restoring...")

	restoreArgs := []string{
		"exec",
		"-t",
		containerName,
		"/opt/mssql-tools/bin/sqlcmd",
		"-S",
		".",
		"-U",
		"sa",
		"-P",
		password,
		"-Q",
		fmt.Sprintf("RESTORE DATABASE %s FROM DISK=N'/var/backups/%s' WITH FILE=1, NOUNLOAD, REPLACE, STATS=5, MOVE '%s' TO '/var/opt/mssql/data/%s.mdf', MOVE '%s' TO '/var/opt/mssql/data/%s.ldf'", databaseName, filename, dataName, databaseName, logName, databaseName),
	}

	_, err := execute(restoreArgs)
	if err != nil {
		return err
	}
	fmt.Println("Restore is completed.")
	return nil
}

func (c DockerSqlClient) GetLogicalNames(params *DatabaseParameters) (string, string, error) {
	dataNameQuery := "SELECT name FROM sys.master_files WHERE database_id = db_id() AND type = 0"
	logNameQuery := "SELECT name FROM sys.master_files WHERE database_id = db_id() AND type = 1"

	outputData, errData := execute(getCommandArgs(params, dataNameQuery))
	if errData != nil {
		return "", "", errData
	}
	dataName := getSQLOutput(outputData)

	outputLog, errLog := execute(getCommandArgs(params, logNameQuery))
	if errLog != nil {
		return "", "", errLog
	}
	logName := getSQLOutput(outputLog)

	return dataName, logName, nil
}

func execute(args []string) (string, error) {
	if viper.GetBool("verbose") {
		fmt.Println("Command executed:", "docker", args)
	}
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

func isDockerInstalled() bool {
	_, err := execute([]string{"help"})
	return err == nil
}

func isDockerContentTrustDisabled() bool {
	return os.Getenv("DOCKER_CONTENT_TRUST") != "1"
}

func getSqlServerContainerName() string {
	args := []string{
		"ps",
		"-f",
		"ancestor=microsoft/mssql-server-linux",
		"-f",
		"status=running",
		"--format",
		"{{.Names}}",
	}
	output, err := execute(args)
	if err != nil {
		return ""
	}
	return strings.Split(output, "\n")[0]
}
