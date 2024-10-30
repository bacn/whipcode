//
//  Copyright 2024 whipcode.app (AnnikaV9)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on
//  an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific
//  language governing permissions and limitations under the License.
//

package server

import (
	"whipcode/control"
	"whipcode/podman"
)

type contextKey string

type ScopedMiddleWareParams struct {
	EnableCache  bool
	KeyAndSalt   []string
	MaxBytesSize int
	KeyStore     *control.KeyStore
	Executor     podman.Executor
}

type MiddleWareParams struct {
	RateLimiter *control.RateLimiter
	Standalone  bool
	RlBurst     int
	RlRefill    int
	Proxy       string
}
