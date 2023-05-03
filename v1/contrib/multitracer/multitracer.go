// Package multitracer provides a way to run more than one OpenTracing tracers, by multiplexing calls across
// multiple Tracer, Span, and SpanContext implementations. The goal is to support a user sending data to
// two tracing vendors simultaneously (e.g., SolarWinds Observability and another implementation).
//
// Currently, baggage propagation is not supported, and the two tracers must use distinct HTTP header
// names, so as not to clobber each other.
package multitracer

import (
	ot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// MultiTracer multiplexes OpenTracing API calls across multiple Tracer implementations.
type MultiTracer struct {
	Tracers []ot.Tracer
}

// MultiSpan represents a list of Spans returned from multiple Tracers, managed by a MultiTracer.
type MultiSpan struct {
	multiTracer *MultiTracer
	Spans       []ot.Span
}

// MultiSpanContext represents a list of SpanContext returned from multiple Spans, managed by a MultiTracer.
type MultiSpanContext struct {
	multiTracer  *MultiTracer
	SpanContexts []ot.SpanContext
}

func (m *MultiTracer) StartSpan(operationName string, opts ...ot.StartSpanOption) ot.Span {
	ret := &MultiSpan{
		multiTracer: m,
		Spans:       make([]ot.Span, len(m.Tracers)),
	}
	for i, t := range m.Tracers {
		// look for ot.StartSpanOptions with a MultiSpanContext in an ot.SpanReference
		tracerSpecificOpts := []ot.StartSpanOption{}
		for _, opt := range opts {
			switch o := opt.(type) {
			case ot.SpanReference:
				// pull out tracer-specific SpanContext from MultiSpanContext
				refCtx := o.ReferencedContext.(*MultiSpanContext)
				tracerSpecificReference := ot.SpanReference{
					Type:              o.Type,
					ReferencedContext: refCtx.SpanContexts[i],
				}
				tracerSpecificOpts = append(tracerSpecificOpts, tracerSpecificReference)
			case ot.StartTime, ot.Tags, ot.Tag:
				tracerSpecificOpts = append(tracerSpecificOpts, opt)
			}
		}
		ret.Spans[i] = t.StartSpan(operationName, tracerSpecificOpts...)
	}
	return ret
}

// Inject propagates context using multiple tracers. Errors from individual tracers are dropped.
func (m *MultiTracer) Inject(sm ot.SpanContext, format interface{}, carrier interface{}) error {
	sc, ok := sm.(*MultiSpanContext)
	if !ok {
		return ot.ErrInvalidSpanContext
	}
	errs := make([]error, len(m.Tracers))
	for i, t := range m.Tracers {
		errs[i] = t.Inject(sc.SpanContexts[i], format, carrier)
	}

	// if every tracer returned an error, then also return the first error.
	var allErrors = true
	for _, err := range errs {
		if err == nil {
			allErrors = false
			break
		}
	}
	if allErrors {
		return errs[0]
	}
	return nil
}

// Extract reads context propagated using multiple tracers. Errors from individual tracers are dropped.
func (m *MultiTracer) Extract(format interface{}, carrier interface{}) (ot.SpanContext, error) {
	ret := &MultiSpanContext{
		multiTracer:  m,
		SpanContexts: make([]ot.SpanContext, len(m.Tracers)),
	}
	errs := make([]error, len(m.Tracers))
	for i, t := range m.Tracers {
		ret.SpanContexts[i], errs[i] = t.Extract(format, carrier)
	}

	// if every tracer returned an error, then also return the first error.
	var allErrors = true
	for _, err := range errs {
		if err == nil {
			allErrors = false
			break
		}
	}
	if allErrors {
		return nil, errs[0]
	}
	return ret, nil
}

func (m *MultiSpan) Finish() {
	for _, s := range m.Spans {
		s.Finish()
	}
}

func (m *MultiSpan) FinishWithOptions(opts ot.FinishOptions) {
	for _, s := range m.Spans {
		s.FinishWithOptions(opts)
	}
}

func (m *MultiSpan) Context() ot.SpanContext {
	ret := &MultiSpanContext{
		multiTracer:  m.multiTracer,
		SpanContexts: make([]ot.SpanContext, len(m.Spans)),
	}
	for i, s := range m.Spans {
		ret.SpanContexts[i] = s.Context()
	}
	return ret
}

func (m *MultiSpan) SetOperationName(operationName string) ot.Span {
	for i, s := range m.Spans {
		m.Spans[i] = s.SetOperationName(operationName)
	}
	return m
}

func (m *MultiSpan) SetTag(key string, value interface{}) ot.Span {
	for i, s := range m.Spans {
		m.Spans[i] = s.SetTag(key, value)
	}
	return m
}

func (m *MultiSpan) LogFields(fields ...log.Field) {
	for _, s := range m.Spans {
		s.LogFields(fields...)
	}
}

func (m *MultiSpan) LogKV(alternatingKeyValues ...interface{}) {
	for _, s := range m.Spans {
		s.LogKV(alternatingKeyValues...)
	}
}

// SetBaggageItem does nothing: baggage propagation is not supported.
func (m *MultiSpan) SetBaggageItem(restrictedKey, value string) ot.Span { return m }

// BaggageItem does nothing: baggage propagation is not supported.
func (m *MultiSpan) BaggageItem(restrictedKey string) string { return "" }

func (m *MultiSpan) Tracer() ot.Tracer {
	return m.multiTracer
}

func (m *MultiSpan) LogEvent(event string) {
	for _, s := range m.Spans {
		s.LogEvent(event) //nolint
	}
}

func (m *MultiSpan) LogEventWithPayload(event string, payload interface{}) {
	for _, s := range m.Spans {
		s.LogEventWithPayload(event, payload) //nolint
	}
}

func (m *MultiSpan) Log(data ot.LogData) {
	for _, s := range m.Spans {
		s.Log(data) //nolint
	}
}

// ForeachBaggageItem does nothing: baggage propagation is not supported.
func (m *MultiSpanContext) ForeachBaggageItem(handler func(k, v string) bool) {}
