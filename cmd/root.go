/*
 * Copyright (c) 2020. Stefan Kiss
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string
var logLevel string

var FilesRoot string
var DownloadCache string
var Sources string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "getbin",
	Short: "Utility to help download various software in binary distribution format.",
	Long: `Since many newer software come in binary distribution needed an utility to:

Manage download url templates (versions / architecture /os / etc). 
Download binary files into appropriate locations.
Extract archives or part of archives based on rules.
Check checksums of downloaded files.
Track files downloaded - have the option of uninstalling them.


`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/getbin.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "L", "info", "log level")
	rootCmd.PersistentFlags().StringVarP(&DownloadCache, "cache", "C", "downloader-cache", "downloads cache")
	rootCmd.PersistentFlags().StringVarP(&FilesRoot, "root", "R", "/", "root directory for downloads")
	rootCmd.PersistentFlags().StringVarP(&Sources, "sources", "S", "sources.yaml", "sources file")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		lvl = log.InfoLevel
	}

	log.SetLevel(lvl)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find cwd directory.
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in cwd directory with name "getbin" (without extension).
		viper.AddConfigPath(cwd)
		viper.SetConfigName("getbin")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
