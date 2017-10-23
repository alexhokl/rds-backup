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

	"github.com/alexhokl/rds-backup/client"
	"github.com/spf13/cobra"
)

type restoreOptions struct {
	databaseName  string
	filename      string
	containerName string
	password      string
	dataName      string
	logName       string
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
			err := runRestore(opts)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := restoreCmd.Flags()
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.StringVarP(&opts.password, "password", "p", "", "Create and download the backup")
	flags.StringVarP(&opts.dataName, "data", "m", "", "Logical name of data")
	flags.StringVarP(&opts.logName, "log", "l", "", "Logical name of log")

	RootCmd.AddCommand(restoreCmd)
}

func runRestore(opts restoreOptions) error {
	err := client.Restore(
		opts.filename,
		opts.containerName,
		opts.password,
		opts.databaseName,
		opts.dataName,
		opts.logName,
	)
	if err != nil {
		return err
	}
	return nil
}

func validateRestoreOptions(opts restoreOptions) error {
	if opts.filename == "" {
		return errors.New("Filename must be specified")
	}
	if opts.containerName == "" {
		return errors.New("Container name must be specified")
	}
	if opts.password == "" {
		return errors.New("Password must be specified")
	}
	if opts.databaseName == "" {
		return errors.New("Database name must be specified")
	}
	if opts.dataName == "" {
		return errors.New("Logical name of data must be specified")
	}
	if opts.logName == "" {
		return errors.New("Logical name of log must be specified")
	}
	return nil
}
