/*
 * Copyright (c) 2020 Thilina Manamgoda. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package inputs handles user inputs.
package inputs

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const (
	// FlagFile is the hosts file path input flag.
	FlagFile = "file"
	// FlagIP is the IP input flag.
	FlagIP = "ip"
	// FlagDomains flag is the list of domains input flag.
	FlagDomains = "domains"
	// FlagComment is the comment input flag.
	FlagComment = "comment"
	// FlagDomain is the domain input flag.
	FlagDomain = "domain"
)

// ETCHostFilePath returns the hosts file path.
func ETCHostFilePath() string {
	return viper.GetString(FlagFile)
}

// MarkFlagRequired marks the given flag required and exists with -1 if an error was occurred.
func MarkFlagRequired(cmd *cobra.Command, flag string) {
	err := cmd.MarkFlagRequired(flag)
	if err != nil {
		fmt.Println(fmt.Errorf("must provide the flag: %s", flag))
		os.Exit(-1)
	}
}

// GetFlagStringVal method returns the String flag value.
func GetFlagStringVal(cmd *cobra.Command, flag string) (string, error) {
	return cmd.Flags().GetString(flag)
}

// GetFlagIntVal method returns the int flag value.
func GetFlagIntVal(cmd *cobra.Command, flag string) (int, error) {
	return cmd.Flags().GetInt(flag)
}

// GetFlagBoolVal method returns the Boolean flag value.
func GetFlagBoolVal(cmd *cobra.Command, flag string) (bool, error) {
	return cmd.Flags().GetBool(flag)
}

// GetFlagStringSliceVal method returns the String slice flag value.
func GetFlagStringSliceVal(cmd *cobra.Command, flag string) ([]string, error) {
	return cmd.Flags().GetStringSlice(flag)
}

// GetFlagsForAdd returns the required flag values for add command.
func GetFlagsForAdd(cmd *cobra.Command) (ip, comment string, domains []string, err error) {
	ip, err = GetFlagStringVal(cmd, FlagIP)
	if err != nil {
		return "", "", nil, err
	}
	domains, err = GetFlagStringSliceVal(cmd, FlagDomains)
	if err != nil {
		return "", "", nil, err
	}
	comment, err = GetFlagStringVal(cmd, FlagComment)
	if err != nil {
		return "", "", nil, err
	}
	return
}
