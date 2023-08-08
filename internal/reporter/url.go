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
package reporter

import (
	"github.com/solarwindscloud/solarwinds-apm-go/internal/config"
	"github.com/solarwindscloud/solarwinds-apm-go/internal/log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/coocood/freecache"
	"github.com/pkg/errors"
)

var urls *urlFilters

func init() {
	urls = newURLFilters()
	urls.LoadConfig(config.GetTransactionFiltering())
}

// ReloadURLsConfig reloads the configuration and build the transaction filtering
// filters and cache.
// This function is used for testing purpose only. It's not thread-safe.
func ReloadURLsConfig(filters []config.TransactionFilter) {
	urls.LoadConfig(filters)
	urls.cache.Clear()
}

// urlCache is a cache to store the disabled url patterns
type urlCache struct{ *freecache.Cache }

const (
	cacheExpireSeconds = 600
)

// setURLTrace sets a url and its trace decision into the cache
func (c *urlCache) setURLTrace(url string, trace tracingMode) {
	_ = c.Set([]byte(url), []byte(trace.ToString()), cacheExpireSeconds)
}

// getURLTrace gets the trace decision of a URL
func (c *urlCache) getURLTrace(url string) (tracingMode, error) {
	traceStr, err := c.Get([]byte(url))
	if err != nil {
		return TraceUnknown, err
	}

	return newTracingMode(config.TracingMode(string(traceStr))), nil
}

// urlFilter defines a URL filter
type urlFilter interface {
	match(url string) bool
	tracingMode() tracingMode
}

// regexFilter is a regular expression based URL filter
type regexFilter struct {
	regex *regexp.Regexp
	trace tracingMode
}

// newRegexFilter creates a new regexFilter instance
func newRegexFilter(regex string, mode tracingMode) (*regexFilter, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse regexp")
	}
	return &regexFilter{regex: re, trace: mode}, nil
}

// match checks if the url matches the filter
func (f *regexFilter) match(url string) bool {
	return f.regex.MatchString(url)
}

// tracingMode returns the tracing mode of this url pattern
func (f *regexFilter) tracingMode() tracingMode {
	return f.trace
}

// extensionFilter is a extension-based filter
type extensionFilter struct {
	Exts  map[string]struct{}
	trace tracingMode
}

// newExtensionFilter create a new instance of extensionFilter
func newExtensionFilter(extensions []string, mode tracingMode) *extensionFilter {
	exts := make(map[string]struct{})
	for _, ext := range extensions {
		exts[ext] = struct{}{}
	}
	return &extensionFilter{Exts: exts, trace: mode}
}

// match checks if the url matches the filter
func (f *extensionFilter) match(url string) bool {
	ext := strings.TrimLeft(filepath.Ext(url), ".")
	_, ok := f.Exts[ext]
	return ok
}

// tracingMode returns the tracing mode of this extension pattern
func (f *extensionFilter) tracingMode() tracingMode {
	return f.trace
}

type urlFilters struct {
	cache   *urlCache
	filters []urlFilter
}

func newURLFilters() *urlFilters {
	return &urlFilters{
		cache: &urlCache{freecache.NewCache(1024 * 1024)},
	}
}

// LoadConfig reads transaction filtering settings from the global configuration
func (f *urlFilters) LoadConfig(filters []config.TransactionFilter) {
	f.loadConfig(filters)
}

func (f *urlFilters) loadConfig(filters []config.TransactionFilter) {
	f.filters = nil

	for _, filter := range filters {
		if filter.RegEx != "" {
			re, err := newRegexFilter(filter.RegEx, newTracingMode(filter.Tracing))
			if err != nil {
				log.Warningf("Ignore bad regex: %s, error=", filter.RegEx, err.Error())
			}
			f.filters = append(f.filters, re)
		} else {
			f.filters = append(f.filters,
				newExtensionFilter(filter.Extensions, newTracingMode(filter.Tracing)))
		}
	}
}

// getTracingMode checks if the URL should be traced or not. It returns TraceUnknown
// if the url is not found.
func (f *urlFilters) getTracingMode(url string) tracingMode {
	if len(f.filters) == 0 || url == "" {
		return TraceUnknown
	}

	trace, err := f.cache.getURLTrace(url)
	if err == nil {
		return trace
	}

	trace = f.lookupTracingMode(url)
	f.cache.setURLTrace(url, trace)

	return trace
}

func (f *urlFilters) lookupTracingMode(url string) tracingMode {
	for _, filter := range f.filters {
		if filter.match(url) {
			return filter.tracingMode()
		}
	}
	return TraceUnknown
}
