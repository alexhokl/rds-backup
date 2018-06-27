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
	verbose        bool
	databaseName   string
	filename       string
	bucketName     string
	isRestore      bool
	containerName  string
	password       string
	dataName       string
	logName        string
	server         string
	serverUsername string
	serverPassword string
}

func init() {

	opts := downloadOptions{}

	var createCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a backup from AWS S3 with option of restore",
		Long:  "Download a backup from AWS S3 with option of restore",
		Run: func(cmd *cobra.Command, args []string) {
			opts = bindDownloadConfiguration(opts)
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

	flags := createCmd.Flags()
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.bucketName, "bucket", "b", "", "Bucket name")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.BoolVarP(&opts.isRestore, "restore", "r", false, "Restore backup in a docker container")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVarP(&opts.password, "password", "p", "", "Password of the MSSQL server in the container to be created")
	flags.StringVarP(&opts.dataName, "data", "m", "", "Logical name of data")
	flags.StringVarP(&opts.logName, "log", "l", "", "Logical name of log")
	flags.StringVarP(&opts.server, "server", "s", viper.GetString("server"), "Source SQL server")
	flags.StringVarP(&opts.serverUsername, "server-username", "n", "", "Source SQL server login name")
	flags.StringVarP(&opts.serverPassword, "server-password", "a", "", "Source SQL server login password")

	RootCmd.AddCommand(createCmd)
}

func runDownload(opts downloadOptions) error {
	if !client.IsAwsCliInstalled() {
		return errors.New("AWS CLI is required")
	}
	if !client.IsAwsCredentialsConfigured() {
		return errors.New("AWS CLI credentials are not configured yet. Please try 'aws configure'")
	}

	fmt.Println("Download started...")

	errDownload := client.DownloadBackup(opts.bucketName, opts.filename)
	if errDownload != nil {
		return errDownload
	}
	fmt.Println("Download of the backup is completed")

	if opts.isRestore {
		fmt.Println("Starting to restore...")
		errRestore := client.Restore(
			opts.filename,
			opts.containerName,
			opts.password,
			opts.databaseName,
			opts.dataName,
			opts.logName,
		)
		if errRestore != nil {
			return errRestore
		}
		fmt.Println("Restore completed.")
	}

	return nil
}

func bindDownloadConfiguration(opts downloadOptions) downloadOptions {
	if opts.server == "" {
		opts.server = viper.GetString("server")
	}
	if opts.serverUsername == "" {
		opts.serverUsername = viper.GetString("username")
	}
	if opts.serverPassword == "" {
		opts.serverPassword = viper.GetString("password")
	}
	if opts.databaseName == "" {
		opts.databaseName = viper.GetString("database")
	}
	if opts.containerName == "" {
		opts.containerName = viper.GetString("container")
	}
	if opts.password == "" {
		opts.password = viper.GetString("restorePassword")
	}
	if opts.dataName == "" {
		opts.dataName = viper.GetString("mdf")
	}
	if opts.logName == "" {
		opts.logName = viper.GetString("ldf")
	}
	if opts.bucketName == "" {
		opts.bucketName = viper.GetString("bucket")
	}
	return opts
}

func validateDownloadOptions(opts downloadOptions) error {
	if opts.server == "" {
		return errors.New("Source SQL server must be specified")
	}

	if opts.serverUsername == "" {
		return errors.New("Source SQL server login name must be specified")
	}

	if opts.serverPassword == "" {
		return errors.New("Source SQL server login password must be specified")
	}

	if opts.bucketName == "" {
		return errors.New("Bucket must be specified")
	}
	if opts.filename == "" {
		return errors.New("Filename must be specified")
	}

	if opts.isRestore {
		if opts.databaseName == "" {
			return errors.New("Database must be specified")
		}
		if opts.containerName == "" {
			return errors.New("Name of container must be specified")
		}
		if opts.password == "" {
			return errors.New("Password must be specified")
		}
		if opts.dataName == "" {
			return errors.New("Logical name of data must be specified")
		}
		if opts.logName == "" {
			return errors.New("Logical name of log must be specified")
		}
	}

	return nil
}
