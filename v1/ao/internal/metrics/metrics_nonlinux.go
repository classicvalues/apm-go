//go:build !linux
// +build !linux

// Copyright (C) 2017 Librato, Inc. All rights reserved.

package metrics

import "github.com/solarwindscloud/swo-golang/v1/ao/internal/bson"

func appendUname(bbuf *bson.Buffer) {}

func addHostMetrics(bbuf *bson.Buffer, index *int) {}
