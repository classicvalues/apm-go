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

package w3cfmt

import (
	"encoding/hex"
	"fmt"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/log"
	"regexp"

	"go.opentelemetry.io/otel/trace"
)

var swTraceStateRegex = regexp.MustCompile(`^([[:xdigit:]]{16})-([[:xdigit:]]{2})$`)

func SwFromCtx(sc trace.SpanContext) string {
	spanID := sc.SpanID()
	traceFlags := sc.TraceFlags()
	return fmt.Sprintf("%x-%x", spanID[:], []byte{byte(traceFlags)})
}

func ParseSwTraceState(s string) SwTraceState {
	matches := swTraceStateRegex.FindStringSubmatch(s)
	if matches != nil {
		flags, err := hex.DecodeString(matches[2])
		if err != nil {
			log.Debug("Could not decode hex!", matches[2])
			matches = nil
		}
		return SwTraceState{isValid: true, spanId: matches[1], flags: trace.TraceFlags(flags[0])}
	}

	return SwTraceState{isValid: false, spanId: "", flags: 0x00}
}

type SwTraceState struct {
	isValid bool
	spanId  string
	flags   trace.TraceFlags
}

func (s SwTraceState) IsValid() bool {
	return s.isValid
}

func (s SwTraceState) SpanId() string {
	return s.spanId
}

func (s SwTraceState) Flags() trace.TraceFlags {
	return s.flags
}
