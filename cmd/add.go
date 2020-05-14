// Copyright © 2020 Thilina Manamgoda
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
	"github.com/ThilinaManamgoda/etchosts/pkg/inputs"
	"github.com/ThilinaManamgoda/etchosts/pkg/parser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a host entry to /etc/hosts file",
	Long:  `Add a host entry to /etc/hosts file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.Parser{Path: viper.GetString(inputs.FlagFile)}
		err := p.Init()
		if err != nil {
			return errors.Wrap(err, "Unable to initialize parser")
		}
		ip, comment, domains, err := inputs.GetFlagsForAdd(cmd)
		if err != nil {
			return errors.Wrap(err, "Unable to get flags")
		}
		err = p.AddNewMapping(ip, domains, comment)
		if err != nil {
			return errors.Wrap(err, "Unable to add host mapping")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP(inputs.FlagIP, "i", "", "IP")
	addCmd.Flags().StringSliceP(inputs.FlagDomains, "d", nil, "Domain list")
	addCmd.Flags().StringP(inputs.FlagComment, "c", "", "Comment for the entry")

	inputs.MarkFlagRequired(addCmd, inputs.FlagIP)
	inputs.MarkFlagRequired(addCmd, inputs.FlagDomains)
}
