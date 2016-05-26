// Copyright (C) 2016 AppNeta, Inc. All rights reserved.
// TraceView HTTP instrumentation for Go

package tv

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

var httpHandlerLayerName = "http.HandlerFunc"

// HTTPHandler wraps an http handler function with entry / exit events,
// returning a new function that can be used in its place.
func HTTPHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	// At wrap time (when binding handler to router): get name of wrapped handler func
	var endArgs []interface{}
	if f := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()); f != nil {
		// e.g. "main.slowHandler", "github.com/appneta/go-appneta/v1/tv_test.handler404"
		fname := f.Name()
		if s := strings.SplitN(fname[strings.LastIndex(fname, "/")+1:], ".", 2); len(s) == 2 {
			endArgs = append(endArgs, "Controller", s[0], "Action", s[1])
		}
	}
	// return wrapped HTTP request handler
	return func(w http.ResponseWriter, r *http.Request) {
		t, w := TraceFromHTTPRequestResponse(httpHandlerLayerName, w, r)
		defer t.End(endArgs...)
		// Call original HTTP handler
		handler(w, r)
	}
}

func TraceFromHTTPRequestResponse(layerName string, w http.ResponseWriter, r *http.Request) (Trace, *httpResponseWriter) {
	t := TraceFromHTTPRequest(layerName, r)
	wrapper := NewResponseWriter(w, t) // wrap writer with response-observing writer
	return t, wrapper
}

// httpResponseWriter observes an http.ResponseWriter when WriteHeader is called to check
// the status code and response headers.
type httpResponseWriter struct {
	http.ResponseWriter
	t      Trace
	Status int
}

func (w *httpResponseWriter) WriteHeader(status int) {
	w.Status = status               // observe HTTP status code
	md := w.Header().Get("X-Trace") // check response for downstream metadata
	if w.t.IsTracing() {            // set trace exit metadata in X-Trace header
		// if downstream response headers mention a different layer, add edge to it
		if md != "" && md != w.t.ExitMetadata() {
			w.t.AddEndArgs("Edge", md)
		}
		w.Header().Set("X-Trace", w.t.ExitMetadata()) // replace downstream MD with ours
	}
	w.ResponseWriter.WriteHeader(status)
}

// NewResponseWriter observes the HTTP Status code of an HTTP response, returning a
// wrapped http.ResponseWriter and a pointer to an int containing the status.
func NewResponseWriter(writer http.ResponseWriter, t Trace) *httpResponseWriter {
	w := &httpResponseWriter{writer, t, http.StatusOK}
	t.AddEndArgs("Status", &w.Status)
	// add exit event metadata to X-Trace header
	if t.IsTracing() {
		// add/replace response header metadata with this trace's
		w.Header().Set("X-Trace", t.ExitMetadata())
	}
	return w
}

// TraceFromHTTPRequest returns a Trace, given an http.Request. If a distributed trace is described
// in the "X-Trace" header, this context will be continued.
func TraceFromHTTPRequest(layerName string, r *http.Request) Trace {
	// start trace, passing in metadata header
	t := NewTraceFromID(layerName, r.Header.Get("X-Trace"), func() KVMap {
		return KVMap{
			"Method":       r.Method,
			"HTTP-Host":    r.Host,
			"URL":          r.URL.Path,
			"Remote-Host":  r.RemoteAddr,
			"Query-String": r.URL.RawQuery,
		}
	})
	// update incoming metadata in request headers for any downstream readers
	r.Header.Set("X-Trace", t.MetadataString())
	return t
}
