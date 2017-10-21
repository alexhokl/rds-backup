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
	"github.com/spf13/viper"
)

type createOptions struct {
	databaseName string
	filename     string
	bucketName   string
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
	output, err := client.StartBackup(params)
	if err != nil {
		return err
	}
	if output == "" {
		return errors.New("Unable to create a backup task")
	}
	fmt.Printf("Backup task [%s] started...", output)

	return nil
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

	return nil
}
