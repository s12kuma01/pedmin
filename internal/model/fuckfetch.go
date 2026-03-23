package model

import "time"

// SystemInfo holds system information for neofetch-style output.
type SystemInfo struct {
	// OS
	OS, Platform, KernelVersion string
	Uptime                      time.Duration
	// CPU
	CPUModel   string
	CPUCores   int // physical
	CPUThreads int // logical
	CPUUsage   float64
	// Memory
	MemTotal, MemUsed, MemAvailable uint64
	MemUsage                        float64
	// Swap
	SwapTotal, SwapUsed uint64
	SwapUsage           float64
	// Disk
	DiskTotal, DiskUsed, DiskFree uint64
	DiskUsage                     float64
	// Network
	NetBytesSent, NetBytesRecv uint64
	// GPU / NPU (optional)
	GPUInfo string
	NPUInfo string
}
