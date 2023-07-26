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

package swohttp

import (
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/reporter"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/swotel"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"regexp"

	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	ACEHdr                = "Access-Control-Expose-Headers"
	XTrace                = "X-Trace"
	XTraceOptionsResponse = "X-Trace-Options-Response"
)

var xtraceRegexp = regexp.MustCompile(`\A00-[[:xdigit:]]{32}-[[:xdigit:]]{16}-01\z`)

func TestHandlerNoXOptsResponse(t *testing.T) {
	r := reporter.SetTestReporter(reporter.TestReporterSettingType(reporter.DefaultST))
	defer r.Close(0)

	cb, err := solarwinds_apm.Start()
	require.NoError(t, err)
	defer cb()

	resp := doRequest(t, "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, XTrace, resp.Header.Get(ACEHdr), XTrace)
	require.Regexp(t, xtraceRegexp, resp.Header.Get(XTrace))
}

func TestHandlerWithXOptsResponse(t *testing.T) {
	r := reporter.SetTestReporter(reporter.TestReporterSettingType(reporter.DefaultST))
	defer r.Close(0)

	cb, err := solarwinds_apm.Start()
	require.NoError(t, err)
	defer cb()

	resp := doRequest(t, "trigger-trace")
	ts := trace.TraceState{}
	ts, err = swotel.SetInternalState(ts, swotel.XTraceOptResp, "my internal state")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, XTrace+","+XTraceOptionsResponse, resp.Header.Get(ACEHdr))
	require.Regexp(t, xtraceRegexp, resp.Header.Get(XTrace))
	require.Regexp(t, "trigger-trace=ok", resp.Header.Get(XTraceOptionsResponse))
}

func doRequest(t *testing.T, xtOpts string) *http.Response {
	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("foo bar baz"))
		require.NoError(t, err)
	})
	handler = NewHandler(handler)
	handler = otelhttp.NewHandler(handler, "foobar")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("", "/", http.NoBody)
	if xtOpts != "" {
		req.Header.Add("X-Trace-Options", xtOpts)
	}
	handler.ServeHTTP(recorder, req)
	return recorder.Result()
}
