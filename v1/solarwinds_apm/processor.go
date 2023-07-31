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
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/entryspans"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/log"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/metrics"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewInboundMetricsSpanProcessor(isAppoptics bool) sdktrace.SpanProcessor {
	return &inboundMetricsSpanProcessor{
		isAppoptics: isAppoptics,
	}
}

var _ sdktrace.SpanProcessor = &inboundMetricsSpanProcessor{}

var recordFunc = metrics.RecordSpan

type inboundMetricsSpanProcessor struct {
	isAppoptics bool
}

func (s *inboundMetricsSpanProcessor) OnStart(_ context.Context, span sdktrace.ReadWriteSpan) {
	if entryspans.IsEntrySpan(span) {
		if err := entryspans.Push(span); err != nil {
			// The only error here should be if it's not an entry span, and we've guarded against that,
			// so it's safe to log the error and move on
			log.Warningf("could not push entry span: %s", err)
		}
	}
}

func popEntrySpan(span sdktrace.ReadOnlySpan) {
	if sid, ok := entryspans.Pop(span.SpanContext().TraceID()); !ok {
		log.Warningf("could not pop entry span!")
	} else if sid != span.SpanContext().SpanID() {
		log.Warningf("span ids did not match! wanted %s, got %s", span.SpanContext().SpanID(), sid)
	}
}

func (s *inboundMetricsSpanProcessor) OnEnd(span sdktrace.ReadOnlySpan) {
	if entryspans.IsEntrySpan(span) {
		defer popEntrySpan(span)
		recordFunc(span, s.isAppoptics)
	}
}

func (s *inboundMetricsSpanProcessor) Shutdown(context.Context) error {
	return nil
}

func (s *inboundMetricsSpanProcessor) ForceFlush(context.Context) error {
	return nil
}
