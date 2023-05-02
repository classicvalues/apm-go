// Copyright (C) 2017 Librato, Inc. All rights reserved.

package config

import (
	"sync/atomic"
)

// ReporterOptions defines the options of a reporter. The fields of it
// must be accessed through atomic operators
type ReporterOptions struct {
	// Events flush interval in seconds
	EventFlushInterval int64 `yaml:"EventFlushInterval,omitempty" env:"SWO_EVENTS_FLUSH_INTERVAL" default:"2"`

	// The maximum bytes per RPC request
	MaxReqBytes int64 `yaml:"MaxReqBytes,omitempty" env:"SWO_MAX_REQUEST_BYTES" default:"2048000"`

	// Metrics flush interval in seconds
	MetricFlushInterval int64 `yaml:"MetricFlushInterval,omitempty" default:"30"`

	// GetSettings interval in seconds
	GetSettingsInterval int64 `yaml:"GetSettingsInterval,omitempty" default:"30"`

	// Settings timeout interval in seconds
	SettingsTimeoutInterval int64 `yaml:"SettingsTimeoutInterval,omitempty" default:"10"`

	// Ping interval in seconds
	PingInterval int64 `yaml:"PingInterval,omitempty" default:"20"`

	// Retry backoff initial delay
	RetryDelayInitial int64 `yaml:"RetryDelayInitial,omitempty" default:"500"`

	// Maximum retry delay
	RetryDelayMax int `yaml:"RetryDelayMax,omitempty" default:"60"`

	// Maximum redirect times
	RedirectMax int `yaml:"RedirectMax,omitempty" default:"20"`

	// The threshold of retries before debug printing
	RetryLogThreshold int `yaml:"RetryLogThreshold,omitempty" default:"10"`

	// The maximum retries
	MaxRetries int `yaml:"MaxRetries,omitempty" default:"20"`
}

// SetEventFlushInterval sets the event flush interval to i
func (r *ReporterOptions) SetEventFlushInterval(i int64) {
	atomic.StoreInt64(&r.EventFlushInterval, i)
}

// SetMaxReqBytes sets the maximum bytes of the PRC request body to i
func (r *ReporterOptions) SetMaxReqBytes(i int64) {
	atomic.StoreInt64(&r.MaxReqBytes, i)
}

// GetEventFlushInterval returns the current event flush interval
func (r *ReporterOptions) GetEventFlushInterval() int64 {
	return atomic.LoadInt64(&r.EventFlushInterval)
}

// GetMaxReqBytes returns the maximum RPC request size
func (r *ReporterOptions) GetMaxReqBytes() int64 {
	return atomic.LoadInt64(&r.MaxReqBytes)
}

func (r *ReporterOptions) validate() error {
	// TODO
	return nil
}
