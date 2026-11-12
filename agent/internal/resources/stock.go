package resources

import (
	"cmp"
	"context"
	"slices"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"

	"hazyio.com/servergram/agent/internal/common"
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

func memoryUsage() (*proto.ServerUsageItem, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return &proto.ServerUsageItem{}, err
	}

	return &proto.ServerUsageItem{
		Usage:     v.Used,
		Total:     v.Total,
		Available: v.Available,
		Percent:   uint32(v.UsedPercent),
	}, nil
}

// getProcessesForTop retrieves the list of processes with their CPU and memory usage for top process calculations.
func getProcessesForTop(ctx context.Context) (*proto.TopProcessResponse, error) {
	ps, err := process.Processes()
	if err != nil {
		return nil, err
	}
	items := make([]*proto.TopProcessResponse_Item, 0, len(ps))

	for _, p := range ps {

		// every call can fail per-process (exited, permission denied): never abort the whole listing
		name, err := p.Name()
		if err != nil {
			continue
		}
		pr := proto.TopProcessResponse_Item{Pid: p.Pid, Name: name}
		// TODO: This can go over 100% for multiple CPU cores, consider normalizing by the number of cores if needed. This is left as-is for now to reflect raw CPU usage.
		pr.CpuUsage, _ = p.CPUPercentWithContext(ctx)

		memInfo, _ := p.MemoryInfoWithContext(ctx)
		if memInfo != nil {
			pr.MemoryUsage = memInfo.RSS
		}
		items = append(items, &pr)
	}
	return &proto.TopProcessResponse{
		Items: items,
	}, nil
}

func getTopProcessesByCpu(ctx context.Context, limit int) (*proto.TopProcessResponse, error) {
	processes, err := getProcessesForTop(ctx)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(processes.Items, func(a, b *proto.TopProcessResponse_Item) int {
		if c := cmp.Compare(b.CpuUsage, a.CpuUsage); c != 0 { // descending
			return c
		}
		return cmp.Compare(a.Pid, b.Pid) // stable tiebreak
	})
	return &proto.TopProcessResponse{
		Items: processes.Items[:limit],
	}, nil
}
func getTopProcessesByMemory(ctx context.Context, limit int) (*proto.TopProcessResponse, error) {
	processes, err := getProcessesForTop(ctx)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(processes.Items, func(a, b *proto.TopProcessResponse_Item) int {
		if c := cmp.Compare(b.MemoryUsage, a.MemoryUsage); c != 0 { // descending
			return c
		}
		return cmp.Compare(a.Pid, b.Pid) // stable tiebreak
	})
	return &proto.TopProcessResponse{
		Items: processes.Items[:limit],
	}, nil
}

func readProcesses(processCache *ProcessCacheManager, columns []proto.WatchProcessesRequest_WatchProcessesColumns, ctx context.Context) ([]*proto.WatchProcessesResponse_Item, error) {
	ps, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]*proto.WatchProcessesResponse_Item, 0, len(ps))
	runningPid := make([]int32, 0, len(ps))
	for _, p := range ps {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		item := &proto.WatchProcessesResponse_Item{Pid: p.Pid}
		if cachedItem, ok := processCache.Get(p.Pid); ok {
			item.Name = cachedItem.Name
			item.Username = cachedItem.Username
			item.Command = cachedItem.Command
			item.Ppid = cachedItem.Ppid
		} else {
			name := common.UnWrapLogErr(p.NameWithContext(ctx))
			if name == nil {
				continue
			}
			item.Name = *name
			ppid := common.UnWrapLogErr(p.PpidWithContext(ctx))
			if ppid == nil {
				continue
			}
			item.Ppid = *ppid

			if slices.Contains(columns, proto.WatchProcessesRequest_WATCH_PROCESSES_COLUMNS_USER) {
				userName := common.UnWrapLogErr(p.UsernameWithContext(ctx))
				if userName == nil {
					continue
				}
				item.Username = *userName

			}

			if slices.Contains(columns, proto.WatchProcessesRequest_WATCH_PROCESSES_COLUMNS_COMMAND) {
				cmd := common.UnWrapLogErr(p.CmdlineWithContext(ctx))
				if cmd == nil {
					continue
				}
				item.Command = *cmd
			}

			// add item to cache
			processCache.Set(p.Pid, ProcessCacheData{
				Name:     item.Name,
				Username: item.Username,
				Command:  item.Command,
				Ppid:     item.Ppid,
			})
		}
		if slices.Contains(columns, proto.WatchProcessesRequest_WATCH_PROCESSES_COLUMNS_STATE) {
			st := common.UnWrapLogErr(p.StatusWithContext(ctx))
			if st == nil || len(*st) == 0 {
				continue
			}
			status := (*st)[0]
			switch status {
			case "S":
				item.Status = proto.WatchProcessesResponse_WATCH_PROCESSES_STATUS_SLEEPING
			case "T":
				item.Status = proto.WatchProcessesResponse_WATCH_PROCESSES_STATUS_STOPPED
			case "Z":
				item.Status = proto.WatchProcessesResponse_WATCH_PROCESSES_STATUS_ZOMBIE
			case "I":
				item.Status = proto.WatchProcessesResponse_WATCH_PROCESSES_STATUS_IDLE
			case "D":
				item.Status = proto.WatchProcessesResponse_WATCH_PROCESSES_STATUS_WAITING
			default:
				item.Status = proto.WatchProcessesResponse_WATCH_PROCESSES_STATUS_OTHER
			}
		}

		if slices.Contains(columns, proto.WatchProcessesRequest_WATCH_PROCESSES_COLUMNS_MEMORY) {
			mem, _ := p.MemoryInfoWithContext(ctx)
			if mem != nil {
				item.MemoryUsage = mem.RSS
			}
		}
		items = append(items, item)
		runningPid = append(runningPid, p.Pid)

	}
	processCache.RemoveAllBut(runningPid)

	return items, nil
}
