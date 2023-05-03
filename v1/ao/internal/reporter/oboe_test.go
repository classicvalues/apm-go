// Copyright (C) 2016 Librato, Inc. All rights reserved.

package reporter

import (
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/solarwindscloud/swo-golang/v1/ao/internal/config"
	"github.com/solarwindscloud/swo-golang/v1/ao/internal/utils"
	"github.com/stretchr/testify/require"
	mbson "gopkg.in/mgo.v2/bson"

	g "github.com/solarwindscloud/swo-golang/v1/ao/internal/graphtest"
	"github.com/stretchr/testify/assert"
)

func newTokenBucket(ratePerSec, size float64) *tokenBucket {
	return &tokenBucket{ratePerSec: ratePerSec, capacity: size, available: size, last: time.Now()}
}

func TestInitMessage(t *testing.T) {
	r := SetTestReporter()

	sendInitMessage()
	r.Close(1)
	assertInitMessage(t, r.EventBufs)
}
func assertInitMessage(t *testing.T, bufs [][]byte) {
	baseline, err := time.Parse(time.RFC3339, "2019-10-02T10:04:05Z")
	require.Nil(t, err)

	g.AssertGraph(t, bufs, 1, g.AssertNodeMap{
		{"go", "single"}: {Edges: g.Edges{}, Callback: func(n g.Node) {
			assert.Equal(t, 1, n.Map["__Init"])
			assert.Equal(t, utils.Version(), n.Map["Go.SolarWindsAPM.Version"])
			assert.NotEmpty(t, n.Map["Go.Version"])
			assert.True(t, strings.HasSuffix(n.Map["Go.InstallDirectory"].(string), "swo-golang/v1/ao"))
			assert.Less(t, baseline.Unix(), n.Map["Go.InstallTimestamp"])
			assert.Less(t, baseline.UnixNano()/1e3, n.Map["Go.LastRestart"])
		}},
	})
}

func TestInitMessageUDP(t *testing.T) {
	assertUDPMode(t)

	var bufs [][]byte
	done := startTestUDPListener(t, &bufs, 2)
	sendInitMessage()
	<-done
	assertInitMessage(t, bufs)
}

func TestTokenBucket(t *testing.T) {
	b := newTokenBucket(5, 2)
	c := b

	consumers := 5
	iters := 100
	sendRate := 30 // test request rate of 30 per second
	sleepInterval := time.Second / time.Duration(sendRate)
	var wg sync.WaitGroup
	wg.Add(consumers)
	var dropped, allowed int64
	for j := 0; j < consumers; j++ {
		go func(id int) {
			perConsumerRate := newTokenBucket(15, 1)
			for i := 0; i < iters; i++ {
				sampled := perConsumerRate.consume(1)
				ok := b.count(sampled, false, true)
				if ok {
					// t.Logf("### OK   id %02d now %v last %v tokens %v", id, time.Now(), b.last, b.available)
					atomic.AddInt64(&allowed, 1)
				} else {
					// t.Logf("--- DROP id %02d now %v last %v tokens %v", id, time.Now(), b.last, b.available)
					atomic.AddInt64(&dropped, 1)
				}
				time.Sleep(sleepInterval)
			}
			wg.Done()
		}(j)
		time.Sleep(sleepInterval / time.Duration(consumers))
	}
	wg.Wait()
	t.Logf("TB iters %d allowed %v dropped %v limited %v", iters, allowed, dropped, c.Limited())
	t.Logf("%+v", c.RateCounts)
	assert.True(t, (allowed == 20 && dropped == 480 && c.Limited() == 230 && c.Traced() == 20) ||
		(allowed == 19 && dropped == 481 && c.Limited() == 231 && c.Traced() == 19) ||
		(allowed == 18 && dropped == 482 && c.Limited() == 232 && c.Traced() == 18))
	assert.Equal(t, int64(500), c.Requested())
	assert.Equal(t, int64(500), c.Sampled())
	assert.Equal(t, int64(0), c.Through())
}

func TestTokenBucketTime(t *testing.T) {
	b := newTokenBucket(5, 2)
	b.consume(1)
	assert.EqualValues(t, 1, b.avail()) // 1 available
	b.last = b.last.Add(time.Second)    // simulate time going backwards
	b.update(time.Now())
	assert.EqualValues(t, 1, b.avail()) // no new tokens added
	assert.True(t, b.consume(1))        // consume available token
	assert.False(t, b.consume(1))       // out of tokens
	assert.True(t, time.Now().After(b.last))
	time.Sleep(200 * time.Millisecond)
	assert.True(t, b.consume(1)) // another token available

	b = newTokenBucket(5, 5)
	assert.EqualValues(t, 5, b.avail())
	b.setRateCap(5, 3)
	assert.EqualValues(t, 3, b.available)
}

func testLayerCount(count int64) interface{} {
	return mbson.D{mbson.DocElem{Name: testLayer, Value: count}}
}

func callShouldTraceRequest(total int, isTraced bool) (traced int) {
	for i := 0; i < total; i++ {
		if ok, _, _, _ := shouldTraceRequest(testLayer, isTraced); ok {
			traced++
		}
	}
	return traced
}

func TestSamplingRate(t *testing.T) {
	r := SetTestReporter(TestReporterDisableDefaultSetting(true))

	// set 2.5% sampling rate
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		25000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))

	total := 100000
	traced := callShouldTraceRequest(total, false)

	// make sure we're within 20% of our expected rate over 1,000,000 trials
	assert.InDelta(t, 2.5, float64(traced)*100/float64(total), 0.2)

	setting, ok := getSetting("")
	assert.True(t, ok)
	c := setting.bucket

	assert.EqualValues(t, c.Requested(), total)
	assert.EqualValues(t, c.Through(), 0)
	assert.EqualValues(t, c.Traced(), traced)
	assert.EqualValues(t, c.Sampled(), total)
	assert.EqualValues(t, c.Limited(), 0)

	r.Close(0)
	// XXX assert bufs
}

func TestSampleNoValidSettings(t *testing.T) {
	r := SetTestReporter(TestReporterDisableDefaultSetting(true))

	total := 1

	// var buf bytes.Buffer
	// log.SetOutput(&buf)
	traced := callShouldTraceRequest(total, false)
	// log.SetOutput(os.Stderr)
	// assert.Contains(t, buf.String(), "Sampling disabled for go_test until valid settings are retrieved")
	assert.EqualValues(t, 0, traced)

	r.Close(0)
}

func TestSampleRateBoundaries(t *testing.T) {
	r := SetTestReporter()

	_, rate, _, _ := shouldTraceRequest(testLayer, false)
	assert.Equal(t, 1000000, rate)

	// check that max value doesn't go above 1000000
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000001, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))

	_, rate, _, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, 1000000, rate)

	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		0, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))

	_, rate, _, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, 0, rate)

	// check that min value doesn't go below 0
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		-1, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))

	_, rate, _, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, 0, rate)

	r.Close(0)
}

func TestSampleSource(t *testing.T) {
	r := SetTestReporter()

	_, _, source, _ := shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_DEFAULT, source)

	resetSettings()
	_, _, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_NONE, source)

	// we're currently only looking up default settings, so this should return NONE sample source
	updateSetting(int32(TYPE_LAYER), testLayer,
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	_, _, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_NONE, source)

	// as soon as we add the default settings back, we get a valid sample source
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	_, _, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_DEFAULT, source)

	r.Close(0)
}

func TestSampleFlags(t *testing.T) {
	r := SetTestReporter(TestReporterDisableDefaultSetting(true))

	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte(""),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	ok, _, _, _ := shouldTraceRequest(testLayer, false)
	assert.False(t, ok)

	setting, ok := getSetting("")
	require.True(t, ok)
	c := setting.bucket
	assert.EqualValues(t, 0, c.Through())
	ok, _, _, _ = shouldTraceRequest(testLayer, true)
	assert.False(t, ok)
	assert.EqualValues(t, 0, c.Through())

	resetSettings()
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START"),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	ok, _, _, _ = shouldTraceRequest(testLayer, false)
	assert.True(t, ok)

	setting, ok = getSetting("")
	require.True(t, ok)
	c = setting.bucket
	assert.EqualValues(t, 0, c.Through())
	ok, _, _, _ = shouldTraceRequest(testLayer, true)
	assert.False(t, ok)
	assert.EqualValues(t, 0, c.Through())

	// Transaction filtering
	urls.loadConfig([]config.TransactionFilter{
		{Type: "url", RegEx: `user\d{3}`, Tracing: config.DisabledTracingMode},
		{Type: "url", Extensions: []string{".png", ".jpg"}, Tracing: config.DisabledTracingMode},
	})
	decision := shouldTraceRequestWithURL(testLayer, false, "http://test.com/user123", ModeTriggerTraceNotPresent)
	assert.False(t, decision.trace)

	resetSettings()

	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_THROUGH_ALWAYS"),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	ok, _, _, _ = shouldTraceRequest(testLayer, false)
	assert.False(t, ok)

	setting, ok = getSetting("")
	require.True(t, ok)
	c = setting.bucket
	assert.EqualValues(t, 0, c.Through())
	ok, _, _, _ = shouldTraceRequest(testLayer, true)
	assert.True(t, ok)
	assert.EqualValues(t, 1, c.Through())

	resetSettings()

	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_THROUGH"),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	ok, _, _, _ = shouldTraceRequest(testLayer, false)
	assert.False(t, ok)

	setting, ok = getSetting("")
	require.True(t, ok)
	c = setting.bucket
	assert.EqualValues(t, 0, c.Through())
	ok, _, _, _ = shouldTraceRequest(testLayer, true)
	assert.True(t, ok)
	assert.EqualValues(t, 1, c.Through())

	r.Close(0)
}

func TestSampleTokenBucket(t *testing.T) {
	r := SetTestReporter()
	setting, ok := getSetting("")
	require.True(t, ok)
	c := setting.bucket

	traced := callShouldTraceRequest(1, false)
	assert.EqualValues(t, 1, traced)
	assert.EqualValues(t, 1, c.Traced())
	assert.EqualValues(t, 1, c.Requested())
	assert.EqualValues(t, 0, c.Limited())

	resetSettings()

	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START"),
		1000000, 120, argsToMap(0, 0, 0, 0, 0, 0, -1, -1, []byte("")))
	traced = callShouldTraceRequest(1, false)
	assert.EqualValues(t, 0, traced)

	setting, ok = getSetting("")
	require.True(t, ok)
	c = setting.bucket
	assert.EqualValues(t, 0, c.Traced())
	assert.EqualValues(t, 1, c.Requested())
	assert.EqualValues(t, 1, c.Limited())

	resetSettings()

	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000000, 120, argsToMap(16, 8, 16, 8, 16, 8, -1, -1, []byte("")))
	traced = callShouldTraceRequest(50, false)
	assert.EqualValues(t, 16, traced)

	setting, ok = getSetting("")
	require.True(t, ok)
	c = setting.bucket
	assert.EqualValues(t, 16, c.Traced())
	assert.EqualValues(t, 50, c.Requested())
	assert.EqualValues(t, 34, c.Limited())
	FlushRateCounts()

	time.Sleep(1 * time.Second)

	traced = callShouldTraceRequest(50, false)
	assert.EqualValues(t, 8, traced)
	assert.EqualValues(t, 8, c.Traced())
	assert.EqualValues(t, 50, c.Requested())
	assert.EqualValues(t, 42, c.Limited())

	r.Close(0)
}

// func TestMetrics(t *testing.T) {
// 	// error sending metrics message: no reporting
// 	r := SetTestReporter()
//
// 	randReader = &errorReader{failOn: map[int]bool{0: true}}
// 	sendMetricsMessage(r)
// 	time.Sleep(100 * time.Millisecond)
// 	r.Close(0)
// 	assert.Len(t, r.EventBufs, 0)
//
// 	r = SetTestReporter()
// 	randReader = &errorReader{failOn: map[int]bool{2: true}}
// 	sendMetricsMessage(r)
// 	time.Sleep(100 * time.Millisecond)
// 	r.Close(0)
// 	assert.Len(t, r.EventBufs, 0)
//
// 	randReader = rand.Reader // set back to normal
// }

// func assertGetNextInterval(t *testing.T, nowTime, expectedDur string) {
// 	t0, err := time.Parse(time.RFC3339Nano, nowTime)
// 	assert.NoError(t, err)
// 	d0 := getNextInterval(t0)
// 	d0e, err := time.ParseDuration(expectedDur)
// 	assert.NoError(t, err)
// 	assert.Equal(t, d0e, d0)
// 	assert.Equal(t, 0, t0.Add(d0).Second()%counterIntervalSecs)
// }
//
// func TestGetNextInterval(t *testing.T) {
// 	assertGetNextInterval(t, "2016-01-02T15:04:05.888-04:00", "24.112s")
// 	assertGetNextInterval(t, "2016-01-02T15:04:35.888-04:00", "24.112s")
// 	assertGetNextInterval(t, "2016-01-02T15:04:00.00-04:00", "30s")
// 	assertGetNextInterval(t, "2016-08-15T23:31:30.00-00:00", "30s")
// 	assertGetNextInterval(t, "2016-01-02T15:04:59.999999999-04:00", "1ns")
// 	assertGetNextInterval(t, "2016-01-07T15:04:29.999999999-00:00", "1ns")
// }

// func TestSendMetrics(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping metrics periodic sender test")
// 	}
// 	// full periodic sender test: wait for next interval & report
// 	r := SetTestReporter(TestReporterTimeout(time.Duration(30) * time.Second))
// 	disableMetrics = false
// 	go sendInitMessage()
// 	d0 := getNextInterval(time.Now()) + time.Second
// 	fmt.Printf("[%v] TestSendMetrics Sleeping for %v\n", time.Now(), d0)
// 	time.Sleep(d0)
// 	fmt.Printf("[%v] TestSendMetrics Closing\n", time.Now())
// 	r.Close(4)
// 	g.AssertGraph(t, r.EventBufs, 4, g.AssertNodeMap{
// 		{"go", "entry"}: {Edges: g.Edges{}, Callback: func(n g.Node) {
// 			assert.Equal(t, 1, n.Map["__Init"])
// 			assert.Equal(t, initVersion, n.Map["Go.Oboe.Version"])
// 			assert.NotEmpty(t, n.Map["Oboe.Version"])
// 			assert.NotEmpty(t, n.Map["Go.Version"])
// 		}},
// 		{"go", "exit"}: {Edges: g.Edges{{"go", "entry"}}},
// 		{metricsLayerName, "entry"}: {Edges: g.Edges{}, Callback: func(n g.Node) {
// 			assert.Equal(t, "go", n.Map["ProcessName"])
// 			assert.IsType(t, int64(0), n.Map["JMX.type=threadcount,name=NumGoroutine"])
// 			assert.IsType(t, int64(0), n.Map["JMX.Memory:MemStats.Alloc"])
// 			assert.True(t, len(n.Map) > 10)
// 		}},
// 		{metricsLayerName, "exit"}: {Edges: g.Edges{{metricsLayerName, "entry"}}},
// 	})
// 	stopMetrics <- struct{}{}
// 	disableMetrics = true
// }

func TestCheckSettingsTimeout(t *testing.T) {
	sc := &oboeSettingsCfg{
		settings: make(map[oboeSettingKey]*oboeSettings),
	}
	k1 := oboeSettingKey{
		sType: TYPE_DEFAULT,
		layer: "expired",
	}
	sc.settings[k1] = &oboeSettings{
		timestamp: time.Now().Add(-time.Second * 2),
		ttl:       1,
	}

	k2 := oboeSettingKey{
		sType: TYPE_DEFAULT,
		layer: "alive",
	}
	sc.settings[k2] = &oboeSettings{
		timestamp: time.Now(),
		ttl:       2,
	}
	sc.checkSettingsTimeout()
	assert.Contains(t, sc.settings, k2, k2.layer)
	assert.NotContains(t, sc.settings, k1, k1.layer)
}

func TestMergeRemoteSettingWithLocalConfig(t *testing.T) {
	// No remote setting
	resetSettings()
	trace, rate, source, _ := shouldTraceRequest(testLayer, false)
	assert.False(t, trace)
	assert.Equal(t, source, SAMPLE_SOURCE_NONE)
	assert.Equal(t, rate, 0)

	resetSettings()
	// Remote setting has the override flag && local config has lower rate
	_ = os.Setenv("SWO_TRACING_MODE", "enabled")
	_ = os.Setenv("SWO_SAMPLE_RATE", "10000")
	_ = config.Load()
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("OVERRIDE,SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	trace, rate, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_FILE, source)
	assert.Equal(t, 10000, rate)

	// Remote setting has the override flag && local config has higher rate
	_ = os.Setenv("SWO_TRACING_MODE", "enabled")
	_ = os.Setenv("SWO_SAMPLE_RATE", "10000")
	_ = config.Load()
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("OVERRIDE,SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	trace, rate, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_DEFAULT, source)
	assert.Equal(t, 1000, rate)

	// Remote setting doesn't have the override flag && local config has lower rate
	_ = os.Setenv("SWO_TRACING_MODE", "enabled")
	_ = os.Setenv("SWO_SAMPLE_RATE", "10000")
	_ = config.Load()
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	trace, rate, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_FILE, source)
	assert.Equal(t, 10000, rate)
	// Remote setting doesn't have the override flag && local config has higher rate
	_ = os.Setenv("SWO_TRACING_MODE", "enabled")
	_ = os.Setenv("SWO_SAMPLE_RATE", "10000")
	_ = config.Load()
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	trace, rate, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_FILE, source)
	assert.Equal(t, 10000, rate)
	// Remote setting has the override flag && no local config
	_ = os.Unsetenv("SWO_TRACING_MODE")
	_ = os.Unsetenv("SWO_SAMPLE_RATE")
	_ = config.Load()
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("OVERRIDE,SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		10000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	trace, rate, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_DEFAULT, source)
	assert.Equal(t, 10000, rate)
	// Remote setting doesn't have the override flag && no local config
	_ = os.Unsetenv("SWO_TRACING_MODE")
	_ = os.Unsetenv("SWO_SAMPLE_RATE")
	_ = config.Load()
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		10000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	trace, rate, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_DEFAULT, source)
	assert.Equal(t, 10000, rate)
	// Remote setting has the override flag && local tracing mode = DISABLED
	_ = os.Setenv("SWO_TRACING_MODE", "disabled")
	_ = os.Setenv("SWO_SAMPLE_RATE", "10000")
	_ = config.Load()
	updateSetting(int32(TYPE_DEFAULT), "",
		[]byte("OVERRIDE,SAMPLE_START,SAMPLE_THROUGH_ALWAYS"),
		1000000, 120, argsToMap(1000000, 1000000, 1000000, 1000000, 1000000, 1000000, -1, -1, []byte("")))
	trace, rate, source, _ = shouldTraceRequest(testLayer, false)
	assert.Equal(t, SAMPLE_SOURCE_FILE, source)
	assert.Equal(t, 10000, rate)
}

func TestAdjustSampleRate(t *testing.T) {
	assert.Equal(t, maxSamplingRate, adjustSampleRate(maxSamplingRate+1))
	assert.Equal(t, 0, adjustSampleRate(-1))
	assert.Equal(t, maxSamplingRate-1, adjustSampleRate(maxSamplingRate-1))
}

func TestBytesToFloat64(t *testing.T) {
	ret, err := bytesToFloat64([]byte{1, 1, 1, 1})
	assert.EqualValues(t, -1, ret)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid length")
}

func TestBytesToInt32(t *testing.T) {
	ret, err := bytesToInt32([]byte{1, 1, 1})
	assert.EqualValues(t, -1, ret)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid length")
}

func TestParseInt32(t *testing.T) {
	args := map[string][]byte{
		"key": {0, 1, 0, 0},
	}
	assert.EqualValues(t, -100, parseInt32(args, "invalidKey", -100))
	assert.EqualValues(t, 256, parseInt32(args, "key", -100))

	args = map[string][]byte{
		"key": {255, 255, 255, 255},
	}
	assert.EqualValues(t, -100, parseInt32(args, "key", -100))
}

func TestParseFloat64(t *testing.T) {
	args := map[string][]byte{
		"key": {0, 1, 0, 0},
	}
	assert.EqualValues(t, 1.01, parseFloat64(args, "key", 1.01))
}
