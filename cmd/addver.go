/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	log "github.com/sirupsen/logrus"
	"github.com/stefan-kiss/getbin/internal/sources"
	"github.com/stefan-kiss/getbin/internal/urlchecksum"

	"github.com/spf13/cobra"
)

// addverCmd represents the addver command
var addverCmd = &cobra.Command{
	Use:   "addver",
	Short: "Add a version for one of the existing sources.",
	Long: `Add a version for one of the existing sources.
The version is passed as: URL. 
Version/os/arch will be extracted from it according to the regexp for that source.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var label, url, hash, srcFile string
		// Flags marked as persistent are not checked by cobra. Hence we should check them here.
		// https://github.com/spf13/cobra/issues/655

		if url, err = cmd.Flags().GetString("url"); err != nil || url == "" {
			return fmt.Errorf("missing url")
		}

		if hash, err = cmd.Flags().GetString("hash"); err != nil || hash == "" {
			hash, err = urlchecksum.Calculate(url)
			if err != nil {
				return fmt.Errorf("error downloading %s: %w", url, err)
			}
		}
		fmt.Printf("Adding Version: [  %s  ]/[  %s  ]\n", url, hash)

		if srcFile, err = rootCmd.Flags().GetString("sources"); err != nil || srcFile == "" {
			return fmt.Errorf("missing sources")
		}

		src, err := sources.LoadSources(srcFile)

		if err != nil {
			return fmt.Errorf("unable to load: %q", err)
		}

		label, hashDef, err := sources.MatchUrl(url, src)
		if err != nil {
			return fmt.Errorf("unable to match: %q", err)
		}

		hashDef.Hash = hash
		src, err = sources.AddHash(label, hashDef, src)
		if err != nil {
			return fmt.Errorf("unable to add: %q", err)
		}

		err = sources.WriteSources("sources.yaml", src)
		if err != nil {
			return fmt.Errorf("unable to write: %q", err)

		}
		log.Infof("Written file: %s", "sources.yaml")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addverCmd)

	addverCmd.PersistentFlags().StringP("url", "u", "", "Download URL")
	addverCmd.MarkFlagRequired("url")

	addverCmd.PersistentFlags().StringP("hash", "H", "", "Hash for artifact")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
