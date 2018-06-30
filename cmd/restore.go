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

type restoreOptions struct {
	verbose              bool
	databaseName         string
	filename             string
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
	opts := restoreOptions{}

	var restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restores the specified backup in a docker container",
		Long:  "Restores the specified backup in a docker container",
		Run: func(cmd *cobra.Command, args []string) {
			errOpt := validateRestoreOptions(opts)
			if errOpt != nil {
				fmt.Println(errOpt.Error())
				cmd.HelpFunc()(cmd, args)
				return
			}
			viper.Set("verbose", opts.verbose)
			err := runRestore(opts)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := restoreCmd.Flags()
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.BoolVarP(&opts.isNative, "native", "n", false, "Restore to local native SQL server")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	viper.BindPFlag("database", flags.Lookup("database"))
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	viper.BindPFlag("container", flags.Lookup("container"))
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	viper.BindPFlag("filename", flags.Lookup("filename"))
	flags.StringVarP(&opts.password, "restore-password", "p", "", "Password of the restored server")
	viper.BindPFlag("restorePassword", flags.Lookup("restore-password"))
	flags.StringVarP(&opts.dataName, "mdf", "m", "", "Logical name of data")
	viper.BindPFlag("mdf", flags.Lookup("mdf"))
	flags.StringVarP(&opts.logName, "ldf", "l", "", "Logical name of log")
	viper.BindPFlag("ldf", flags.Lookup("ldf"))
	flags.StringVar(&opts.restoreDatabaseName, "restore-database", "", "Name of restored database")
	viper.BindPFlag("restoreDatabase", flags.Lookup("restore-database"))
	flags.StringVar(&opts.restoreDataDirectory, "restore-data-directory", "", "Path to the directory where MDF and LDF files to be located")
	viper.BindPFlag("restoreDataDirectory", flags.Lookup("restore-data-directory"))
	flags.IntVar(&opts.port, "port", 1433, "port of restored server container")
	viper.BindPFlag("port", flags.Lookup("port"))

	RootCmd.AddCommand(restoreCmd)
}

func runRestore(opts restoreOptions) error {
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
	err := client.Restore(
		viper.GetString("filename"),
		viper.GetString("container"),
		viper.GetString("restorePassword"),
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

func validateRestoreOptions(opts restoreOptions) error {
	if viper.GetString("filename") == "" {
		return errors.New("Filename must be specified")
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
	if viper.GetString("database") == "" {
		return errors.New("Database name must be specified")
	}
	if viper.GetString("mdf") == "" {
		return errors.New("Logical name of data must be specified")
	}
	if viper.GetString("ldf") == "" {
		return errors.New("Logical name of log must be specified")
	}
	return nil
}
