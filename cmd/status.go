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

type statusOptions struct {
	verbose        bool
	databaseName   string
	server         string
	serverUsername string
	serverPassword string
}

func init() {

	opts := statusOptions{}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show the status of the latest backup",
		Long:  "Show the status of the latest backup",
		Run: func(cmd *cobra.Command, args []string) {
			opts = bindStatusConfiguration(opts)
			errOpt := validateStatusOptions(opts)
			if errOpt != nil {
				fmt.Println(errOpt.Error())
				cmd.HelpFunc()(cmd, args)
				return
			}
			viper.Set("verbose", opts.verbose)
			err := runStatus(opts)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := statusCmd.Flags()
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	flags.StringVarP(&opts.databaseName, "database", "d", "", "Name of database")
	flags.StringVarP(&opts.server, "server", "s", "", "Source SQL server")
	flags.StringVarP(&opts.serverUsername, "server-username", "n", "", "Source SQL server login name")
	flags.StringVarP(&opts.serverPassword, "server-password", "a", "", "Source SQL server login password")

	RootCmd.AddCommand(statusCmd)
}

func runStatus(opts statusOptions) error {
	params := &client.DatabaseParameters{
		Server:       viper.GetString("server"),
		Username:     viper.GetString("username"),
		Password:     viper.GetString("password"),
		DatabaseName: opts.databaseName,
	}

	c := client.GetClient()
	if c == nil {
		return errors.New("Unable to find a sqlcmd client")
	}

	output, err := c.GetStatus(params, "")
	if err != nil {
		return err
	}

	if output == "ERROR" {
		errorMessage, errErr := c.GetTaskMessage(params)
		if errErr != nil {
			return errErr
		}
		fmt.Println(errorMessage)
		return nil
	}

	fmt.Println(output)

	return nil
}

func bindStatusConfiguration(opts statusOptions) statusOptions {
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
	return opts
}

func validateStatusOptions(opts statusOptions) error {
	if opts.server == "" {
		return errors.New("Source SQL server must be specified")
	}

	if opts.serverUsername == "" {
		return errors.New("Source SQL server login name must be specified")
	}

	if opts.serverPassword == "" {
		return errors.New("Source SQL server login password must be specified")
	}

	if opts.databaseName == "" {
		return errors.New("Database must be specified")
	}

	return nil
}
