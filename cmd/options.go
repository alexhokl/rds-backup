package cmd

type basicOptions struct {
	verbose bool
}

type serverOptions struct {
	databaseName   string
	server         string
	serverUsername string
	serverPassword string
}

type statusOptions struct {
	basicOptions
	serverOptions
}

type containerOptions struct {
	containerName string
	password      string
}

type nativeOptions struct {
	restoreDatabaseName  string
	restoreDataDirectory string
	port                 int
}

type basicRestoreOptions struct {
	databaseName string
	filename     string
	dataName     string
	logName      string
	isNative     bool
}

type restoreOptions struct {
	basicOptions
	basicRestoreOptions
	containerOptions
	nativeOptions
}

type s3Options struct {
	bucketName string
}

type downloadOptions struct {
	basicOptions
	s3Options
	containerOptions
	nativeOptions
	basicRestoreOptions
	isRestore bool
}

type createOptions struct {
	basicOptions
	s3Options
	serverOptions
	downloadOptions
	containerOptions
	nativeOptions
	filename            string
	isDownload          bool
	isWaitForCompletion bool
	isNative            bool
}
