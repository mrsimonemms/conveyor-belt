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
	"github.com/mrsimonemms/conveyor-belt/pkg/config"
	"github.com/mrsimonemms/conveyor-belt/pkg/pipeline"
	"github.com/mrsimonemms/conveyor-belt/pkg/server"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(rootOpts.CfgFile)
		if err != nil {
			return err
		}

		s := server.New(cfg)

		p, err := pipeline.Build(cfg)
		if err != nil {
			return err
		}

		if err := s.Triggers(p); err != nil {
			return err
		}

		return s.Start()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
