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
package solarwinds_apm_test

import (
	"testing"
	"time"

	"context"

	solarwinds_apm "github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm"
	g "github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/graphtest"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/reporter"
	"github.com/stretchr/testify/assert"
)

func childExample(ctx context.Context) {
	// create a new trace, and a context to carry it around
	l1, _ := solarwinds_apm.BeginSpan(ctx, "L1")
	l2 := l1.BeginSpan("DBx", "Query", "SELECT * FROM tbl")
	time.Sleep(20 * time.Millisecond)
	l2.End()
	l1.End()

	// test attempting to start a child from a span that has ended
	// currently we don't allow this, so nothing should be reported
	l3 := l1.BeginSpan("invalidSpan", "notReported", true)
	l3.End()

	// test attempting to start a profile from a span that has ended
	// similarly we don't allow this, so nothing should be reported

	// end the trace
	solarwinds_apm.EndTrace(ctx)
}

func childExampleCtx(ctx context.Context) {
	// create a new trace, and a context to carry it around
	_, ctxL1 := solarwinds_apm.BeginSpan(ctx, "L1")
	_, ctxL2 := solarwinds_apm.BeginSpan(ctxL1, "DBx", "Query", "SELECT * FROM tbl")
	time.Sleep(20 * time.Millisecond)
	solarwinds_apm.End(ctxL2)
	solarwinds_apm.End(ctxL1)

	// test attempting to start a child from a span that has ended
	// currently we don't allow this, so nothing should be reported
	_, ctxL3 := solarwinds_apm.BeginSpan(ctxL1, "invalidSpan", "notReported", true)
	solarwinds_apm.End(ctxL3)

	// end the trace
	solarwinds_apm.EndTrace(ctx)
}

func assertTraceChild(t *testing.T, bufs [][]byte) {
	// validate events reported
	g.AssertGraph(t, bufs, 6, g.AssertNodeMap{
		{"childExample", "entry"}: {},
		{"L1", "entry"}:           {Edges: g.Edges{{"childExample", "entry"}}},
		{"DBx", "entry"}:          {Edges: g.Edges{{"L1", "entry"}}},
		{"DBx", "exit"}:           {Edges: g.Edges{{"DBx", "entry"}}},
		{"L1", "exit"}:            {Edges: g.Edges{{"DBx", "exit"}, {"L1", "entry"}}},
		{"childExample", "exit"}:  {Edges: g.Edges{{"L1", "exit"}, {"childExample", "entry"}}},
	})
}

func TestTraceChild(t *testing.T) {
	r := reporter.SetTestReporter() // enable test reporter
	ctx := solarwinds_apm.NewContext(context.Background(), solarwinds_apm.NewTrace("childExample"))
	childExample(ctx) // generate events
	r.Close(6)
	assertTraceChild(t, r.EventBufs)
}

func TestTraceChildCtx(t *testing.T) {
	r := reporter.SetTestReporter() // enable test reporter
	ctx := solarwinds_apm.NewContext(context.Background(), solarwinds_apm.NewTrace("childExample"))
	childExampleCtx(ctx) // generate events
	r.Close(6)
	assertTraceChild(t, r.EventBufs)
}

func TestNoTraceChild(t *testing.T) {
	r := reporter.SetTestReporter()
	ctx := context.Background()
	childExample(ctx)
	assert.Len(t, r.EventBufs, 0)

	r = reporter.SetTestReporter()
	ctx = context.Background()
	childExampleCtx(ctx)
	assert.Len(t, r.EventBufs, 0)
}
