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
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var entryFileName string

// entryCmd represents the entry command
var entryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Generate entry file (main.go) for the project",
	Long:  `This command helps generating a main file for your web project.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := viper.Sub("entry").GetString("dir")
		if moduleRootPath != "" {
			dir = moduleRootPath
		}

		fileName := filepath.Join(viper.GetString("root"), dir, entryFileName)
		if err := checkFileName(fileName); err != nil {
			er(err)
		}

		if err := createEntryFile(fileName); err != nil {
			er(err)
		}

		fmt.Printf("Entry file created: %s\n", fileName)
	},
}

func init() {
	makeCmd.AddCommand(entryCmd)

	entryCmd.PersistentFlags().StringVar(&entryFileName, "name", "main.go", "Entry file name")
	viper.SetDefault("entry.dir", "")
}

func createEntryFile(fileName string) error {
	template := `package main

import (
	"github.com/fpay/gopress"
)

func main() {
	// create server
	s := gopress.NewServer(gopress.ServerOptions{
		Port: 3000,
	})

	// init and register services
	// s.RegisterServices(
	// 	services.NewDatabaseService(),
	// )

	// register middlewares
	s.RegisterGlobalMiddlewares(
		gopress.NewLoggingMiddleware("global"),
	)

	// init and register controllers
	// s.RegisterControllers(
	// 	controllers.NewUsersController(),
	// )

	s.Start()
}
`

	entryScript, err := render(template, nil)
	if err != nil {
		er(err)
	}

	err = writeStringToFile(fileName, entryScript)
	if err != nil {
		er(err)
	}

	return nil
}
