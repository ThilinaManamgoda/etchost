// Copyright © 2020 NAME HERE <EMAIL ADDRESS>
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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [ID]",
	Short: "Removes a host mapping entry",
	Long:  `Removes a host mapping entry.`,
	Args:  validateInput,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := parser.Parser{Path: viper.GetString(inputs.FlagFile)}
		err := p.Init()
		if err != nil {
			return errors.Wrap(err, "Unable to initialize parser")
		}

		isDomain, err := inputs.GetFlagBoolVal(cmd, inputs.FlagDomain)
		if err != nil {
			return errors.Wrapf(err, "Unable to get flag %s", inputs.FlagDomain)
		}
		if isDomain {
			err = p.RemoveDomainFromHostMapping(args[0])
			if err != nil {
				return errors.Wrap(err, "Unable to initialize parser")
			}
		} else {
			err = p.RemoveHostMapping(args[0])
			if err != nil {
				return errors.Wrap(err, "Unable to initialize parser")
			}
		}
		return nil
	},
}

func validateInput(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Please provide a valid number of inputs")
	}
	if args[0] == "" {
		return errors.New("Please provide a valid input")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP(inputs.FlagDomain, "d", false, "Delete a domain")
}
