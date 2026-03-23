package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/disgoorg/disgo/gateway"
	"github.com/shirou/gopsutil/v4/process"
)

func (b *Bot) startPresenceUpdater(ctx context.Context) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		b.Logger.Error("failed to create process handle for presence", slog.Any("error", err))
		return
	}

	b.updatePresence(ctx, proc)

	ticker := time.NewTicker(b.Cfg.PresenceInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.updatePresence(ctx, proc)
		}
	}
}

func (b *Bot) updatePresence(ctx context.Context, proc *process.Process) {
	cpuPercent, err := proc.CPUPercentWithContext(ctx)
	if err != nil {
		b.Logger.Warn("failed to get CPU percent", slog.Any("error", err))
		return
	}

	memInfo, err := proc.MemoryInfoWithContext(ctx)
	if err != nil {
		b.Logger.Warn("failed to get memory info", slog.Any("error", err))
		return
	}

	ramMB := memInfo.RSS / 1024 / 1024
	status := fmt.Sprintf("RAM: %d MB | CPU: %.1f%%", ramMB, cpuPercent)

	if err := b.Client.SetPresence(ctx, gateway.WithWatchingActivity(status)); err != nil {
		b.Logger.Warn("failed to set presence", slog.Any("error", err))
	}
}
