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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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
	bindDownloadOptions(flags, &opts)

	downloadCmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f)
	})

	RootCmd.AddCommand(downloadCmd)
}

func runDownload(opts downloadOptions) error {
	if !client.IsAwsCliInstalled() {
		return errors.New("AWS CLI is required")
	}
	if !client.IsAwsCredentialsConfigured() {
		return errors.New("AWS CLI credentials are not configured yet. Please try 'aws configure'")
	}

	errDownload := client.DownloadBackup(viper.GetString("bucket"), viper.GetString("filename"))
	if errDownload != nil {
		return errDownload
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
			viper.GetString("mdf"),
			viper.GetString("ldf"),
			viper.GetInt("port"),
		)
		if errRestore != nil {
			return errRestore
		}
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
