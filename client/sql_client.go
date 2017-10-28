package client

type SqlClient interface {
	IsEnvironmentSatisfied() bool
	GetStatus(*DatabaseParameters, string) (string, error)
	GetCompletionPercentage(*DatabaseParameters) (string, error)
	GetTaskMessage(*DatabaseParameters) (string, error)
	StartBackup(*BackupParameters) (string, error)
	GetLogicalNames(*DatabaseParameters) (string, string, error)
}

func GetClient() SqlClient {
	nativeCli := &NativeClient{}
	if nativeCli.IsEnvironmentSatisfied() {
		return nativeCli
	}
	dockerCli := &DockerSqlClient{}
	if dockerCli.IsEnvironmentSatisfied() {
		return dockerCli
	}
	return nil
}
