// © 2023 SolarWinds Worldwide, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package xtrace

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/log"
)

const (
	XTraceOptionsHeaderName    = "x-trace-options"
	XTraceOptionsSigHeaderName = "x-trace-options-signature"
)

var optRegex = regexp.MustCompile(";+")
var customKeyRegex = regexp.MustCompile(`^custom-[^\s]*$`)

type XTraceOptions interface {
	SwKeys() string
	CustomKVs() map[string]string
	Timestamp() int64
	TriggerTrace() bool
	Signature() string
}

func NewXTraceOptions(opts string, sig string) XTraceOptions {
	return xTraceOptions{
		opts:        opts,
		sig:         sig,
		initialized: false,
		swKeys:      "",
		customKVs:   make(map[string]string),
		timestamp:   0,
	}
}

type xTraceOptions struct {
	opts        string
	sig         string
	initialized bool
	swKeys      string
	customKVs   map[string]string
	timestamp   int64
	tt          bool
}

func (x *xTraceOptions) init() {
	if x.opts != "" {
		x.extractOpts()
	}

	x.initialized = true
}

func (x *xTraceOptions) extractOpts() {
	opts := optRegex.Split(x.opts, -1)
	for _, opt := range opts {
		k, v, found := strings.Cut(opt, "=")
		k = strings.TrimSpace(k)
		if !found {
			// Only support trigger-trace without an equals sign
			if k == "trigger-trace" {
				x.tt = true
			}
			continue
		}
		v = strings.TrimSpace(v)
		if k == "sw-keys" {
			x.swKeys = v
		} else if k == "ts" {
			ts, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				log.Debug("Invalid x-trace timestamp value", ts)
			} else {
				x.timestamp = ts
			}
		} else if k == "trigger-trace" {
			log.Debug("trigger-trace must be standalone flag, ignoring.")
		} else if customKeyRegex.MatchString(k) {
			x.customKVs[k] = strings.TrimSpace(v)
		}

	}
}

func (x xTraceOptions) SwKeys() string {
	if !x.initialized {
		x.init()
	}
	return x.swKeys
}

func (x xTraceOptions) CustomKVs() map[string]string {
	if !x.initialized {
		x.init()
	}
	return x.customKVs
}

func (x xTraceOptions) Timestamp() int64 {
	if !x.initialized {
		x.init()
	}
	return x.timestamp
}

func (x xTraceOptions) TriggerTrace() bool {
	if !x.initialized {
		x.init()
	}
	return x.tt
}

func (x xTraceOptions) Signature() string {
	// This is set on instantiation, no need to initialize
	return x.sig
}
