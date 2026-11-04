package resources

import (
	"context"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"

	"hazyio.com/servergram/agent/internal/common"
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
func getResourcesUsage(req *proto.GetServerResourcesRequest, stream proto.ServerResources_GetServer) error {
	iv := max(time.Duration(req.IntervalMs)*time.Millisecond,
		// clamp: don't let a client ask for 1ms
		time.Second)
	t := time.NewTicker(iv)
	defer t.Stop()
	ctx := stream.Context()

	workFunc := func(stream proto.ServerResources_GetServer) error {
		cpuUsage, err := cpuUsage(ctx)
		if err != nil {
			log.Println("Failed to get CPU usage", "error", err)

		}
		run := common.RunAll(func() (*proto.ServerUsageItem, error) {
			return memoryUsage()
		}, func() (*proto.ServerUsageItem, error) {
			return diskUsage(ctx)
		}, func() (*proto.ServerUsageItem, error) {
			return swapUsage(ctx)
		})

		stream.Send(&proto.GetServerResourcesResponse{
			Cpu:    cpuUsage,
			Memory: run[0].Val,
			Disk:   run[1].Val,
			Swap:   run[2].Val,
		})
		return nil
	}
	// Do the first work immediately, then start the ticker for subsequent work
	workFunc(stream)
	for {
		select {
		case <-stream.Context().Done():
			return nil // client cancelled
		case <-t.C:
			if err := workFunc(stream); err != nil {
				return err
			}
		}
	}
}
func topProcesses(req *proto.TopProcessRequest, stream proto.ServerResources_TopProcessesServer) error {
	iv := max(time.Duration(req.IntervalMs)*time.Millisecond,
		// clamp: don't let a client ask for 1ms
		time.Second)
	t := time.NewTicker(iv)
	defer t.Stop()
	ctx := stream.Context()
	workFunc := func(stream proto.ServerResources_TopProcessesServer) error {
		switch req.SortBy {
		case proto.TopProcessSortBy_CPU_USAGE:
			data, err := getTopProcessesByCpu(ctx, int(req.Limit))
			if err != nil {
				return err
			}
			stream.Send(data)
		case proto.TopProcessSortBy_MEMORY_USAGE:
			data, err := getTopProcessesByMemory(ctx, int(req.Limit))
			if err != nil {
				return err
			}
			stream.Send(data)
		}
		return nil
	}
	// Do the first work immediately, then start the ticker for subsequent work
	workFunc(stream)
	for {
		select {
		case <-stream.Context().Done():
			return nil // client cancelled
		case <-t.C:
			if err := workFunc(stream); err != nil {
				return err
			}
		}
	}

}
