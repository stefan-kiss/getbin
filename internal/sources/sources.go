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
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
)

type BinSources map[string]Source

type Source struct {
	Url      string    `yaml:"url,omitempty"`
	UrlMatch string    `yaml:"urlmatch,omitempty"`
	Format   string    `yaml:"format,omitempty"`
	Owner    string    `yaml:"owner,omitempty"`
	Group    string    `yaml:"group,omitempty"`
	Perms    string    `yaml:"perms,omitempty"`
	Filemap  []FileMap `yaml:"filemap,omitempty"`
	Hashes   []BinHash `yaml:"hashes,omitempty"`
}

type BinHash struct {
	Version  string `yaml:"version,omitempty"`
	Os       string `yaml:"os,omitempty"`
	Arch     string `yaml:"arch,omitempty"`
	Hash     string `yaml:"hash,omitempty"`
	HashType string `yaml:"hashType,omitempty"`
}

type FileMap struct {
	Name          string `yaml:"name,omitempty"`
	ExtractedPath string `yaml:"-"`
	Destination   string `yaml:"destination,omitempty"`
}

func LoadSources(path string) (binSources BinSources, err error) {
	binSources = make(BinSources, 0)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &binSources)
	if err != nil {
		return nil, err
	}
	return binSources, nil
}

func WriteSources(path string, binSources BinSources) (err error) {
	bytes, err := yaml.Marshal(binSources)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func MatchUrl(url string, binSources BinSources) (label string, hash BinHash, err error) {
	for label, val := range binSources {
		if val.UrlMatch == "" {
			continue
		}
		reg, err := regexp.Compile(val.UrlMatch)
		if err != nil {
			log.Warnf("%-30s: urlmatch: unable to compile regexp: %q", label, err)
			continue
		}
		if len(reg.SubexpNames()) < 2 {
			log.Warnf("%-30s: urlmatch: no capture groups found", label)
			continue
		}
		matches := reg.FindStringSubmatch(url)
		if len(matches) < 1 {
			continue
		}
		idx := 0

		hash := BinHash{
			Version:  "",
			Os:       "",
			Arch:     "",
			Hash:     "",
			HashType: "",
		}

		idx = reg.SubexpIndex("version")
		if idx > 0 && len(matches) >= idx {
			hash.Version = matches[idx]
		}

		idx = reg.SubexpIndex("os")
		if idx > 0 && len(matches) >= idx {
			hash.Os = matches[idx]
		}

		idx = reg.SubexpIndex("arch")
		if idx > 0 && len(matches) >= idx {
			hash.Arch = matches[idx]
		}
		return label, hash, nil

	}
	return "", BinHash{}, fmt.Errorf("not found")
}
func FindHash(label string, inHash BinHash, binSources BinSources) (index int, err error) {

	srcData, ok := binSources[label]
	if !ok {
		return -1, fmt.Errorf("label not found: %s", label)
	}

	for idx, h := range srcData.Hashes {
		if inHash.Version == h.Version && inHash.Os == h.Os && inHash.Arch == h.Arch {
			return idx, nil
		}
	}

	return -1, nil
}

func AddHash(label string, inHash BinHash, inSources BinSources) (binSources BinSources, err error) {
	ret, err := FindHash(label, inHash, inSources)

	if err != nil {
		return nil, err
	}

	if ret < 0 {
		val := inSources[label]
		val.Hashes = append(inSources[label].Hashes, inHash)
		inSources[label] = val
		return inSources, nil
	}
	inSources[label].Hashes[ret] = inHash
	return inSources, nil
}
