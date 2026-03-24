// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import "errors"

var (
	ErrPanelUnauthorized = errors.New("unauthorized (invalid API key)")
	ErrPanelNotFound     = errors.New("server not found")
	ErrPanelRateLimited  = errors.New("rate limited")
)

// Server represents a game server from the Pelican API.
type Server struct {
	Identifier  string
	Name        string
	Description string
	Status      string // "running", "starting", "stopping", "offline"
	Node        string
	Limits      ServerLimits
	IsSuspended bool
}

// ServerLimits holds resource limits for a server.
type ServerLimits struct {
	Memory int // MB
	Disk   int // MB
	CPU    int // percent (e.g. 400 = 4 cores)
}

// Resources holds live resource usage for a server.
type Resources struct {
	CurrentState   string
	MemoryBytes    int64
	CPUAbsolute    float64
	DiskBytes      int64
	NetworkRxBytes int64
	NetworkTxBytes int64
	Uptime         int64 // milliseconds
}
