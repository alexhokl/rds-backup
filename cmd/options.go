package cmd

import (
	"github.com/spf13/pflag"
)

type statusOptions struct {
	verbose        bool
	server         string
	serverUsername string
	serverPassword string
	databaseName   string
}

type restoreOptions struct {
	verbose              bool
	databaseName         string
	filename             string
	dataName             string
	logName              string
	isNative             bool
	containerName        string
	password             string
	restoreDatabaseName  string
	restoreDataDirectory string
	port                 int
}

type downloadOptions struct {
	verbose              bool
	bucketName           string
	containerName        string
	password             string
	restoreDatabaseName  string
	restoreDataDirectory string
	port                 int
	databaseName         string
	filename             string
	dataName             string
	logName              string
	isNative             bool
	isRestore            bool
}

type createOptions struct {
	verbose              bool
	bucketName           string
	server               string
	serverUsername       string
	serverPassword       string
	databaseName         string
	containerName        string
	password             string
	restoreDatabaseName  string
	restoreDataDirectory string
	port                 int
	filename             string
	isDownload           bool
	isWaitForCompletion  bool
	isNative             bool
	isRestore            bool
}

func bindStatusOptions(flags *pflag.FlagSet, opts *statusOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.StringVarP(&opts.server, "server", "s", "", "Source SQL server")
	flags.StringVarP(&opts.serverUsername, "username", "u", "", "Source SQL server login name")
	flags.StringVarP(&opts.serverPassword, "password", "p", "", "Source SQL server login password")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
}

func bindRestoreOptions(flags *pflag.FlagSet, opts *restoreOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.StringVarP(&opts.dataName, "mdf", "m", "", "Logical name of data")
	flags.StringVarP(&opts.logName, "ldf", "l", "", "Logical name of log")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVar(&opts.password, "restore-password", "", "Password of the MSSQL server in the container to be created")
	flags.StringVar(&opts.restoreDatabaseName, "restoreDatabase", "", "Name of restored database")
	flags.StringVar(&opts.restoreDataDirectory, "restore-data-directory", "", "Path to the directory where MDF and LDF files to be located")
	flags.IntVar(&opts.port, "port", 1433, "port of restored server container")
}

func bindDownloadOptions(flags *pflag.FlagSet, opts *downloadOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.StringVarP(&opts.bucketName, "bucket", "b", "", "Bucket name")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVar(&opts.password, "restore-password", "", "Password of the MSSQL server in the container to be created")
	flags.StringVar(&opts.restoreDatabaseName, "restore-database", "", "Name of restored database")
	flags.StringVar(&opts.restoreDataDirectory, "restore-data-directory", "", "Path to the directory where MDF and LDF files to be located")
	flags.IntVar(&opts.port, "port", 1433, "port of restored server container")
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.StringVarP(&opts.dataName, "mdf", "m", "", "Logical name of data")
	flags.StringVarP(&opts.logName, "ldf", "l", "", "Logical name of log")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.BoolVarP(&opts.isRestore, "restore", "r", false, "Restore backup in a docker container")
}

func bindCreateOptions(flags *pflag.FlagSet, opts *createOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.BoolVarP(&opts.isWaitForCompletion, "wait", "w", false, "Wait for backup to complete")
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.BoolVar(&opts.isDownload, "download", false, "Create and download the backup")
	flags.BoolVarP(&opts.isRestore, "restore", "r", false, "Restore backup in a docker container")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.bucketName, "bucket", "b", "", "Bucket name")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVar(&opts.password, "restore-password", "", "Password of the MSSQL server in the container to be created")
	flags.StringVarP(&opts.server, "server", "s", "", "Source SQL server")
	flags.StringVarP(&opts.serverUsername, "username", "u", "", "Source SQL server login name")
	flags.StringVarP(&opts.serverPassword, "password", "p", "", "Source SQL server login password")
	flags.StringVar(&opts.restoreDatabaseName, "restore-database", "", "Name of restored database")
	flags.StringVar(&opts.restoreDataDirectory, "restore-data-directory", "", "Path to the directory where MDF and LDF files to be located")
	flags.IntVar(&opts.port, "port", 1433, "port of restored server container")
}
