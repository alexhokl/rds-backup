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
	verbose       bool
	databaseName  string
	filename      string
	bucketName    string
	isRestore     bool
	containerName string
	password      string
	dataName      string
	logName       string
}

func init() {

	opts := downloadOptions{}

	var createCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a backup from AWS S3 with option of restore",
		Long:  "Download a backup from AWS S3 with option of restore",
		Run: func(cmd *cobra.Command, args []string) {
			errConfig := validateConfig()
			if errConfig != nil {
				fmt.Println(errConfig.Error())
				return
			}
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

	RootCmd.AddCommand(createCmd)
}

func runDownload(opts downloadOptions) error {
	c := client.GetClient()
	if c == nil {
		return errors.New("Unable to find a sqlcmd client")
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

func validateDownloadOptions(opts downloadOptions) error {
	if opts.databaseName == "" {
		return errors.New("Database must be specified")
	}
	if opts.bucketName == "" {
		return errors.New("Bucket must be specified")
	}
	if opts.filename == "" {
		return errors.New("Filename must be specified")
	}

	if opts.isRestore {
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
