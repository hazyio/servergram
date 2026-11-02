package resources

import (
	"context"
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"

	// "github.com/shirou/gopsutil/v4/host"

	"hazyio.com/servergram/shared/gen/proto"
)

func getInfo(ctx context.Context) (*proto.GetServerInfoResponse, error) {
	cpuInfo, _ := cpu.InfoWithContext(ctx)
	h, _ := host.InfoWithContext(ctx)
	// h.PlatformVersion
	// Implement the logic to gather server info here
	return &proto.GetServerInfoResponse{
		CpuModel:      cpuInfo[0].ModelName,
		Os:            h.OS,
		Architecture:  h.KernelArch,
		KernelVersion: h.KernelVersion,
		Hostname:      h.Hostname,
		Uptime:        int64(h.Uptime),
		Platform:      h.Platform,
	}, nil
}

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
		Usage: usage.Used,
		Total: usage.Total,
		Available: usage.Free,
		Percent: uint32(usage.UsedPercent),
	}, nil
}

func memoryUsage() (proto.ServerUsageItem, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return proto.ServerUsageItem{}, err
	}
	log.Println(v.Total, v.Used, v.Available, v.UsedPercent)

	return proto.ServerUsageItem{
		Usage:     v.Used,
		Total:     v.Total,
		Available: v.Available,
		Percent:   uint32(v.UsedPercent),
	}, nil
}
