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
	"time"

	"github.com/spf13/pflag"

	"github.com/alexhokl/rds-backup/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {

	opts := createOptions{}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Creates a new backup",
		Long:  "Creates a new backup",
		Run: func(cmd *cobra.Command, args []string) {
			errOpt := validateCreateOptions(opts)
			if errOpt != nil {
				fmt.Println(errOpt.Error())
				cmd.HelpFunc()(cmd, args)
				return
			}
			viper.Set("verbose", opts.verbose)
			err := runCreate(opts)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := createCmd.Flags()
	bindCreateOptions(flags, &opts)

	createCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	RootCmd.AddCommand(createCmd)
}

func runCreate(opts createOptions) error {
	if opts.isDownload || opts.isRestore {
		if !client.IsAwsCliInstalled() {
			return errors.New("AWS CLI is required")
		}
		if !client.IsAwsCredentialsConfigured() {
			return errors.New("AWS CLI credentials are not configured yet. Please try 'aws configure'")
		}
	}

	params := &client.BackupParameters{
		DatabaseParameters: client.DatabaseParameters{
			Server:       viper.GetString("server"),
			Username:     viper.GetString("username"),
			Password:     viper.GetString("password"),
			DatabaseName: viper.GetString("database"),
		},
		BucketName: viper.GetString("bucket"),
		Filename:   viper.GetString("filename"),
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
		fmt.Printf("Backup completed (on AWS S3 at s3://%s/%s).\n", params.BucketName, params.Filename)
	}

	if opts.isDownload || opts.isRestore {
		errDownload := client.DownloadBackup(params.BucketName, params.Filename)
		if errDownload != nil {
			return errDownload
		}
	}

	if opts.isRestore {
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
			dataLogicalName,
			logLogicalName,
			viper.GetInt("port"),
		)
		if errRestore != nil {
			return errRestore
		}
	}

	return nil
}

func isBackupCompleted(c client.SQLClient, params *client.BackupParameters, taskID string) error {
	done := false
	var err error

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

func isBackupDone(c client.SQLClient, params *client.DatabaseParameters, taskID string) (bool, error) {
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
	if viper.GetString("server") == "" {
		return errors.New("Source SQL server must be specified")
	}

	if viper.GetString("username") == "" {
		return errors.New("Source SQL server login name must be specified")
	}

	if viper.GetString("password") == "" {
		return errors.New("Source SQL server login password must be specified")
	}

	if viper.GetString("database") == "" {
		return errors.New("Database must be specified")
	}
	if viper.GetString("bucket") == "" {
		return errors.New("Bucket must be specified")
	}
	if viper.GetString("filename") == "" {
		return errors.New("Filename must be specified")
	}

	if opts.isRestore {
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
