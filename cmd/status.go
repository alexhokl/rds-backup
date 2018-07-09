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
	"strings"

	"github.com/alexhokl/rds-backup/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {

	opts := statusOptions{}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show the status of the latest backup",
		Long:  "Show the status of the latest backup",
		Run: func(cmd *cobra.Command, args []string) {
			bindConfiguration(cmd)
			viper.Set("verbose", opts.verbose)
			if viper.GetBool("verbose") {
				dumpParameters(cmd)
			}
			errOpt := validateStatusOptions()
			if errOpt != nil {
				fmt.Println(errOpt.Error())
				cmd.HelpFunc()(cmd, args)
				return
			}
			cmdLine := &client.CommandLine{}
			err := runStatus(cmdLine)
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}

	flags := statusCmd.Flags()
	bindStatusOptions(flags, &opts)

	RootCmd.AddCommand(statusCmd)
}

func runStatus(cmdLine client.Command) error {
	params := &client.DatabaseParameters{
		Server:       viper.GetString("server"),
		Username:     viper.GetString("username"),
		Password:     viper.GetString("password"),
		DatabaseName: viper.GetString("database"),
	}

	c := client.GetClient(cmdLine)
	if c == nil {
		return errors.New("Unable to find a sqlcmd client")
	}

	output, err := c.GetStatus(cmdLine, params, "")
	if err != nil {
		return err
	}

	if output == "ERROR" {
		errorMessage, errErr := c.GetTaskMessage(cmdLine, params)
		if errErr != nil {
			return errErr
		}
		fmt.Println(errorMessage)
		return nil
	}

	fmt.Println(output)

	return nil
}

func validateStatusOptions() error {
	messages := strings.Builder{}

	if viper.GetString("server") == "" {
		messages.WriteString("--server AWS RDS SQL server must be specified\n")
	}
	if viper.GetString("username") == "" {
		messages.WriteString("--username AWS RDS SQL server login name must be specified\n")
	}
	if viper.GetString("password") == "" {
		messages.WriteString("--password AWS RDS SQL server login password must be specified\n")
	}
	if viper.GetString("database") == "" {
		messages.WriteString("--database Name of database must be specified\n")
	}

	if messages.String() != "" {
		return errors.New(messages.String())
	}

	return nil
}
