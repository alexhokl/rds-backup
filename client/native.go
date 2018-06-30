package client

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// NativeClient is a SQL client runs sqlcmd on a machine
type NativeClient struct{}

// IsEnvironmentSatisfied returns if this client can be run on this machine
func (c *NativeClient) IsEnvironmentSatisfied() bool {
	if runtime.GOOS == "linux" {
		return false
	}
	args := []string{"-?"}
	_, err := executeSQLCmd(args)
	if err != nil {
		return false
	}
	return true
}

// GetStatus returns the status of the latest backup
func (c *NativeClient) GetStatus(params *DatabaseParameters, taskID string) (string, error) {
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

	args := getSQLCommandArgs(params, statement)
	output, err := executeSQLCmd(args)
	if err != nil {
		return "", err
	}
	return getSQLOutput(output), nil
}

// GetCompletionPercentage returns the percentage of completion of the latest backup
func (c *NativeClient) GetCompletionPercentage(params *DatabaseParameters) (string, error) {
	statement := fmt.Sprintf(`SET NOCOUNT ON

	%s

	INSERT INTO @s
	exec msdb.dbo.rds_task_status @db_name='%s'

	SELECT TOP 1 complete FROM @s

	SET NOCOUNT OFF`, statusTableDeclaration, params.DatabaseName)

	args := getSQLCommandArgs(params, statement)
	output, err := executeSQLCmd(args)
	if err != nil {
		return "", err
	}
	return getSQLOutput(output), nil
}

// GetTaskMessage returns the message of the latest backup task
func (c *NativeClient) GetTaskMessage(params *DatabaseParameters) (string, error) {
	statement := fmt.Sprintf(`SET NOCOUNT ON

	%s

	INSERT INTO @s
	exec msdb.dbo.rds_task_status @db_name='%s'

	SELECT TOP 1 task_info FROM @s

	SET NOCOUNT OFF`, statusTableDeclaration, params.DatabaseName)

	args := getSQLCommandArgs(params, statement)
	output, err := executeSQLCmd(args)
	if err != nil {
		return "", err
	}
	return getSQLOutput(output), nil
}

// StartBackup creates a new backup
func (c *NativeClient) StartBackup(params *BackupParameters) (string, error) {
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

	args := getSQLCommandArgs(&params.DatabaseParameters, statement)
	output, err := executeSQLCmd(args)
	if err != nil {
		return "", err
	}
	lines := strings.Split(output, "\n")
	if len(lines) < 4 {
		return "", errors.New(output)
	}
	return strings.TrimSpace(lines[3]), nil
}

// GetLogicalNames returns the logical names of MDF and LDF
func (c *NativeClient) GetLogicalNames(params *DatabaseParameters) (string, string, error) {
	dataNameQuery := "SELECT name FROM sys.master_files WHERE database_id = db_id() AND type = 0"
	logNameQuery := "SELECT name FROM sys.master_files WHERE database_id = db_id() AND type = 1"

	outputData, errData := executeSQLCmd(getSQLCommandArgs(params, dataNameQuery))
	if errData != nil {
		return "", "", errData
	}
	dataName := getSQLOutput(outputData)

	outputLog, errLog := executeSQLCmd(getSQLCommandArgs(params, logNameQuery))
	if errLog != nil {
		return "", "", errLog
	}
	logName := getSQLOutput(outputLog)

	return dataName, logName, nil
}

// RestoreNative restores a backup onto a local instance of SQL server
func RestoreNative(filename string, databaseName string, dataName string, logName string, renameDatabase string, customDataPath string) error {
	_, errFile := os.Stat(filename)
	if errFile != nil {
		return errFile
	}

	serverDirectory := "C:\\Program Files\\Microsoft SQL Server\\MSSQL13.MSSQLSERVER\\MSSQL\\"
	serverBackupDirectory := filepath.Join(serverDirectory, "Backup\\")
	serverMdfDirectory := filepath.Join(serverDirectory, "DATA\\")
	serverLdfDirectory := filepath.Join(serverDirectory, "LOG\\")

	currentDirectory, _ := os.Getwd()
	pathToBak := filepath.Join(currentDirectory, filename)
	pathToBackup := filepath.Join(serverBackupDirectory, filename)

	fmt.Println("Starting to restore onto local SQL Server...")

	copyFile(pathToBak, pathToBackup)

	fmt.Printf("Copied from %s to %s to prepare restoration.\n", pathToBak, pathToBackup)

	mdfDirectory := serverMdfDirectory
	ldfDirectory := serverLdfDirectory

	if customDataPath != "" {
		if _, errCustomPath := os.Stat(customDataPath); os.IsNotExist(errCustomPath) {
			return errCustomPath
		}
		mdfDirectory = customDataPath
		ldfDirectory = customDataPath
	}

	mdfPath := filepath.Join(mdfDirectory, fmt.Sprintf("%s.mdf", databaseName))
	ldfPath := filepath.Join(ldfDirectory, fmt.Sprintf("%s.ldf", databaseName))

	database := databaseName
	if renameDatabase != "" {
		database = renameDatabase
	}

	fmt.Println("Restoring...")

	restoreArgs := []string{
		"-Q",
		fmt.Sprintf("RESTORE DATABASE %s FROM DISK=N'%s' WITH FILE=1, NOUNLOAD, REPLACE, STATS=5, MOVE '%s' TO '%s', MOVE '%s' TO '%s'", database, pathToBackup, dataName, mdfPath, logName, ldfPath),
	}

	_, err := executeSQLCmd(restoreArgs)
	if err != nil {
		return err
	}
	fmt.Printf("Restore has been completed (as database '%s').\n", database)

	errRemove := os.Remove(pathToBackup)
	if errRemove != nil {
		return errRemove
	}
	fmt.Printf("Removed file %s. Clean up done.\n", pathToBackup)

	return nil
}

func executeSQLCmd(args []string) (string, error) {
	if viper.GetBool("verbose") {
		fmt.Println("Command executed:", "sqlcmd", args)
	}
	byteOutput, err := exec.Command("sqlcmd", args...).Output()
	return string(byteOutput), err
}

func getSQLCommandArgs(params *DatabaseParameters, statement string) []string {
	return []string{
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

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
