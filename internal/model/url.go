// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

// VTResult holds VirusTotal scan results.
type VTResult struct {
	Harmless   int
	Malicious  int
	Suspicious int
	Undetected int
	Timeout    int
}
