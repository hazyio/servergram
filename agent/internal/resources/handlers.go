package resources

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"

	"hazyio.com/servergram/shared/gen/proto"
)

func getInfo(ctx context.Context) (*proto.GetServerInfoResponse, error) {
	cpuInfo, _ := cpu.InfoWithContext(ctx)
	h, _ := host.InfoWithContext(ctx)
	now := time.Now().UTC()

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
		SystemTime:    now.Format(time.RFC3339),
	}, nil
}

func topProcesses(ctx context.Context, req *proto.TopProcessRequest) (*proto.TopProcessResponse, error) {
	switch req.SortBy {
	case proto.TopProcessSortBy_CPU_USAGE:
		// Implement CPU usage sorting logic here
	case proto.TopProcessSortBy_MEMORY_USAGE:
		// Implement memory usage sorting logic here
	case proto.TopProcessSortBy_DISK_USAGE:
		// Implement disk usage sorting logic here
	}
	// Implement the logic to gather top processes here
	return nil, nil
}
