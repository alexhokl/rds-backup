package cmd

import (
	"fmt"

	"github.com/alexhokl/rds-backup/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type basicOptions struct {
	verbose      bool
	databaseName string
}

type nativeRestoreOptions struct {
	restoreDataDirectory               string
	restoreServerInstallationDirectory string
}

type dockerRestoreOptions struct {
	containerName string
	password      string
	port          int
}

type basicBackupOptions struct {
	filename string
}

type basicRestoreOptions struct {
	restoreDatabaseName string
	dataName            string
	logName             string
	isNative            bool
}

type basicDownloadOptions struct {
	bucketName string
}

type localDownloadOptions struct {
	downloadDirectory string
}

type serverOptions struct {
	server         string
	serverUsername string
	serverPassword string
}

type statusOptions struct {
	basicOptions
	serverOptions
}

type restoreOptions struct {
	basicOptions
	nativeRestoreOptions
	dockerRestoreOptions
	basicBackupOptions
	basicRestoreOptions
	localDownloadOptions
}

type downloadOptions struct {
	basicOptions
	basicRestoreOptions
	nativeRestoreOptions
	dockerRestoreOptions
	basicBackupOptions
	basicDownloadOptions
	localDownloadOptions
	isRestore bool
}

type createOptions struct {
	basicOptions
	nativeRestoreOptions
	dockerRestoreOptions
	basicBackupOptions
	serverOptions
	basicDownloadOptions
	localDownloadOptions
	isNative            bool
	isDownload          bool
	isWaitForCompletion bool
	isRestore           bool
}

func bindBasicOptions(flags *pflag.FlagSet, opts *basicOptions) {
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
}

func bindNativeRestoreOptions(flags *pflag.FlagSet, opts *nativeRestoreOptions) {
	flags.StringVar(&opts.restoreDataDirectory, "restore-data-directory", "", "Path to the directory where MDF and LDF files to be located")
	flags.StringVar(&opts.restoreServerInstallationDirectory, "restore-server-directory", client.DefaultServerInstallationPath, "Path to the directory the native SQL server")
}

func bindDockerRestoreOptions(flags *pflag.FlagSet, opts *dockerRestoreOptions) {
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVar(&opts.password, "restore-password", "", "Password of the MSSQL server in the container to be created")
	flags.IntVar(&opts.port, "port", client.DefaultServerPort, "port of restored server container")
}

func bindBasicBackupOptions(flags *pflag.FlagSet, opts *basicBackupOptions) {
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
}

func bindBasicRestoreOptions(flags *pflag.FlagSet, opts *basicRestoreOptions) {
	flags.StringVar(&opts.restoreDatabaseName, "restore-database", "", "Name of restored database")
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.StringVarP(&opts.dataName, "mdf", "m", "", "Logical name of data")
	flags.StringVarP(&opts.logName, "ldf", "l", "", "Logical name of log")
}

func bindBasicDownloadOptions(flags *pflag.FlagSet, opts *basicDownloadOptions) {
	flags.StringVarP(&opts.bucketName, "bucket", "b", "", "Bucket name")
}

func bindLocalDownloadOptions(flags *pflag.FlagSet, opts *localDownloadOptions) {
	flags.StringVar(&opts.downloadDirectory, "download-directory", "", "Path to the directory where backup from AWS S3 located")
}

func bindServerOptions(flags *pflag.FlagSet, opts *serverOptions) {
	flags.StringVarP(&opts.server, "server", "s", "", "Source SQL server")
	flags.StringVarP(&opts.serverUsername, "username", "u", "", "Source SQL server login name")
	flags.StringVarP(&opts.serverPassword, "password", "p", "", "Source SQL server login password")
}

func bindStatusOptions(flags *pflag.FlagSet, opts *statusOptions) {
	bindBasicOptions(flags, &opts.basicOptions)
	bindServerOptions(flags, &opts.serverOptions)
}

func bindRestoreOptions(flags *pflag.FlagSet, opts *restoreOptions) {
	bindBasicOptions(flags, &opts.basicOptions)
	bindNativeRestoreOptions(flags, &opts.nativeRestoreOptions)
	bindDockerRestoreOptions(flags, &opts.dockerRestoreOptions)
	bindBasicBackupOptions(flags, &opts.basicBackupOptions)
	bindBasicRestoreOptions(flags, &opts.basicRestoreOptions)
	bindLocalDownloadOptions(flags, &opts.localDownloadOptions)
}

func bindDownloadOptions(flags *pflag.FlagSet, opts *downloadOptions) {
	bindBasicOptions(flags, &opts.basicOptions)
	bindBasicRestoreOptions(flags, &opts.basicRestoreOptions)
	bindNativeRestoreOptions(flags, &opts.nativeRestoreOptions)
	bindDockerRestoreOptions(flags, &opts.dockerRestoreOptions)
	bindBasicBackupOptions(flags, &opts.basicBackupOptions)
	bindLocalDownloadOptions(flags, &opts.localDownloadOptions)
	flags.BoolVarP(&opts.isRestore, "restore", "r", false, "Restore backup in a docker container")
}

func bindCreateOptions(flags *pflag.FlagSet, opts *createOptions) {
	bindBasicOptions(flags, &opts.basicOptions)
	bindNativeRestoreOptions(flags, &opts.nativeRestoreOptions)
	bindDockerRestoreOptions(flags, &opts.dockerRestoreOptions)
	bindBasicBackupOptions(flags, &opts.basicBackupOptions)
	bindServerOptions(flags, &opts.serverOptions)
	bindBasicDownloadOptions(flags, &opts.basicDownloadOptions)
	bindLocalDownloadOptions(flags, &opts.localDownloadOptions)
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.BoolVarP(&opts.isWaitForCompletion, "wait", "w", false, "Wait for backup to complete")
	flags.BoolVar(&opts.isDownload, "download", false, "Create and download the backup")
	flags.BoolVarP(&opts.isRestore, "restore", "r", false, "Restore backup in a docker container")
}

func bindConfiguration(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})
}

func dumpParameters(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		fmt.Printf("%s: %s\n", f.Name, viper.GetString(f.Name))
	})
}
