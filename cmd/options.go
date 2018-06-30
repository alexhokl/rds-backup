package cmd

import (
	"github.com/spf13/pflag"
)

type basicOptions struct {
	verbose bool
}

type databaseOptions struct {
	databaseName string
}

type backupFileOptions struct {
	filename string
}

type serverOptions struct {
	server         string
	serverUsername string
	serverPassword string
}

type statusOptions struct {
	basicOptions
	serverOptions
	databaseOptions
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
	databaseOptions
	backupFileOptions
	dataName string
	logName  string
	isNative bool
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
	backupFileOptions
	isDownload          bool
	isWaitForCompletion bool
	isNative            bool
}

func bindStatusOptions(flags *pflag.FlagSet, opts statusOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.StringVarP(&opts.server, "server", "s", "", "Source SQL server")
	flags.StringVarP(&opts.serverUsername, "username", "u", "", "Source SQL server login name")
	flags.StringVarP(&opts.serverPassword, "password", "p", "", "Source SQL server login password")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
}

func bindRestoreOptions(flags *pflag.FlagSet, opts restoreOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.StringVarP(&opts.dataName, "mdf", "m", "", "Logical name of data")
	flags.StringVarP(&opts.logName, "ldf", "l", "", "Logical name of log")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVarP(&opts.password, "restore-password", "a", "", "Password of the MSSQL server in the container to be created")
	flags.StringVar(&opts.restoreDatabaseName, "restore-database", "", "Name of restored database")
	flags.StringVar(&opts.restoreDataDirectory, "restore-data-directory", "", "Path to the directory where MDF and LDF files to be located")
	flags.IntVar(&opts.port, "port", 1433, "port of restored server container")
}

func bindDownloadOptions(flags *pflag.FlagSet, opts downloadOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.StringVarP(&opts.bucketName, "bucket", "b", "", "Bucket name")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVarP(&opts.password, "restore-password", "a", "", "Password of the MSSQL server in the container to be created")
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

func bindCreateOptions(flags *pflag.FlagSet, opts createOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.BoolVarP(&opts.isWaitForCompletion, "wait", "w", false, "Wait for backup to complete")
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.BoolVar(&opts.isDownload, "download", false, "Create and download the backup")
	flags.BoolVarP(&opts.isRestore, "restore", "r", false, "Restore backup in a docker container")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.bucketName, "bucket", "b", "", "Bucket name")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVarP(&opts.password, "restore-password", "a", "", "Password of the MSSQL server in the container to be created")
	flags.StringVarP(&opts.server, "server", "s", "", "Source SQL server")
	flags.StringVarP(&opts.serverUsername, "username", "u", "", "Source SQL server login name")
	flags.StringVarP(&opts.serverPassword, "password", "p", "", "Source SQL server login password")
	flags.StringVar(&opts.restoreDatabaseName, "restore-database", "", "Name of restored database")
	flags.StringVar(&opts.restoreDataDirectory, "restore-data-directory", "", "Path to the directory where MDF and LDF files to be located")
	flags.IntVar(&opts.port, "port", 1433, "port of restored server container")
}
