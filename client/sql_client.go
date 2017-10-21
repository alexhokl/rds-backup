package client

type SqlClient interface {
	IsEnvironmentSatisfied() bool
	GetStatus(*DatabaseParameters, string) (string, error)
	GetCompletionPercentage(*DatabaseParameters) (string, error)
	GetTaskMessage(*DatabaseParameters) (string, error)
	StartBackup(*BackupParameters) (string, error)
}
