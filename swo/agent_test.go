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

package swo

import (
	"context"
	"github.com/solarwinds/apm-go/internal/entryspans"
	"github.com/solarwinds/apm-go/internal/log"
	"github.com/solarwinds/apm-go/internal/testutils"
	"github.com/solarwinds/apm-go/internal/utils"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetGetLogLevel(t *testing.T) {
	oldLevel := GetLogLevel()

	err := SetLogLevel("INVALID")
	assert.Equal(t, err, errInvalidLogLevel)

	nl := "ERROR"
	err = SetLogLevel(nl)
	assert.Nil(t, err)

	newLevel := GetLogLevel()
	assert.Equal(t, newLevel, nl)

	SetLogLevel(oldLevel)
}

func TestShutdown(t *testing.T) {
	Shutdown(context.Background())
	assert.True(t, Closed())
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*24)
	defer cancel()
	assert.False(t, WaitForReady(ctx))
}

func TestSetLogOutput(t *testing.T) {
	oldLevel := GetLogLevel()
	_ = SetLogLevel("DEBUG")
	defer SetLogLevel(oldLevel)

	var buf utils.SafeBuffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	log.Info("hello world")
	assert.True(t, strings.Contains(buf.String(), "hello world"))
}

func TestCreateResource(t *testing.T) {
	r, err := createResource(attribute.String("foo", "bar"))
	require.NoError(t, err)
	collected := make(map[string]attribute.KeyValue)
	for _, kv := range r.Attributes() {
		collected[string(kv.Key)] = kv
	}

	require.Equal(t, attribute.String("foo", "bar"), collected["foo"])
	// These are all generated by the Otel SDK, so we check existence but not the value
	require.NotNil(t, collected["os.description"])
	require.NotNil(t, collected["os.type"])
	require.NotNil(t, collected["process.command_args"])
	require.NotNil(t, collected["process.executable.name"])
	require.NotNil(t, collected["process.executable.path"])
	require.NotNil(t, collected["process.owner"])
	require.NotNil(t, collected["process.pid"])
	require.NotNil(t, collected["telemetry.sdk.language"])
	require.NotNil(t, collected["telemetry.sdk.name"])
	require.NotNil(t, collected["telemetry.sdk.version"])
}

func TestSetTransactionName(t *testing.T) {
	err := SetTransactionName(context.Background(), "    ")
	require.Error(t, err)
	require.Equal(t, "invalid transaction name", err.Error())

	err = SetTransactionName(context.Background(), "valid")
	require.Error(t, err)
	require.Equal(t, "could not obtain OpenTelemetry SpanContext from given context", err.Error())

	tr, teardown := testutils.TracerSetup()
	defer teardown()
	ctx, s := tr.Start(context.Background(), "span name")
	err = entryspans.Push(s.(trace.ReadOnlySpan))
	require.NoError(t, err)
	err = SetTransactionName(ctx, "this should work")
	require.NoError(t, err)
	require.Equal(t, "this should work", entryspans.GetTransactionName(s.SpanContext().TraceID()))

}
