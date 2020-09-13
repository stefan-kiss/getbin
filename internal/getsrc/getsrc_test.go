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

package getsrc

import (
	"github.com/stefan-kiss/getbin/internal/sources"
	"io/ioutil"
	"os"
	"testing"
)

func TestGetSrc(t *testing.T) {

	cacheDir, err := ioutil.TempDir(os.TempDir(), "getbin_")
	if err != nil {
		t.Errorf("unable to create temp cacheDir: %w", err)
		return
	}

	defer os.RemoveAll(cacheDir)

	srcAll, err := sources.LoadSources("../../test/sources.yaml")
	if err != nil {
		t.Errorf("unable to load sources.yaml")
		return
	}
	version := sources.BinHash{
		Version:  "1.19.0",
		Os:       "linux",
		Arch:     "amd64",
		Hash:     "",
		HashType: "",
	}
	type args struct {
		srcLabel string
		version  sources.BinHash
		timeout  int
		root     string
		cache    string
		srcAll   sources.BinSources
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetSucces", args: struct {
				srcLabel string
				version  sources.BinHash
				timeout  int
				root     string
				cache    string
				srcAll   sources.BinSources
			}{
				srcLabel: "crictl",
				version:  version,
				timeout:  60,
				root:     "../../test",
				cache:    cacheDir,
				srcAll:   srcAll,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GetSrc(tt.args.srcLabel, tt.args.version, tt.args.timeout, tt.args.root, tt.args.cache, tt.args.srcAll); (err != nil) != tt.wantErr {
				t.Errorf("GetSrc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
