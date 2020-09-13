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
	"github.com/stefan-kiss/getbin/internal/getsrc"
	"github.com/stefan-kiss/getbin/internal/sources"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a binary from a configured source.",
	Long: ` Get a binary from a configured source.

It needs
- label
- version
- os
- architecture

The file will be downloaded in the cache directory and then moved to the configured path relative to the configured root.

If the file exists in cache dir checksum will be checked. If it does not match it will be downloaded again.
If it's an archive it will be extracted and files moved to the configured locations (relative to the configured root).

If the file already exists in the end destination and it matches the file downloaded (or extracted from archive) it wont be overwritten. 
When (and only when) a file is copied to the destination a log line is printed. 
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Flags marked as persistent are not checked by cobra. Hence we should check them here.
		// https://github.com/spf13/cobra/issues/655
		var err error
		var label, version, os, arch, root, cache, srcFile string
		if label, err = cmd.Flags().GetString("label"); err != nil || label == "" {
			return fmt.Errorf("missing label")
		}

		if version, err = cmd.Flags().GetString("version"); err != nil || version == "" {
			return fmt.Errorf("missing version")
		}

		if os, err = cmd.Flags().GetString("os"); err != nil || os == "" {
			return fmt.Errorf("missing os")
		}

		if arch, err = cmd.Flags().GetString("arch"); err != nil || arch == "" {
			return fmt.Errorf("missing arch")
		}

		if srcFile, err = rootCmd.Flags().GetString("sources"); err != nil || srcFile == "" {
			return fmt.Errorf("missing sources")
		}

		src, err := sources.LoadSources(srcFile)
		if err != nil {
			return fmt.Errorf("unable to load: %q", err)
		}

		if root, err = cmd.Flags().GetString("root"); err != nil {
			return fmt.Errorf("error retriving files root")
		}

		if cache, err = cmd.Flags().GetString("cache"); err != nil {
			return fmt.Errorf("error retriving files cache")
		}

		copied, err := getsrc.GetSrc(label, sources.BinHash{
			Version:  version,
			Os:       os,
			Arch:     arch,
			Hash:     "",
			HashType: "",
		}, 60, root, cache, src)
		for _, file := range copied {
			log.Infof("Copied file to [ %-70s ]", file)
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")
	getCmd.PersistentFlags().StringP("label", "l", "", "source label")
	getCmd.MarkFlagRequired("label")
	getCmd.PersistentFlags().StringP("version", "V", "", "source version")
	getCmd.MarkFlagRequired("version")
	getCmd.PersistentFlags().StringP("os", "o", "", "source os")
	getCmd.MarkFlagRequired("os")
	getCmd.PersistentFlags().StringP("arch", "a", "", "source architecture")
	getCmd.MarkFlagRequired("arch")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
