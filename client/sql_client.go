package client

// SQLClient performs SQL operations
type SQLClient interface {
	IsEnvironmentSatisfied() bool
	GetStatus(*DatabaseParameters, string) (string, error)
	GetCompletionPercentage(*DatabaseParameters) (string, error)
	GetTaskMessage(*DatabaseParameters) (string, error)
	StartBackup(*BackupParameters) (string, error)
	GetLogicalNames(*DatabaseParameters) (string, string, error)
}

// GetClient returns a SQL client which can be run on this machine
func GetClient() SQLClient {
	nativeCli := &NativeClient{}
	if nativeCli.IsEnvironmentSatisfied() {
		return nativeCli
	}
	dockerCli := &DockerSQLClient{}
	if dockerCli.IsEnvironmentSatisfied() {
		return dockerCli
	}
	return nil
}
