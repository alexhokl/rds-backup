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
	"os"
	"strings"

	"github.com/alexhokl/rds-backup/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {

	opts := downloadOptions{}

	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download a backup from AWS S3 with option of restore",
		Long:  "Download a backup from AWS S3 with option of restore",
		Run: func(cmd *cobra.Command, args []string) {
			bindConfiguration(cmd)
			viper.Set("verbose", opts.verbose)
			if viper.GetBool("verbose") {
				dumpParameters(cmd)
			}
			errOpt := validateDownloadOptions()
			if errOpt != nil {
				fmt.Println(errOpt.Error())
				cmd.HelpFunc()(cmd, args)
				return
			}
			err := runDownload()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := downloadCmd.Flags()
	bindDownloadOptions(flags, &opts)

	RootCmd.AddCommand(downloadCmd)
}

func runDownload() error {
	if !client.IsAwsCliInstalled() {
		return errors.New("AWS CLI is required")
	}
	if !client.IsAwsCredentialsConfigured() {
		return errors.New("AWS CLI credentials are not configured yet. Please try 'aws configure'")
	}

	errDownload := client.DownloadBackup(viper.GetString("bucket"), viper.GetString("filename"), viper.GetString("download-directory"))
	if errDownload != nil {
		return errDownload
	}

	basicRestoreParameters := client.BaseRestoreParameters{
		Filename:          viper.GetString("filename"),
		DatabaseName:      viper.GetString("database"),
		DataName:          viper.GetString("mdf"),
		LogName:           viper.GetString("ldf"),
		DownloadDirectory: viper.GetString("download-directory"),
	}

	if viper.GetBool("restore") {
		if viper.GetBool("native") {
			nativeParameters := &client.NativeRestoreParameters{
				BaseRestoreParameters: basicRestoreParameters,
				CustomDataPath:        viper.GetString("restore-data-directory"),
				ServerPath:            viper.GetString("restore-server-directory"),
			}
			errNative := client.RestoreNative(nativeParameters)
			if errNative != nil {
				return errNative
			}
			return nil
		}
		restoreParameters := &client.RestoreParameters{
			BaseRestoreParameters: basicRestoreParameters,
			ContainerName:         viper.GetString("container"),
			Password:              viper.GetString("restore-password"),
			Port:                  viper.GetInt("port"),
		}
		errRestore := client.Restore(restoreParameters)
		if errRestore != nil {
			return errRestore
		}
	}

	return nil
}

func validateDownloadOptions() error {
	messages := strings.Builder{}

	if viper.GetString("bucket") == "" {
		messages.WriteString("--bucket AWS S3 Bucket must be specified\n")
	}
	if viper.GetString("filename") == "" {
		messages.WriteString("--filename Filename must be specified\n")
	}

	downloadDirectory := viper.GetString("download-directory")
	if downloadDirectory != "" {
		if _, errDownloadDirectory := os.Stat(downloadDirectory); os.IsNotExist(errDownloadDirectory) {
			messages.WriteString(fmt.Sprintf("the specified download-directory (%s) does not exist\n", downloadDirectory))
		}
	}

	if viper.GetBool("restore") {
		if viper.GetString("database") == "" {
			messages.WriteString("--database Name of database must be specified\n")
		}
		if viper.GetString("mdf") == "" {
			messages.WriteString("--mdf Logical name of data must be specified\n")
		}
		if viper.GetString("ldf") == "" {
			messages.WriteString("--ldf Logical name of log must be specified\n")
		}
		if viper.GetBool("native") {
			if viper.GetInt("port") != client.DefaultServerPort {
				messages.WriteString("--port Port cannot be used in restoring to local native SQL server\n")
			}
			restoreServerDirectory := viper.GetString("restore-server-directory")
			if restoreServerDirectory != "" {
				if _, errServerDirectory := os.Stat(restoreServerDirectory); os.IsNotExist(errServerDirectory) {
					messages.WriteString("the specified restore-server-directory does not exist\n")
				}
			}
		} else {
			if viper.GetString("container") == "" {
				messages.WriteString("--container Container name must be specified\n")
			}
			if viper.GetString("restore-password") == "" {
				messages.WriteString("restore-password Password of the restored SQL server must be specified\n")
			}
			if viper.GetString("restore-database") != "" {
				messages.WriteString("--restore-database cannot be used in Docker container restore\n")
			}
			if viper.GetString("restore-data-directory") != "" {
				messages.WriteString("--restore-data-directory cannot be used in Docker container restore\n")
			}
		}
	}

	if messages.String() != "" {
		return errors.New(messages.String())
	}

	return nil
}
