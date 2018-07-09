package client

// SQLClient performs SQL operations
type SQLClient interface {
	IsEnvironmentSatisfied(Command) bool
	GetStatus(Command, *DatabaseParameters, string) (string, error)
	GetCompletionPercentage(Command, *DatabaseParameters) (string, error)
	GetTaskMessage(Command, *DatabaseParameters) (string, error)
	StartBackup(Command, *BackupParameters) (string, error)
	GetLogicalNames(Command, *DatabaseParameters) (string, string, error)
}

// GetClient returns a SQL client which can be run on this machine
func GetClient(cmdLine Command) SQLClient {
	nativeCli := &NativeClient{}
	if nativeCli.IsEnvironmentSatisfied(cmdLine) {
		return nativeCli
	}
	dockerCli := &DockerSQLClient{}
	if dockerCli.IsEnvironmentSatisfied(cmdLine) {
		return dockerCli
	}
	return nil
}
