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
package solarwinds_apm

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/config"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/reporter"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

const TestServiceKey = "ae38315f6116585d64d82ec2455aa3ec61e02fee25d286f74ace9e4fea189217:go"

func init() {
	os.Setenv("SW_APM_SERVICE_KEY", TestServiceKey)
}

func TestErrorSpec(t *testing.T) {
	r := reporter.SetTestReporter()

	tr := NewTrace("test")
	tr.Error("testClass", "testMsg")
	tr.End()

	r.Close(3)

	var foundErrEvt = false
	for _, evt := range r.EventBufs {
		m := make(map[string]interface{})
		bson.Unmarshal(evt, m)
		if m["Label"] != reporter.LabelError {
			continue
		}
		foundErrEvt = true
		// check error spec
		assert.Equal(t, reporter.LabelError, m["Label"])
		assert.Equal(t, "test", m["Layer"])
		assert.Equal(t, "error", m["Spec"])
		assert.Equal(t, "testMsg", m["ErrorMsg"])
		assert.Contains(t, m, "Timestamp_u")
		assert.Contains(t, m, "X-Trace")
		assert.Contains(t, m, KeyBackTrace)
		assert.Equal(t, "1", m["_V"])
		assert.Equal(t, "testClass", m["ErrorClass"])
		assert.IsType(t, "", m[KeyBackTrace])
	}
	assert.True(t, foundErrEvt)
}

func TestBeginSpan(t *testing.T) {
	r := reporter.SetTestReporter()

	ctx := NewContext(context.Background(), NewTrace("baseSpan"))
	s, _ := BeginSpan(ctx, "testSpan")

	subSpan, _ := BeginSpanWithOptions(ctx, "spanWithBackTrace", SpanOptions{WithBackTrace: true})
	subSpan.BeginSpanWithOptions("spanWithBackTrace", SpanOptions{WithBackTrace: true}).End()
	subSpan.End()

	s.End()

	EndTrace(ctx)

	r.Close(10)

	for _, evt := range r.EventBufs {
		m := make(map[string]interface{})
		bson.Unmarshal(evt, m)
		if m["Label"] != "entry" {
			continue
		}
		layer := m["Layer"]
		switch layer {
		case "testSpan":
			assert.Nil(t, m[KeyBackTrace], layer)
		case "spanWithBackTrace":
			assert.NotNil(t, m[KeyBackTrace], layer)
		case "spanWithBackTrace2":
			assert.NotNil(t, m[KeyBackTrace], layer)
		}
	}
}

func TestSpanInfo(t *testing.T) {
	r := reporter.SetTestReporter()

	ctx := NewContext(context.Background(), NewTrace("baseSpan"))
	s, _ := BeginSpan(ctx, "testSpan")

	s.InfoWithOptions(SpanOptions{WithBackTrace: true})

	s.End()
	EndTrace(ctx)

	r.Close(5)

	for _, evt := range r.EventBufs {
		m := make(map[string]interface{})
		bson.Unmarshal(evt, m)
		if m["Label"] != "info" {
			continue
		}
		layer := m["Layer"]
		switch layer {
		case "testSpan":
			assert.NotNil(t, m[KeyBackTrace], layer)
		}
	}
}

func TestFromKVs(t *testing.T) {
	assert.Equal(t, 0, len(fromKVs()))
	assert.Equal(t, 0, len(fromKVs("hello")))
	m := fromKVs("hello", 1)
	assert.Equal(t, 1, len(m))
	assert.Equal(t, 1, m["hello"])

	m = fromKVs("hello", 1, 2)
	assert.Equal(t, 1, len(m))
	assert.Equal(t, 1, m["hello"])

	m = fromKVs("hello", 1, "world", 2)
	assert.Equal(t, 2, len(m))
	assert.Equal(t, 1, m["hello"])
	assert.Equal(t, 2, m["world"])

	m = fromKVs(1, 1, 2)
	assert.Equal(t, 0, len(m))

	m = fromKVs(nil, "hello", 1)
	assert.Equal(t, 1, len(m))
	assert.Equal(t, 1, m["hello"])

	m = fromKVs(1.1, "hello", 1)
	assert.Equal(t, 1, len(m))
	assert.Equal(t, 1, m["hello"])

	m = fromKVs([]string{"1", "2"}, "hello", 1)
	assert.Equal(t, 1, len(m))
	assert.Equal(t, 1, m["hello"])
}

func TestAddKVsFromOpts(t *testing.T) {
	assert.Equal(t, 0, len(addKVsFromOpts(SpanOptions{})))

	kvs := addKVsFromOpts(SpanOptions{}, "hello")
	assert.Equal(t, []interface{}{"hello"}, kvs)

	kvs = addKVsFromOpts(SpanOptions{}, "hello", 1)
	assert.Equal(t, []interface{}{"hello", 1}, kvs)

	kvs = addKVsFromOpts(SpanOptions{WithBackTrace: true})
	assert.Equal(t, 2, len(kvs))

	kvs = addKVsFromOpts(SpanOptions{WithBackTrace: true}, "hello", 1)
	assert.Equal(t, 4, len(kvs))
}

func TestMergeKVs(t *testing.T) {
	cases := []struct {
		left   []interface{}
		right  []interface{}
		merged []interface{}
	}{
		{nil, nil, []interface{}{}},
		{nil, []interface{}{}, []interface{}{}},
		{[]interface{}{}, nil, []interface{}{}},
		{[]interface{}{}, []interface{}{}, []interface{}{}},
		{[]interface{}{"a"}, []interface{}{}, []interface{}{"a"}},
		{[]interface{}{}, []interface{}{"a"}, []interface{}{"a"}},
		{[]interface{}{"a"}, []interface{}{"b"}, []interface{}{"a", "b"}},
		{[]interface{}{"a", "b"}, []interface{}{"c", "d"}, []interface{}{"a", "b", "c", "d"}},
	}

	for idx, c := range cases {
		assert.Equal(t, c.merged, mergeKVs(c.left, c.right), fmt.Sprintf("Test case: #%d", idx))
	}
}

func TestWithTransactionFiltering(t *testing.T) {
	const layerName = "transaction_filtering"
	samplingID := "2BF4CAA9299299E3D38A58A9821BD34F6268E576CFAB2198D447EA220301"
	nonSamplingID := "2BF4CAA9299299E3D38A58A9821BD34F6268E576CFAB2198D447EA220300"

	// 1. no transaction settings
	r := reporter.SetTestReporter()
	tr := NewTrace(layerName)
	tr.End()
	r.Close(2)
	assert.Equal(t, 2, len(r.EventBufs))

	// Reload config with transaction filtering settings
	reporter.ReloadURLsConfig([]config.TransactionFilter{
		{"url", `test\d{1}`, nil, "disabled"},
		{"url", "", []string{"jpg"}, "disabled"},
	})

	// 2. “disabled” transaction settings not matched
	r = reporter.SetTestReporter()
	tr = NewTraceWithOptions(layerName, SpanOptions{WithBackTrace: false, ContextOptions: ContextOptions{URL: "/eric"}})
	tr.End()
	r.Close(2)
	assert.Equal(t, 2, len(r.EventBufs))

	// 3.1 “disabled” transaction settings matched
	r = reporter.SetTestReporter()
	tr = NewTraceWithOptions(layerName, SpanOptions{WithBackTrace: false, ContextOptions: ContextOptions{URL: "/test1"}})
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// 3.2 “disabled” transaction settings matched
	r = reporter.SetTestReporter()
	tr = NewTraceWithOptions(layerName, SpanOptions{WithBackTrace: false, ContextOptions: ContextOptions{URL: "/eric.jpg"}})
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// 4.incoming sampling xtrace + “disabled” transaction settings not matched
	r = reporter.SetTestReporter()
	tr = NewTraceFromIDForURL(layerName, samplingID, "/eric", nil)
	tr.End()
	r.Close(2)
	assert.Equal(t, 2, len(r.EventBufs))

	// 5.incoming sampling xtrace + “disabled” transaction settings matched
	r = reporter.SetTestReporter()
	tr = NewTraceFromIDForURL(layerName, samplingID, "/test2", nil)
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// 6.incoming non-sampling xtrace + “disabled” transaction settings not matched
	r = reporter.SetTestReporter()
	tr = NewTraceFromIDForURL(layerName, nonSamplingID, "/eric", nil)
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// 7. incoming non-sampling xtrace + “disabled” transaction settings matched
	r = reporter.SetTestReporter()
	tr = NewTraceFromIDForURL(layerName, nonSamplingID, "/test3", nil)
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// service level trace mode is disabled
	reporter.ReloadURLsConfig([]config.TransactionFilter{})
	os.Setenv("SW_APM_TRACING_MODE", "disabled")
	config.Load()

	// 8.no transaction settings
	r = reporter.SetTestReporter()
	tr = NewTrace(layerName)
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// service level trace mode is disabled
	reporter.ReloadURLsConfig([]config.TransactionFilter{
		{"url", `test\d{1}`, nil, "enabled"},
		{"url", "", []string{"jpg"}, "enabled"},
	})

	// 9.“enabled” transaction settings not matched
	r = reporter.SetTestReporter()
	tr = NewTraceWithOptions(layerName, SpanOptions{WithBackTrace: false, ContextOptions: ContextOptions{URL: "/eric"}})
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// 10.“enabled” transaction settings matching
	r = reporter.SetTestReporter()
	tr = NewTraceWithOptions(layerName, SpanOptions{WithBackTrace: false, ContextOptions: ContextOptions{URL: "/test1"}})
	tr.End()
	r.Close(2)
	assert.Equal(t, 2, len(r.EventBufs))

	// 11.incoming sampling xtrace + “enabled” transaction settings not matched
	r = reporter.SetTestReporter()
	tr = NewTraceFromIDForURL(layerName, samplingID, "/eric", nil)
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// 12.incoming sampling xtrace + “enabled” transaction settings matched
	r = reporter.SetTestReporter()
	tr = NewTraceFromIDForURL(layerName, samplingID, "/test2", nil)
	tr.End()
	r.Close(2)
	assert.Equal(t, 2, len(r.EventBufs))

	// 13.incoming non-sampling xtrace + “enabled” transaction settings not matched
	r = reporter.SetTestReporter()
	tr = NewTraceFromIDForURL(layerName, nonSamplingID, "/eric", nil)
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// 14.incoming non-sampling xtrace + “enabled” transaction settings matched
	r = reporter.SetTestReporter()
	tr = NewTraceFromIDForURL(layerName, nonSamplingID, "/test3", nil)
	tr.End()
	r.Close(0)
	assert.Equal(t, 0, len(r.EventBufs))

	// reset configurations
	os.Unsetenv("SW_APM_TRACING_MODE")
	config.Load()
	reporter.ReloadURLsConfig([]config.TransactionFilter{})
}
