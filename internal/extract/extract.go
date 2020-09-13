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
	"archive/tar"
	"fmt"
	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stefan-kiss/getbin/internal/sources"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func ExtractSrc(root string, downloadedFile string, format string, filesIn []sources.FileMap) (files []sources.FileMap, err error) {

	if format == "binary" {
		files = make([]sources.FileMap, 0)
		// only one file for binary. the rest is ignored.
		// so we updated the extracted path and pass this on to the next function.
		file := filesIn[0]
		file.ExtractedPath = downloadedFile
		files = append(files, file)
		return files, nil
	}

	root, err = filepath.Abs(root)

	rd, err := os.Stat(root)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(root, 0755)
		if err != nil {
			return nil, fmt.Errorf("cache directory: %s does not exist and unable to create: %w", root, err)
		}
	}
	if rd != nil && !rd.IsDir() {
		return nil, fmt.Errorf("cache directory: %s does not exist and unable to create: %w", root, err)
	}

	if err != nil {
		return nil, err
	}
	if format == "tgz" {
		return extractTgz(root, downloadedFile, filesIn)
	}

	return nil, fmt.Errorf("unknown archive format: %s", format)
}

func extractTgz(root string, downloadedFile string, filesIn []sources.FileMap) (files []sources.FileMap, err error) {
	tgz := archiver.NewTarGz()
	tgz.ContinueOnError = false

	err = tgz.Walk(downloadedFile, func(f archiver.File) error {
		th, ok := f.Header.(*tar.Header)
		if !ok {
			return fmt.Errorf("expected header to be *tar.Header but was %T", f.Header)
		}

		name := path.Clean(th.Name)
		log.Debugf("found file in archive: %s", name)
		if f.IsDir() {
			return nil
		}
		for _, fileDef := range filesIn {
			matched, _ := regexp.MatchString(fileDef.Name, name)
			if !matched {
				continue
			}
			fileName := filepath.Join(root, filepath.Base(name))

			status, err := os.Stat(fileName)

			if err == nil && status.Size() == th.Size {
				log.Debugf("skipping extracting %s already exists and has same size", fileName)
			} else {

				wh, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
				if err != nil {
					return fmt.Errorf("unable to open destination file: %s: %w", "kubectl", err)
				}
				defer wh.Close()
				_, err = io.Copy(wh, f.ReadCloser)

				if err != nil {
					return fmt.Errorf("unable to write to file: %s: %w", "kubectl", err)
				}

				if err != nil {
					return fmt.Errorf("unable to write to file: %s: %w", "kubectl", err)
				}
			}

			files = append(files, sources.FileMap{
				Name:          fileDef.Name,
				ExtractedPath: fileName,
				Destination:   fileDef.Destination,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
