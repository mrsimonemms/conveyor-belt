/*
Copyright © 2023 Simon Emms <simon@simonemms.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootOpts struct {
	CfgFile  string
	LogLevel string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	exec, err := os.Executable()
	cobra.CheckErr(err)

	rootCmd = &cobra.Command{
		Use:   filepath.Base(exec),
		Short: "Build your own pipelines",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set log level
			lvl, err := zerolog.ParseLevel(rootOpts.LogLevel)
			if err != nil {
				return err
			}
			zerolog.SetGlobalLevel(lvl)

			return nil
		},
	}

	dir, err := os.Getwd()
	cobra.CheckErr(err)

	viper.SetDefault("config", filepath.Join(dir, "config.yaml"))
	viper.SetDefault("log-level", zerolog.InfoLevel)
	rootCmd.PersistentFlags().StringVarP(&rootOpts.CfgFile, "config", "c", viper.GetString("config"), "path to config file")
	rootCmd.PersistentFlags().StringVar(&rootOpts.LogLevel, "log-level", viper.GetString("log-level"), "log level")
}
