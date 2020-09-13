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
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stefan-kiss/getbin/internal/extract"
	"github.com/stefan-kiss/getbin/internal/geturl"
	"github.com/stefan-kiss/getbin/internal/sources"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

func GetSrc(srcLabel string, version sources.BinHash, timeout int, root string, cache string, srcAll sources.BinSources) (target []string, err error) {

	// log.Debugf("cache: %s", pp.Sprint(viper.Get("cache")))
	// log.Debugf("root: %s", pp.Sprint(viper.Get("root")))
	idx, err := sources.FindHash(srcLabel, version, srcAll)

	if err != nil {
		return nil, err
	}
	if idx < 0 {
		return nil, fmt.Errorf("given version/os/arch for label: %s not found in config", srcLabel)
	}
	tUrl := srcAll[srcLabel].Url
	tmpl, err := template.New("urlTpl").Parse(tUrl)
	var w bytes.Buffer
	tmpl.Execute(&w, version)

	url := string(w.Bytes())

	filename := filepath.Base(url)

	sHash := srcAll[srcLabel].Hashes[idx].Hash

	destdir := filepath.Join(cache, srcLabel, version.Version, version.Os, version.Arch)
	err = os.MkdirAll(destdir, 0755)
	if err != nil {
		return nil, err
	}

	dest := filepath.Join(destdir, filename)
	urlExt := fmt.Sprintf("%s?checksum=sha512:%s&archive=false", url, sHash)

	err = geturl.GetUrl(dest, urlExt, timeout)
	if err != nil {
		return nil, err
	}
	extractDir := filepath.Dir(dest)

	extractFiles, err := extract.ExtractSrc(extractDir, dest, srcAll[srcLabel].Format, srcAll[srcLabel].Filemap)
	destinations, err := copyFiles(extractFiles, root)
	if err != nil {
		return nil, err
	}
	return destinations, nil
}

func copyFiles(files []sources.FileMap, root string) (destinations []string, err error) {
	destinations = make([]string, 0)
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("could not determine absolute path for: %s :%w", root, err)
	}
	rd, err := os.Stat(root)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(root, 0755)
		if err != nil {
			return nil, fmt.Errorf("root directory: %s does not exist and unable to create: %w", root, err)
		}
	}
	if rd != nil && !rd.IsDir() {
		return nil, fmt.Errorf("root directory: %s does not exist and unable to create: %w", root, err)
	}

	for _, file := range files {
		src, err := os.Stat(file.ExtractedPath)
		if err != nil {
			return nil, fmt.Errorf("source file: %s does not exist: %w", file.ExtractedPath, err)
		}
		file.Destination = filepath.Join(root, file.Destination)
		dst, err := os.Stat(file.Destination)
		if err == nil && src.Size() == dst.Size() {
			log.Debugf("skipping copy %s file already exists and has the same size", file.Destination)
			continue
		}

		if err != nil {
			parent := filepath.Dir(file.Destination)
			_, err := os.Stat(parent)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(parent, 0755)
					if err != nil {
						return nil, fmt.Errorf("parent directory: %s does not exist and unable to create: %w", root, err)
					}
				} else {
					return nil, fmt.Errorf("error for directory: %s : %w", root, err)
				}
			}
		}
		err = copy(file.ExtractedPath, file.Destination)
		if err != nil {
			return nil, fmt.Errorf("error coppying file to: %s : %w", file.Destination, err)
		}
		destinations = append(destinations, file.Destination)
	}
	return destinations, nil
}

// quick copy from StackOverflow :)))

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copy(src, dst string) error {
	log.Debugf("copy: %s to %s", src, dst)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
