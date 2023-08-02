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
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/reporter"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

func TestLoggableTraceIDFromContext(t *testing.T) {
	r := reporter.SetTestReporter(reporter.TestReporterSettingType(reporter.DefaultST))
	defer r.Close(0)

	ctx := context.Background()
	lt := LoggableTrace(ctx)
	require.Equal(t, LoggableTraceContext{
		TraceID:     trace.TraceID{},
		SpanID:      trace.SpanID{},
		TraceFlags:  0,
		ServiceName: "test-reporter-service",
	}, lt)
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x22},
		SpanID:     trace.SpanID{0x11},
		TraceFlags: trace.FlagsSampled,
	})
	require.False(t, lt.IsValid())
	require.Equal(t,
		"trace_id=00000000000000000000000000000000 span_id=0000000000000000 trace_flags=00 resource.service.name=test-reporter-service",
		lt.String())

	ctx = trace.ContextWithSpanContext(ctx, sc)
	lt = LoggableTrace(ctx)
	require.Equal(t, LoggableTraceContext{
		TraceID:     sc.TraceID(),
		SpanID:      sc.SpanID(),
		TraceFlags:  sc.TraceFlags(),
		ServiceName: "test-reporter-service",
	}, lt)
	require.True(t, lt.IsValid())
	require.Equal(t,
		"trace_id=22000000000000000000000000000000 span_id=1100000000000000 trace_flags=01 resource.service.name=test-reporter-service",
		lt.String())

	sc = trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x33},
		SpanID:     trace.SpanID{0xAA},
		TraceFlags: trace.FlagsSampled,
	})
	ctx = trace.ContextWithSpanContext(ctx, sc)
	lt = LoggableTrace(ctx)
	require.Equal(t, LoggableTraceContext{
		TraceID:     sc.TraceID(),
		SpanID:      sc.SpanID(),
		TraceFlags:  sc.TraceFlags(),
		ServiceName: "test-reporter-service",
	}, lt)
	require.True(t, lt.IsValid())
	require.Equal(t,
		"trace_id=33000000000000000000000000000000 span_id=aa00000000000000 trace_flags=01 resource.service.name=test-reporter-service",
		lt.String())
}
