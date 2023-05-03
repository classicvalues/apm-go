// Copyright (C) 2023 SolarWinds Worldwide, LLC. All rights reserved.

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReporterOptions(t *testing.T) {
	r := &ReporterOptions{}

	r.SetEventFlushInterval(20)
	assert.Equal(t, r.GetEventFlushInterval(), int64(20))

	r.SetMaxReqBytes(2000)
	assert.Equal(t, r.GetMaxReqBytes(), int64(2000))

	assert.Nil(t, r.validate())
}
