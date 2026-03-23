// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package config

import "time"

// Operational defaults — hardcoded values that do not need external configuration.
const (
	DefaultVolume             = 50
	DefaultAutoLeaveTimeout   = 180 * time.Second
	DefaultPresenceInterval   = 30 * time.Second
	DefaultLavalinkTimeout    = 2 * time.Second
	DefaultLavalinkLoadTimeout = 10 * time.Second
	DefaultHTTPClientTimeout  = 10 * time.Second
	DefaultPanelPowerTimeout  = 15 * time.Second
	DefaultRSSPollInterval    = 300 * time.Second
	DefaultRSSFeedTimeout     = 30 * time.Second
	DefaultShortenTimeout     = 10 * time.Second
	DefaultScanTimeout        = 15 * time.Second
)
