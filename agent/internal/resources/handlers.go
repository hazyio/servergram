package resources

import (
	"context"
	"log"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"hazyio.com/servergram/agent/internal/common"
	"hazyio.com/servergram/shared/gen/proto"
)

func getInfo(ctx context.Context) (*proto.GetServerInfoResponse, error) {
	cpuInfo, _ := cpu.InfoWithContext(ctx)
	h, _ := host.InfoWithContext(ctx)
	now := time.Now().UTC()
	if len(cpuInfo) == 0 {
		return nil, status.Error(codes.Internal, "no CPU info available")
	}
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
func getResourcesUsage(req *proto.GetServerResourcesRequest, stream proto.ResourcesService_GetServer) error {
	iv := max(time.Duration(req.IntervalMs)*time.Millisecond,
		// clamp: don't let a client ask for 1ms
		time.Second)
	t := time.NewTicker(iv)
	defer t.Stop()
	ctx := stream.Context()

	workFunc := func(stream proto.ResourcesService_GetServer) error {
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
func topProcesses(req *proto.TopProcessRequest, stream proto.ResourcesService_TopProcessesServer) error {
	iv := max(time.Duration(req.IntervalMs)*time.Millisecond,
		// clamp: don't let a client ask for 1ms
		time.Second)
	t := time.NewTicker(iv)
	defer t.Stop()
	ctx := stream.Context()
	workFunc := func(stream proto.ResourcesService_TopProcessesServer) error {
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

func watchProcesses(req *proto.WatchProcessesRequest, stream proto.ResourcesService_WatchProcessesServer) error {
	computeDuration := max(time.Duration(req.GetIntervalMs())*time.Millisecond,
		time.Second) // clamp: don't let a client ask for 1ms

	t := time.NewTicker(computeDuration)
	defer t.Stop()
	ctx := stream.Context()
	var running atomic.Bool
	processCache := NewProcessCacheManager()
	columns := req.GetColumns()
	if len(columns) == 0 {
		columns = make([]proto.WatchProcessesRequest_WatchProcessesColumns, 0)
	}
	workFunc := func(stream proto.ResourcesService_WatchProcessesServer) error {
		if running.Load() {
			return nil
		}
		running.Store(true)
		defer running.Store(false)
		timeOutCtx, cancel := context.WithTimeout(ctx, 2*computeDuration) // last for twice the duration
		defer cancel()

		items, err := readProcesses(processCache, columns, timeOutCtx)
		if err != nil {
			// do not close the stream on timeout or errors
			log.Println("Failed to read processes", "error", err)
			return nil
		}
		stream.Send(&proto.WatchProcessesResponse{
			Items: items,
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
func triggerSignal(pid int32, signal syscall.Signal) error {
	processId := int(pid)
	if processId <= 1 {
		return status.Error(codes.InvalidArgument, "invalid pid")
	}
	if err := syscall.Kill(processId, signal); err != nil {
		switch err {
		case syscall.ESRCH:
			return status.Error(codes.NotFound, "no such process")
		case syscall.EPERM:
			return status.Error(codes.PermissionDenied, "not permitted to signal process")
		}
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func stopProcess(req *proto.StopProcessRequest) (*proto.StopProcessResponse, error) {

	if err := triggerSignal(req.Pid, syscall.SIGTERM); err != nil {
		return nil, err
	}

	return &proto.StopProcessResponse{Success: true}, nil
}

func forceStopProcess(req *proto.ForceStopProcessRequest) (*proto.ForceStopProcessResponse, error) {
	if err := triggerSignal(req.Pid, syscall.SIGKILL); err != nil {
		return nil, err
	}

	return &proto.ForceStopProcessResponse{Success: true}, nil
}

func sendSignal(req *proto.SignalProcessRequest) (*proto.SignalProcessResponse, error) {
	var signalMap = map[proto.SignalProcessRequest_SignalProcessRequestSignals]syscall.Signal{
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_HUP:   syscall.SIGHUP,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_INT:   syscall.SIGINT,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_KILL:  syscall.SIGKILL,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_TERM:  syscall.SIGTERM,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_QUIT:  syscall.SIGQUIT,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_USR1:  syscall.SIGUSR1,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_USR2:  syscall.SIGUSR2,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_STOP:  syscall.SIGSTOP,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_CONT:  syscall.SIGCONT,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_TSTP:  syscall.SIGTSTP,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_ALRM:  syscall.SIGALRM,
		proto.SignalProcessRequest_SIGNAL_PROCESS_REQUEST_SIGNALS_ABORT: syscall.SIGABRT,
	}
	sig, ok := signalMap[req.Signal]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "unsupported signal")
	}
	if err := triggerSignal(req.Pid, sig); err != nil {
		return nil, err
	}
	return &proto.SignalProcessResponse{Success: true}, nil
}
