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

package utils

import (
	"context"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/entryspans"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/testutils"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"strings"
	"testing"
)

func TestGetTransactionName(t *testing.T) {
	tr, teardown := testutils.TracerSetup()
	defer teardown()

	ctx := context.Background()
	ctx, span := tr.Start(ctx, "derived")
	roSpan, ok := span.(trace.ReadOnlySpan)
	err := entryspans.Push(roSpan)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "derived", GetTransactionName(roSpan))
	err = entryspans.SetTransactionName(span.SpanContext().TraceID(), "custom")
	require.NoError(t, err)
	require.Equal(t, "custom", GetTransactionName(roSpan))
}

func TestDeriveTransactionName(t *testing.T) {
	// Defaults to `unknown`
	var attrs []attribute.KeyValue
	name := ""
	require.Equal(t, "unknown", deriveTransactionName(name, attrs))

	// Favors span name
	name = "foo"
	require.Equal(t, name, deriveTransactionName(name, attrs))

	// Favors span name over `http.url`
	attrs = append(attrs, attribute.String("http.url", "https://user:pass@example.com/foo/bar"))
	require.Equal(t, name, deriveTransactionName(name, attrs))

	// Will use `http.url` when name is blank, and it strips user:pass
	name = ""
	require.Equal(t, "https://example.com/foo/bar", deriveTransactionName(name, attrs))

	// Will use `http.route`
	attrs = []attribute.KeyValue{
		attribute.String("http.route", "/foo/bar"),
	}
	require.Equal(t, "/foo/bar", deriveTransactionName(name, attrs))

	// Favors `http.route` over `http.url
	attrs = append(attrs, attribute.String("http.url", "https://user:pass@example.com/foo/bar"))
	require.Equal(t, "/foo/bar", deriveTransactionName(name, attrs))

	// Does not use an invalid URL
	attrs = []attribute.KeyValue{
		attribute.String("http.url", ":/"),
	}
	require.Equal(t, "unknown", deriveTransactionName(name, attrs))

	// Trims spaces
	name = " my transaction "
	attrs = []attribute.KeyValue{
		attribute.String("http.url", "https://user:pass@example.com/foo/bar"),
	}
	require.Equal(t, "my transaction", deriveTransactionName(name, attrs))

	// Truncates long transaction names
	name = strings.Repeat("a", 1024)
	expected := strings.Repeat("a", 255)
	require.Equal(t, expected, deriveTransactionName(name, attrs))
}
