// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"time"

	"github.com/alexhokl/rds-backup/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type createOptions struct {
	databaseName        string
	filename            string
	bucketName          string
	isDownload          bool
	isWaitForCompletion bool
	isRestore           bool
	containerName       string
	password            string
}

func init() {

	opts := createOptions{}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Creates a new backup",
		Long:  "Creates a new backup",
		Run: func(cmd *cobra.Command, args []string) {
			err := runCreate(opts)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := createCmd.Flags()
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.bucketName, "bucket", "b", "", "Bucket name")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.BoolVar(&opts.isDownload, "download", false, "Create and download the backup")
	flags.BoolVarP(&opts.isWaitForCompletion, "wait", "w", false, "Wait for backup to complete")
	flags.BoolVarP(&opts.isRestore, "restore", "r", false, "Restore backup in a docker container")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVarP(&opts.password, "password", "p", "", "Password of the MSSQL server in the container to be created")

	RootCmd.AddCommand(createCmd)
}

func runCreate(opts createOptions) error {
	errOpt := validateCreateOptions(opts)
	if errOpt != nil {
		return errOpt
	}
	params := &client.BackupParameters{
		DatabaseParameters: client.DatabaseParameters{
			Server:       viper.GetString("server"),
			Username:     viper.GetString("username"),
			Password:     viper.GetString("password"),
			DatabaseName: opts.databaseName,
		},
		BucketName: opts.bucketName,
		Filename:   opts.filename,
	}

	c := client.GetClient()
	if c == nil {
		return errors.New("Unable to find a sqlcmd client")
	}

	dataLogicalName := ""
	logLogicalName := ""
	if opts.isRestore {
		dataName, logName, errLogicalNames := c.GetLogicalNames(&params.DatabaseParameters)
		if errLogicalNames != nil {
			return errLogicalNames
		}
		dataLogicalName = dataName
		logLogicalName = logName
	}

	taskID, err := c.StartBackup(params)
	if err != nil {
		return err
	}
	if taskID == "" {
		return errors.New("Unable to create a backup task")
	}
	fmt.Printf("Backup task [%s] started...", taskID)

	if opts.isDownload || opts.isWaitForCompletion || opts.isRestore {
		errBackup := isBackupCompleted(c, params, taskID)
		if errBackup != nil {
			return errBackup
		}
		fmt.Println("Backup completed.")
	}

	if opts.isDownload || opts.isRestore {
		errDownload := client.DownloadBackup(params.BucketName, params.Filename)
		if errDownload != nil {
			return errDownload
		}
		fmt.Println("Download of the backup is completed")
	}

	fmt.Println("")

	if opts.isRestore {
		fmt.Println("Starting to restore...")
		errRestore := client.Restore(
			opts.filename,
			opts.containerName,
			opts.password,
			opts.databaseName,
			dataLogicalName,
			logLogicalName,
		)
		if errRestore != nil {
			return errRestore
		}
		fmt.Println("Restore completed.")
	}

	return nil
}

func isBackupCompleted(c client.SqlClient, params *client.BackupParameters, taskID string) error {
	done := false
	var err error = nil

	for !done {
		fmt.Printf(".")
		time.Sleep(5 * time.Second)
		done, err = isBackupDone(c, &params.DatabaseParameters, taskID)
		if err != nil {
			return err
		}
	}

	fmt.Println("")
	if err != nil {
		return err
	}

	return nil
}

func isBackupDone(c client.SqlClient, params *client.DatabaseParameters, taskID string) (bool, error) {
	status, err := c.GetStatus(params, taskID)
	if err != nil {
		return false, err
	}
	if status == "ERROR" {
		errorMessage, errErr := c.GetTaskMessage(params)
		if errErr != nil {
			return false, errErr
		}
		fmt.Println(errorMessage)
		return false, errors.New(errorMessage)
	}
	return status == "SUCCESS", nil
}

func validateCreateOptions(opts createOptions) error {
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
	}

	return nil
}
