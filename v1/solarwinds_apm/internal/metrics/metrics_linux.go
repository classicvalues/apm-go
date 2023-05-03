//go:build linux
// +build linux

// Copyright (C) 2023 SolarWinds Worldwide, LLC. All rights reserved.

package metrics

import (
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/bson"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/solarwinds_apm/internal/utils"
)

// gets and appends UnameSysName/UnameVersion to a BSON buffer
// bbuf	the BSON buffer to append the KVs to
func appendUname(bbuf *bson.Buffer) {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err == nil {
		sysname := utils.Byte2String(uname.Sysname[:])
		release := utils.Byte2String(uname.Release[:])
		bbuf.AppendString("UnameSysName", strings.TrimRight(sysname, "\x00"))
		bbuf.AppendString("UnameVersion", strings.TrimRight(release, "\x00"))
	}
}

func addHostMetrics(bbuf *bson.Buffer, index *int) {
	// system load of last minute
	if s := utils.GetStrByKeyword("/proc/loadavg", ""); s != "" {
		load, err := strconv.ParseFloat(strings.Fields(s)[0], 64)
		if err == nil {
			addMetricsValue(bbuf, index, "Load1", load)
		}
	}

	// system total memory
	if s := utils.GetStrByKeyword("/proc/meminfo", "MemTotal"); s != "" {
		memTotal := strings.Fields(s) // MemTotal: 7657668 kB
		if len(memTotal) == 3 {
			if total, err := strconv.Atoi(memTotal[1]); err == nil {
				addMetricsValue(bbuf, index, "TotalRAM", int64(total*1024))
			}
		}
	}

	// free memory
	if s := utils.GetStrByKeyword("/proc/meminfo", "MemFree"); s != "" {
		memFree := strings.Fields(s) // MemFree: 161396 kB
		if len(memFree) == 3 {
			if free, err := strconv.Atoi(memFree[1]); err == nil {
				addMetricsValue(bbuf, index, "FreeRAM", int64(free*1024)) // bytes
			}
		}
	}

	// process memory
	if s := utils.GetStrByKeyword("/proc/self/statm", ""); s != "" {
		processRAM := strings.Fields(s)
		if len(processRAM) != 0 {
			for _, ps := range processRAM {
				if p, err := strconv.Atoi(ps); err == nil {
					addMetricsValue(bbuf, index, "ProcessRAM", p*os.Getpagesize())
					break
				}
			}
		}
	}
}
