package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	installDir        string
	installTsInSec    int64
	lastRestartInUSec int64
)

func init() {
	installDir = initInstallDir()
	installTsInSec = initInstallTsInSec()
	lastRestartInUSec = initLastRestartInUSec()
}

func initInstallDir() string {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return "unknown"
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "unknown"
	}

	prevPath := string(os.PathSeparator)
	for path != prevPath {
		base := filepath.Base(path)
		if base == "solarwinds_apm" {
			return path
		}
		prevPath = path
		path = filepath.Dir(path)
	}
	return "unknown"
}

func initInstallTsInSec() int64 {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return 0
	}
	fileStat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fileStat.ModTime().Unix()
}

func initLastRestartInUSec() int64 {
	return time.Now().UnixNano() / 1e3
}

func InstallDir() string {
	return installDir
}

func InstallTsInSec() int64 {
	return installTsInSec
}

func LastRestartInUSec() int64 {
	return lastRestartInUSec
}
