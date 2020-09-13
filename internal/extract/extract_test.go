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

package extract

import (
	"github.com/stefan-kiss/getbin/internal/sources"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var inFilemap = []sources.FileMap{
	{
		Name:        "kubectl",
		Destination: "/usr/local/bin/kubectl",
	},
}

//func TestExtractSrc(t *testing.T) {
//	srcAll, err := sources.LoadSources("sources.yaml")
//	if err != nil {
//		t.Errorf("ExtractSrc() error loading sources = %v", err)
//	}
//
//	type args struct {
//		srcLabel       string
//		downloadedFile string
//		root string
//		srcAll         sources.BinSources
//	}
//	tests := []struct {
//		name      string
//		args      args
//		wantFiles []sources.FileMap
//		wantErr   bool
//	}{
//		{
//			name: "ExtractSucces",
//			args: args{
//				srcLabel:       "kubeclient",
//				downloadedFile: "kubernetes-client-linux-amd64.tar.gz",
//				root:	 "../../test",
//				srcAll:         srcAll,
//			},
//			wantFiles: gotFilemap,
//			wantErr: false,
//		},
//
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotFiles, err := ExtractSrc(tt.args.root, tt.args.downloadedFile, "", nil)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ExtractSrc() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
//				t.Errorf("ExtractSrc() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
//			}
//		})
//	}
//}

func TestExtractSrc(t *testing.T) {
	type args struct {
		root           string
		downloadedFile string
		format         string
		filesIn        []sources.FileMap
	}
	tests := []struct {
		name      string
		args      args
		wantFiles []sources.FileMap
		wantErr   bool
	}{
		{
			name: "ExtractSucces",
			args: args{
				root:           "../../test/cache",
				downloadedFile: "../../test/kubernetes-client-linux-amd64.tar.gz",
				format:         "tgz",
				filesIn:        inFilemap,
			},
			wantFiles: nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {

		wantFilemap := inFilemap
		testdirAbsPath, err := filepath.Abs("../../test/cache")

		if err != nil {
			t.Errorf("ExtractSrc() unable to determine absolute filepath for: %v, err %v", "../../test", err)
		}

		wantFilemap[0].ExtractedPath = filepath.Join(testdirAbsPath, "kubectl")

		t.Run(tt.name, func(t *testing.T) {
			func() {

				tt.wantFiles = wantFilemap
				gotFiles, err := ExtractSrc(tt.args.root, tt.args.downloadedFile, tt.args.format, tt.args.filesIn)
				defer os.Remove(gotFiles[0].ExtractedPath)
				if (err != nil) != tt.wantErr {
					t.Errorf("ExtractSrc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
					t.Errorf("ExtractSrc() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
				}
			}()
		})
	}
}
