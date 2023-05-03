// Copyright (C) 2023 SolarWinds Worldwide, LLC. All rights reserved.
// test usage example from doc.go

package solarwinds_apm_test

import (
	"errors"
	"testing"

	"context"

	solarwinds_apm "github.com/solarwindscloud/solarwinds-apm-go/v1/ao"
	g "github.com/solarwindscloud/solarwinds-apm-go/v1/ao/internal/graphtest"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/ao/internal/reporter"
	"github.com/stretchr/testify/assert"
)

func testDocSpanExample() {
	// create trace and bind to new context
	ctx := solarwinds_apm.NewContext(context.Background(), solarwinds_apm.NewTrace("myApp"))
	// create new span for this trace
	l, ctxL := solarwinds_apm.BeginSpan(ctx, "mySpan")

	// Start a new span, given a parent span
	db1L := l.BeginSpan("myDB1", "Query", "SELECT * FROM tbl1")
	// perform a query
	db1L.End()

	// Start a new span, given a context.Context
	db2L, _ := solarwinds_apm.BeginSpan(ctxL, "myDB2", "Query", "SELECT * FROM tbl2")
	// perform a query
	db2L.End()

	l.End()
	solarwinds_apm.Err(ctx, errors.New("Got bad error!"))
	solarwinds_apm.EndTrace(ctx)
}

func testDocSpanExampleCtx() {
	// create trace and bind to new context
	ctx := solarwinds_apm.NewContext(context.Background(), solarwinds_apm.NewTrace("myApp"))
	// create new span for this trace
	_, ctxL := solarwinds_apm.BeginSpan(ctx, "mySpan")

	// Start a new span, given a parent span
	_, ctxQ1 := solarwinds_apm.BeginSpan(ctxL, "myDB1", "Query", "SELECT * FROM tbl1")
	// perform a query
	solarwinds_apm.End(ctxQ1)

	// Start a new span, given a context.Context
	_, ctxQ2 := solarwinds_apm.BeginSpan(ctxL, "myDB2", "Query", "SELECT * FROM tbl2")
	// perform a query
	solarwinds_apm.End(ctxQ2)

	solarwinds_apm.End(ctxL)
	solarwinds_apm.Err(ctx, errors.New("Got bad error!"))
	solarwinds_apm.EndTrace(ctx)
}

func TestDocSpanExample(t *testing.T) {
	r := reporter.SetTestReporter()
	testDocSpanExample()
	r.Close(11)
	assertDocSpanExample(t, r.EventBufs)
}
func TestDocSpanExampleCtx(t *testing.T) {
	r := reporter.SetTestReporter()
	testDocSpanExampleCtx()
	r.Close(11)
	assertDocSpanExample(t, r.EventBufs)
}
func assertDocSpanExample(t *testing.T, bufs [][]byte) {
	g.AssertGraph(t, bufs, 9, g.AssertNodeMap{
		{"myApp", "entry"}:  {},
		{"mySpan", "entry"}: {Edges: g.Edges{{"myApp", "entry"}}},
		{"myDB1", "entry"}: {Edges: g.Edges{{"mySpan", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "SELECT * FROM tbl1", n.Map["Query"])
		}},
		{"myDB1", "exit"}: {Edges: g.Edges{{"myDB1", "entry"}}},
		{"myDB2", "entry"}: {Edges: g.Edges{{"mySpan", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "SELECT * FROM tbl2", n.Map["Query"])
		}},
		{"myDB2", "exit"}:  {Edges: g.Edges{{"myDB2", "entry"}}},
		{"mySpan", "exit"}: {Edges: g.Edges{{"myDB1", "exit"}, {"myDB2", "exit"}, {"mySpan", "entry"}}},
		{"myApp", "error"}: {Edges: g.Edges{{"myApp", "entry"}}, Callback: func(n g.Node) {
			assert.Equal(t, "error", n.Map["ErrorClass"])
			assert.Equal(t, "Got bad error!", n.Map["ErrorMsg"])
		}},
		{"myApp", "exit"}: {Edges: g.Edges{{"mySpan", "exit"}, {"myApp", "error"}}},
	})
}
