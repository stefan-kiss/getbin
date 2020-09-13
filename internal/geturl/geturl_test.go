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

package geturl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGetUrl(t *testing.T) {
	type args struct {
		dest string
		url  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GetSucces",
			args: args{
				dest: "1Mio.dat",
				url:  "http://www.ovh.net/files/1Mio.dat?checksum=sha512:d2b1a1688a1ea4476c797b4e6d0d5f7992cbd144e47ebef9710a6c8b882068cb984d5da703f0988249e61648e60a130e2b7ddeafbec14c220c18a19475db5816",
			},
			wantErr: false,
		},
	}
	dir, err := ioutil.TempDir(os.TempDir(), "getbin_")
	if err != nil {
		t.Errorf("unable to create temp dir: %w", err)
		return
	}

	defer os.RemoveAll(dir)
	filedest := filepath.Join(dir, "testfile")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetUrlEX(filedest+tt.args.dest, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("GetUrlEX() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
