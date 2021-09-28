/*
Copyright © 2021 Joseph Cracchiolo <joe@cracchiolo.com>

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
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	errMissingArgument = fmt.Errorf("missing argument")
	cfgFile            string
	flagShowHTTPResult bool
	flagVerbose        bool
	flagShowXML        bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arincli",
	Short: "CLI for ARIN API",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(viper.GetString("apikey")) == 0 {
			return fmt.Errorf("missing APIKEY in .arincli or Environment")
		}
		return nil
	},
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.arincli.yaml)")
	rootCmd.PersistentFlags().BoolVar(&flagShowHTTPResult, "dump", false, "show output of all REST calls to stderr")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "v", false, "show verbose output")
	rootCmd.PersistentFlags().BoolVar(&flagShowXML, "xml", false, "show XML returned")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".arincli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".arincli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && flagVerbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
