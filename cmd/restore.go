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

func init() {
	opts := restoreOptions{}

	var restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restores the specified backup in a docker container",
		Long:  "Restores the specified backup in a docker container",
		Run: func(cmd *cobra.Command, args []string) {
			bindConfiguration(cmd)
			viper.Set("verbose", opts.verbose)
			if viper.GetBool("verbose") {
				dumpParameters(cmd)
			}
			errOpt := validateRestoreOptions()
			if errOpt != nil {
				fmt.Println(errOpt.Error())
				cmd.HelpFunc()(cmd, args)
				return
			}
			err := runRestore()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := restoreCmd.Flags()
	bindRestoreOptions(flags, &opts)

	RootCmd.AddCommand(restoreCmd)
}

func runRestore() error {
	if viper.GetBool("native") {
		errNative := client.RestoreNative(
			viper.GetString("filename"),
			viper.GetString("database"),
			viper.GetString("mdf"),
			viper.GetString("ldf"),
			viper.GetString("restore-database"),
			viper.GetString("restore-data-directory"),
		)
		if errNative != nil {
			return errNative
		}
		return nil
	}
	err := client.Restore(
		viper.GetString("filename"),
		viper.GetString("container"),
		viper.GetString("restore-password"),
		viper.GetString("database"),
		viper.GetString("mdf"),
		viper.GetString("ldf"),
		viper.GetInt("port"),
	)
	if err != nil {
		return err
	}
	return nil
}

func validateRestoreOptions() error {
	if viper.GetString("filename") == "" {
		return errors.New("--filename Filename must be specified")
	}
	if viper.GetBool("native") {
		if viper.GetInt("port") != client.DefaultServerPort {
			return errors.New("--port Port cannot be used in restoring to local native SQL server")
		}
	} else {
		if viper.GetString("container") == "" {
			return errors.New("--container Container name must be specified")
		}
		if viper.GetString("restore-password") == "" {
			return errors.New("restore-password Password of the restored SQL server must be specified")
		}
		if viper.GetString("restore-database") != "" {
			return errors.New("--restore-database cannot be used in Docker container restore")
		}
		if viper.GetString("restore-data-directory") != "" {
			return errors.New("--restore-data-directory cannot be used in Docker container restore")
		}
	}
	if viper.GetString("database") == "" {
		return errors.New("--database Name of database must be specified")
	}
	if viper.GetString("mdf") == "" {
		return errors.New("--mdf Logical name of data must be specified")
	}
	if viper.GetString("ldf") == "" {
		return errors.New("--ldf Logical name of log must be specified")
	}
	return nil
}
