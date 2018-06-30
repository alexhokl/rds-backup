// Copyright © 2017 Alex Ho <alexhokl@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"

	"github.com/alexhokl/rds-backup/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type downloadOptions struct {
	verbose              bool
	databaseName         string
	filename             string
	bucketName           string
	isRestore            bool
	containerName        string
	password             string
	dataName             string
	logName              string
	isNative             bool
	restoreDatabaseName  string
	restoreDataDirectory string
	port                 int
}

func init() {

	opts := downloadOptions{}

	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a backup from AWS S3 with option of restore",
		Long:  "Download a backup from AWS S3 with option of restore",
		Run: func(cmd *cobra.Command, args []string) {
			errOpt := validateDownloadOptions(opts)
			if errOpt != nil {
				fmt.Println(errOpt.Error())
				cmd.HelpFunc()(cmd, args)
				return
			}
			viper.Set("verbose", opts.verbose)
			err := runDownload(opts)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := downloadCmd.Flags()
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.BoolVarP(&opts.isRestore, "restore", "r", false, "Restore backup in a docker container")
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	viper.BindPFlag("database", downloadCmd.Flags().Lookup("database"))
	flags.StringVarP(&opts.bucketName, "bucket", "b", "", "Bucket name")
	viper.BindPFlag("bucket", downloadCmd.Flags().Lookup("bucket"))
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	viper.BindPFlag("filename", downloadCmd.Flags().Lookup("filename"))
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	viper.BindPFlag("container", downloadCmd.Flags().Lookup("container"))
	flags.StringVarP(&opts.password, "restore-password", "p", "", "Password of the restored server")
	viper.BindPFlag("restorePassword", downloadCmd.Flags().Lookup("restore-password"))
	flags.StringVarP(&opts.dataName, "mdf", "m", "", "Logical name of data")
	viper.BindPFlag("mdf", downloadCmd.Flags().Lookup("mdf"))
	flags.StringVarP(&opts.logName, "ldf", "l", "", "Logical name of log")
	viper.BindPFlag("ldf", downloadCmd.Flags().Lookup("ldf"))
	flags.StringVar(&opts.restoreDatabaseName, "restore-database", "", "Name of restored database")
	viper.BindPFlag("restoreDatabase", downloadCmd.Flags().Lookup("restore-database"))
	flags.StringVar(&opts.restoreDataDirectory, "restore-data-directory", "", "Path to the directory where MDF and LDF files to be located")
	viper.BindPFlag("restoreDataDirectory", downloadCmd.Flags().Lookup("restore-data-directory"))
	flags.IntVar(&opts.port, "port", 1433, "port of restored server container")
	viper.BindPFlag("port", downloadCmd.Flags().Lookup("port"))

	RootCmd.AddCommand(downloadCmd)
}

func runDownload(opts downloadOptions) error {
	if !client.IsAwsCliInstalled() {
		return errors.New("AWS CLI is required")
	}
	if !client.IsAwsCredentialsConfigured() {
		return errors.New("AWS CLI credentials are not configured yet. Please try 'aws configure'")
	}

	fmt.Println("Download started...")

	errDownload := client.DownloadBackup(viper.GetString("bucket"), viper.GetString("filename"))
	if errDownload != nil {
		return errDownload
	}
	fmt.Println("Download of the backup is completed")

	if opts.isRestore {
		fmt.Println("Starting to restore...")
		if opts.isNative {
			errNative := client.RestoreNative(
				viper.GetString("filename"),
				viper.GetString("database"),
				viper.GetString("mdf"),
				viper.GetString("ldf"),
				viper.GetString("restoreDatabase"),
				viper.GetString("restoreDataDirectory"),
			)
			if errNative != nil {
				return errNative
			}
			return nil
		}
		errRestore := client.Restore(
			viper.GetString("filename"),
			viper.GetString("container"),
			viper.GetString("restorePassword"),
			viper.GetString("database"),
			viper.GetString("mdf"),
			viper.GetString("ldf"),
			viper.GetInt("port"),
		)
		if errRestore != nil {
			return errRestore
		}
		fmt.Println("Restore completed.")
	}

	return nil
}

func validateDownloadOptions(opts downloadOptions) error {
	if viper.GetString("bucket") == "" {
		return errors.New("Bucket must be specified")
	}
	if viper.GetString("filename") == "" {
		return errors.New("Filename must be specified")
	}

	if opts.isRestore {
		if viper.GetString("database") == "" {
			return errors.New("Database must be specified")
		}
		if viper.GetString("mdf") == "" {
			return errors.New("Logical name of data must be specified")
		}
		if viper.GetString("ldf") == "" {
			return errors.New("Logical name of log must be specified")
		}
		if opts.isNative {
			if viper.GetString("port") != "" {
				return errors.New("Port cannot be used in restoring to local native SQL server")
			}
		} else {
			if viper.GetString("container") == "" {
				return errors.New("Container name must be specified")
			}
			if viper.GetString("restorePassword") == "" {
				return errors.New("Password must be specified")
			}
			if viper.GetString("restoreDatabase") != "" {
				return errors.New("restore-database cannot be used in Docker container restore")
			}
			if viper.GetString("restoreDataDirectory") != "" {
				return errors.New("restore-data-directory cannot be used in Docker container restore")
			}
		}
	}

	return nil
}
