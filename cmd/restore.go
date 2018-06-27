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
	verbose        bool
	databaseName   string
	filename       string
	containerName  string
	password       string
	dataName       string
	logName        string
	server         string
	serverUsername string
	serverPassword string
}

func init() {
	opts := restoreOptions{}

	var restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restores the specified backup in a docker container",
		Long:  "Restores the specified backup in a docker container",
		Run: func(cmd *cobra.Command, args []string) {
			opts = bindRestoreConfiguration(opts)
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
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.containerName, "container", "c", "", "Name of container to be created")
	flags.StringVarP(&opts.filename, "filename", "f", "", "File name of the backup")
	flags.StringVarP(&opts.password, "password", "p", "", "Create and download the backup")
	flags.StringVarP(&opts.dataName, "data", "m", "", "Logical name of data")
	flags.StringVarP(&opts.logName, "log", "l", "", "Logical name of log")
	flags.StringVarP(&opts.server, "server", "s", "", "Source SQL server")
	flags.StringVarP(&opts.serverUsername, "server-username", "n", "", "Source SQL server login name")
	flags.StringVarP(&opts.serverPassword, "server-password", "a", "", "Source SQL server login password")

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

func bindRestoreConfiguration(opts restoreOptions) restoreOptions {
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
	return opts
}

func validateRestoreOptions(opts restoreOptions) error {
	if opts.server == "" {
		return errors.New("Source SQL server must be specified")
	}

	if opts.serverUsername == "" {
		return errors.New("Source SQL server login name must be specified")
	}

	if opts.serverPassword == "" {
		return errors.New("Source SQL server login password must be specified")
	}

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
