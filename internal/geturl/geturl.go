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
	"context"
	"github.com/hashicorp/go-getter"
	"time"
)

func GetUrlEX(dest string, url string) error {
	timeout := 60
	d := time.Now().Add(time.Duration(timeout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	// Even though ctx will be expired, it is good practice to call its
	// cancellation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()

	err := getter.GetFile(dest, url, getter.WithContext(ctx))

	if err != nil {
		return err
	}

	return nil
}

func GetUrl(dest string, url string, timeout int) error {

	d := time.Now().Add(time.Duration(timeout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	// Even though ctx will be expired, it is good practice to call its
	// cancellation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()

	err := getter.GetFile(dest, url, getter.WithContext(ctx))

	if err != nil {
		return err
	}

	return nil
}
