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

package sources

import (
	"reflect"
	"testing"
)

var UnmarshaledSources = BinSources{
	"kubectl": Source{
		Url:      "https://dl.k8s.io/{{.Version}}/kubernetes-client-{{.Os}}-{{.Arch}}.tar.gz",
		UrlMatch: "https://dl.k8s.io/(?P<version>.*)/kubernetes-client-(?P<os>.*)-(?P<arch>.*).tar.gz",
		Format:   "archive",
		Owner:    "root",
		Group:    "root",
		Perms:    "0755",
		Filemap: []FileMap{
			FileMap{
				Name:        "kubectl",
				Destination: "/usr/local/bin/kubectl",
			},
		},
		Hashes: []BinHash{
			BinHash{
				Version:  "v1.19.0",
				Os:       "linux",
				Arch:     "amd64",
				Hash:     "1590d4357136a71a70172e32820c4a68430d1b94cf0ac941ea17695fbe0c5440d13e26e24a2e9ebdd360c231d4cd16ffffbbe5b577c898c78f7ebdc1d8d00fa3",
				HashType: "",
			},
			BinHash{
				Version:  "v1.19.0",
				Os:       "darwin",
				Arch:     "amd64",
				Hash:     "7093a34298297e46bcd1ccb77a9c83ca93b8ccb63ce2099d3d8cd8911ccc384470ac202644843406f031c505a8960d247350a740d683d8910ca70a0b58791a1b",
				HashType: "",
			},
		},
	},
}

func TestLoadSources(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		wantSources BinSources
		wantErr     bool
	}{
		{
			name:        "LoadSuccess",
			args:        args{"../../test/simple.srcin.yaml"},
			wantSources: UnmarshaledSources,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSources, err := LoadSources(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSources, tt.wantSources) {
				t.Errorf("LoadSources() gotSources = %v, want %v", gotSources, tt.wantSources)
			}
		})
	}
}

func TestMatchUrl(t *testing.T) {
	type args struct {
		url     string
		sources BinSources
	}
	tests := []struct {
		name      string
		args      args
		wantLabel string
		wantHash  BinHash
		wantErr   bool
	}{
		{
			name: "MatchSuccess",
			args: args{
				url:     "https://dl.k8s.io/v1.19.0/kubernetes-client-darwin-amd64.tar.gz",
				sources: UnmarshaledSources,
			},
			wantLabel: "kubectl",
			wantErr:   false,
			wantHash: BinHash{
				Version:  "v1.19.0",
				Os:       "darwin",
				Arch:     "amd64",
				Hash:     "",
				HashType: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLabel, gotHash, err := MatchUrl(tt.args.url, tt.args.sources)
			if (err != nil) != tt.wantErr {
				t.Errorf("MatchUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLabel != tt.wantLabel {
				t.Errorf("MatchUrl() gotLabel = %v, want %v", gotLabel, tt.wantLabel)
			}
			if gotHash != tt.wantHash {
				t.Errorf("MatchUrl() gotHash = %v, want %v", gotLabel, tt.wantHash)
			}
		})
	}
}

func TestFindHash(t *testing.T) {
	type args struct {
		label   string
		inHash  BinHash
		sources BinSources
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "FindSuccess", args: struct {
				label   string
				inHash  BinHash
				sources BinSources
			}{
				label:   "kubectl",
				inHash:  UnmarshaledSources["kubectl"].Hashes[0],
				sources: UnmarshaledSources,
			},
			want: 0,
		},
		{
			name: "FindSuccess", args: struct {
				label   string
				inHash  BinHash
				sources BinSources
			}{
				label:   "kubectl",
				inHash:  BinHash{},
				sources: UnmarshaledSources,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := FindHash(tt.args.label, tt.args.inHash, tt.args.sources); got != tt.want || err != nil {
				t.Errorf("FindHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndToEnd(t *testing.T) {
	sources, err := LoadSources("../../test/simple.srcin.yaml")
	if err != nil {
		t.Errorf("unable to load: %q", err)
	}

	label, hash, err := MatchUrl("https://dl.k8s.io/v1.19.0/kubernetes-client-windows-amd64.tar.gz", sources)
	if err != nil {
		t.Errorf("unable to match: %q", err)
		return
	}
	sources, err = AddHash(label, hash, sources)
	if err != nil {
		t.Errorf("unable to add: %q", err)
		return
	}

	err = WriteSources("../../test/simple.srcout.yaml", sources)
	if err != nil {
		t.Errorf("unable to write: %q", err)
		return
	}
}
