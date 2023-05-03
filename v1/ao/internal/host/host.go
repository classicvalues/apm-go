// Copyright (C) 2023 SolarWinds Worldwide, LLC. All rights reserved.

// Package host has all the functions to get host metadata needed by traces and
// metrics. It also maintains an update-to-date global HostID object, which is
// refreshed periodically.
package host

import (
	"net"
	"os"
	"sync"

	"github.com/solarwindscloud/solarwinds-apm-go/v1/ao/internal/config"
	"github.com/solarwindscloud/solarwinds-apm-go/v1/ao/internal/log"
)

const (
	REDHAT    = "/etc/redhat-release"
	AMAZON    = "/etc/system-release"
	UBUNTU    = "/etc/lsb-release"
	DEBIAN    = "/etc/debian_version"
	SUSE      = "/etc/os-release"
	SUSE_OLD  = "/etc/SuSE-release"
	SLACKWARE = "/etc/slackware-version"
	GENTOO    = "/etc/gentoo-release"
	OTHER     = "/etc/issue"
)

// logging texts
const (
	stopHostIdObserverByUser = "Host observer is closing per request."
)

var (
	// hostId stores the up-to-date ID info, which is updated periodically
	hostId = newLockedID()

	// exit indicates the ID observer should exit when it's closed
	exit = make(chan struct{})

	// make sure the channel exit is not closed twice, it's effectively immutable.
	exitClosed sync.Once

	// make sure the host observer starts only once
	startOnce sync.Once

	// the cache for initDistro information and its lock
	distro     string
	distroOnce sync.Once

	// the cache for pid, it's only modified/initialized when this package is
	// imported.
	pid = getPid()

	// hostname and its lock
	hostname, _ = os.Hostname()
	hm          sync.RWMutex
)

// CurrentID returns a copyID of the current ID
func CurrentID() ID {
	hostId.waitForReady()
	return hostId.copyID()
}

// BestEffortCurrentID returns the current host ID with the best effort. It
// doesn't wait until the ID is ready. This function is used mainly by getSettings
// where the ID may not be fully initialized immediately after starting up but
// it's acceptable.
func BestEffortCurrentID() ID {
	return hostId.copyID()
}

// PID returns the cached process ID
func PID() int {
	return pid
}

// Start starts the host observer as a standalone goroutine, which will refresh
// the host metadata periodically
func Start() {
	startOnce.Do(func() {
		go observer()
	})
}

// Stop stops the host metadata refreshing goroutine
func Stop() {
	exitClosed.Do(func() {
		close(exit)
		log.Info(stopHostIdObserverByUser)
	})
}

// ConfiguredHostname returns the hostname configured by user
func ConfiguredHostname() string {
	return config.GetHostAlias()
}

// Hostname returns the hostname
func Hostname() string {
	hm.RLock()
	defer hm.RUnlock()
	return hostname
}

// IPAddresses gets the system's IP addresses
func IPAddresses() []string {
	ifaces, err := FilteredIfaces()
	if err != nil {
		return nil
	}

	var addresses []string

	for _, iface := range ifaces {
		// get unicast addresses associated with the current network interface
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				addresses = append(addresses, ipnet.IP.String())
			}
		}
	}

	return addresses
}

// FilteredIfaces returns a list of Interface which contains only interfaces
// required. See https://swicloud.atlassian.net/browse/AO-9021
func FilteredIfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var filtered []net.Interface
	for _, iface := range ifaces {
		// skip over local interface
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// skip over point-to-point interface
		if iface.Flags&net.FlagPointToPoint != 0 {
			continue
		}
		// skip over virtual interface
		if physical := IsPhysicalInterface(iface.Name); !physical {
			continue
		}
		// skip over interfaces without unicast IP addresses
		addrs, err := iface.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}
		filtered = append(filtered, iface)
	}
	return filtered, nil
}

// Distro returns the distro information
func Distro() string {
	distroOnce.Do(func() {
		distro = initDistro()
	})
	return distro
}
