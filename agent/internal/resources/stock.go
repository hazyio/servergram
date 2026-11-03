package resources

import (
	"context"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"

	"hazyio.com/servergram/shared/gen/proto"
)
func cpuUsage(ctx context.Context) (*proto.ServerCPUUsage, error) {
	percentages, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, err
	}
	physical, _ := cpu.CountsWithContext(ctx, false)
	logical, _ := cpu.CountsWithContext(ctx, true)
	return &proto.ServerCPUUsage{
		Usage:   percentages[0],
		Cores:   int32(physical),
		Threads: int32(logical),
	}, nil
}

func diskUsage(ctx context.Context) (*proto.ServerUsageItem, error) {
	usage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return nil, err
	}
	return &proto.ServerUsageItem{
		Usage:     usage.Used,
		Total:     usage.Total,
		Available: usage.Free,
		Percent:   uint32(usage.UsedPercent),
	}, nil
}

func swapUsage(ctx context.Context) (*proto.ServerUsageItem, error) {
	usage, err := mem.SwapMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}
	if usage.Total == 0 {
		return nil, nil
	}
	return &proto.ServerUsageItem{
		Usage:     usage.Used,
		Total:     usage.Total,
		Available: usage.Free,
		Percent:   uint32(usage.UsedPercent),
	}, nil
}

func memoryUsage() (proto.ServerUsageItem, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return proto.ServerUsageItem{}, err
	}

	return proto.ServerUsageItem{
		Usage:     v.Used,
		Total:     v.Total,
		Available: v.Available,
		Percent:   uint32(v.UsedPercent),
	}, nil
}
