package model

// VTResult holds VirusTotal scan results.
type VTResult struct {
	Harmless   int
	Malicious  int
	Suspicious int
	Undetected int
	Timeout    int
}
