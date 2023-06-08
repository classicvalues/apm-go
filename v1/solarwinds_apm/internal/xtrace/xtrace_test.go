package xtrace_test

import (
	"testing"

	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/xtrace"
	"github.com/stretchr/testify/assert"
)

func TestNoKeyNoValue(t *testing.T) {
	xto := xtrace.NewXTraceOptions("=", "")
	assert.Empty(t, xto.CustomKVs())
	assert.Empty(t, xto.SwKeys())
}

func TestOrphanValue(t *testing.T) {
	xto := xtrace.NewXTraceOptions("=oops", "")
	assert.Empty(t, xto.CustomKVs())
	assert.Empty(t, xto.SwKeys())
}

func TestValidTT(t *testing.T) {
	xto := xtrace.NewXTraceOptions("trigger-trace", "")
	assert.True(t, xto.TriggerTrace())
	assert.Empty(t, xto.CustomKVs())
	assert.Empty(t, xto.SwKeys())
}

func TestTTKeyIgnored(t *testing.T) {
	xto := xtrace.NewXTraceOptions("trigger-trace=1", "")
	assert.False(t, xto.TriggerTrace())
	assert.Empty(t, xto.CustomKVs())
	assert.Empty(t, xto.SwKeys())
}

func TestSwKeysKVStrip(t *testing.T) {
	xto := xtrace.NewXTraceOptions("sw-keys=   foo:key   ", "")
	assert.Equal(t, "foo:key", xto.SwKeys())
}

func TestSwKeysContainingSemicolonIgnoreAfter(t *testing.T) {
	xto := xtrace.NewXTraceOptions("sw-keys=check-id:check-1013,website-id;booking-demo", "")
	assert.Equal(t, "check-id:check-1013,website-id", xto.SwKeys())
}

func TestCustomKeysMatchStoredInOptionsHeaderAndCustomKVs(t *testing.T) {
	xto := xtrace.NewXTraceOptions("custom-awesome-key=    foo ", "")
	assert.Equal(t, map[string]string{"custom-awesome-key": "foo"}, xto.CustomKVs())
}

func TestCustomKeysMatchButNoValueIgnored(t *testing.T) {
	xto := xtrace.NewXTraceOptions("custom-no-value", "")
	assert.Equal(t, map[string]string{}, xto.CustomKVs())
}

func TestCustomKeysMatchEqualInValue(t *testing.T) {
	xto := xtrace.NewXTraceOptions("custom-and=a-value=12345containing_equals=signs", "")
	assert.Equal(t, map[string]string{"custom-and": "a-value=12345containing_equals=signs"}, xto.CustomKVs())
}

func TestCustomKeysSpacesInKeyDisallowed(t *testing.T) {
	xto := xtrace.NewXTraceOptions("custom- key=this_is_bad;custom-key 7=this_is_bad_too", "")
	assert.Equal(t, map[string]string{}, xto.CustomKVs())
}

func TestValidTs(t *testing.T) {
	xto := xtrace.NewXTraceOptions("ts=12345", "")
	assert.Equal(t, int64(12345), xto.Timestamp())
}

func TestInvalidTs(t *testing.T) {
	xto := xtrace.NewXTraceOptions("ts=invalid", "")
	assert.Equal(t, int64(0), xto.Timestamp())
}

func TestSig(t *testing.T) {
	xto := xtrace.NewXTraceOptions("foo bar baz", "signature123")
	assert.Equal(t, "signature123", xto.Signature())
}

func TestSigWithoutOptions(t *testing.T) {
	xto := xtrace.NewXTraceOptions("", "signature123")
	assert.Equal(t, "signature123", xto.Signature())
}

func TestDocumentedExample1(t *testing.T) {
	xto := xtrace.NewXTraceOptions("trigger-trace;sw-keys=check-id:check-1013,website-id:booking-demo", "")
	assert.True(t, xto.TriggerTrace())
	assert.Empty(t, xto.CustomKVs())
	assert.Equal(t, "check-id:check-1013,website-id:booking-demo", xto.SwKeys())
}

func TestDocumentedExample2(t *testing.T) {
	xto := xtrace.NewXTraceOptions("trigger-trace;custom-key1=value1", "")
	assert.True(t, xto.TriggerTrace())
	assert.Equal(t, map[string]string{"custom-key1": "value1"}, xto.CustomKVs())
	assert.Empty(t, xto.SwKeys())
}

func TestDocumentedExample3(t *testing.T) {
	xto := xtrace.NewXTraceOptions(
		"trigger-trace;sw-keys=check-id:check-1013,website-id:booking-demo;ts=1564432370",
		"5c7c733c727e5038d2cd537630206d072bbfc07c",
	)
	assert.True(t, xto.TriggerTrace())
	assert.Empty(t, xto.CustomKVs())
	assert.Equal(t, "check-id:check-1013,website-id:booking-demo", xto.SwKeys())
	assert.Equal(t, int64(1564432370), xto.Timestamp())
}

func TestStripAllOptions(t *testing.T) {
	xto := xtrace.NewXTraceOptions(
		" trigger-trace ;  custom-something=value; custom-OtherThing = other val ;  sw-keys = 029734wr70:9wqj21,0d9j1   ; ts = 12345 ; foo = bar ",
		"",
	)
	assert.Empty(t, xto.Signature())
	assert.Equal(t, map[string]string{
		"custom-something":  "value",
		"custom-OtherThing": "other val",
	}, xto.CustomKVs())
	assert.Equal(t, "029734wr70:9wqj21,0d9j1", xto.SwKeys())
	assert.True(t, xto.TriggerTrace())
	assert.Equal(t, int64(12345), xto.Timestamp())
}

func TestAllOptionsHandleSequentialSemicolons(t *testing.T) {
	xto := xtrace.NewXTraceOptions(
		";foo=bar;;;custom-something=value_thing;;sw-keys=02973r70:1b2a3;;;;custom-key=val;ts=12345;;;;;;;trigger-trace;;;",
		"",
	)
	assert.Empty(t, xto.Signature())
	assert.Equal(t, map[string]string{
		"custom-something": "value_thing",
		"custom-key":       "val",
	}, xto.CustomKVs())
	assert.Equal(t, "02973r70:1b2a3", xto.SwKeys())
	assert.True(t, xto.TriggerTrace())
	assert.Equal(t, int64(12345), xto.Timestamp())
}

func TestAllOptionsHandleSingleQuotes(t *testing.T) {
	xto := xtrace.NewXTraceOptions(
		"trigger-trace;custom-foo='bar;bar';custom-bar=foo",
		"",
	)
	assert.Empty(t, xto.Signature())
	assert.Equal(t, map[string]string{
		"custom-foo": "'bar",
		"custom-bar": "foo",
	}, xto.CustomKVs())
	assert.Empty(t, xto.SwKeys())
	assert.True(t, xto.TriggerTrace())
	assert.Equal(t, int64(0), xto.Timestamp())
}

func TestAllOptionsHandleMissingValuesAndSemicolons(t *testing.T) {
	xto := xtrace.NewXTraceOptions(
		";trigger-trace;custom-something=value_thing;sw-keys=02973r70:9wqj21,0d9j1;1;2;3;4;5;=custom-key=val?;=",
		"",
	)
	assert.Empty(t, xto.Signature())
	assert.Equal(t, map[string]string{
		"custom-something": "value_thing",
	}, xto.CustomKVs())
	assert.Equal(t, "02973r70:9wqj21,0d9j1", xto.SwKeys())
	assert.True(t, xto.TriggerTrace())
	assert.Equal(t, int64(0), xto.Timestamp())
}
