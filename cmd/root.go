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
	"fmt"
	"github.com/ThilinaManamgoda/etchosts/pkg/inputs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	// EnvPrefix is the environment variable prefix.
	EnvPrefix = "ETCHOSTS"
	//DefaultETCHostsFilePath is the default host entry path.
	DefaultETCHostsFilePath = "/etc/hosts"
)

// Version of the password manager. Should be initialized at build time.
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "etchosts",
	Short: "A CLI to manage /etc/host file",
	Long:  `A simple CLi to add, search, remove host entries in the /etc/host file.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP(inputs.FlagFile, "f", DefaultETCHostsFilePath,
		"Hosts file path(Default is "+DefaultETCHostsFilePath)
	err := viper.BindPFlag(inputs.FlagFile, rootCmd.PersistentFlags().Lookup(inputs.FlagFile))
	if err != nil {
		fmt.Println(err)
	}
	rootCmd.Version = Version
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("file", DefaultETCHostsFilePath)
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}
