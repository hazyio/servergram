import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import type { components } from "#lib/types/request-schema.d.ts";
import { formatBytes } from "#lib/utils.ts";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";

export type ProcessItem = {
  pid: number;
  ppid: number;
  name: string;
  cpu_usage: number;
  memory_usage: number;
  command: string;
  username: string;
  status: number;
  uptime: number;
};

type RawProcess = components["schemas"]["WatchProcessesResponse_Item"];

export function toProcessItem(process: RawProcess): ProcessItem {
  return {
    pid: process.pid ?? 0,
    ppid: process.ppid ?? 0,
    name: process.name ?? "",
    cpu_usage: process.cpu_usage ?? 0,
    memory_usage: process.memory_usage ?? 0,
    command: process.command ?? "",
    username: process.username ?? "",
    status: process.status ?? 0,
    uptime: process.uptime ?? 0,
  };
}

export function formatProcessStatus(status: number): string {
  switch (status) {
    case 1:
      return "Sleeping";
    case 2:
      return "Stopped";
    case 3:
      return "Zombie";
    case 4:
      return "Idle";
    case 5:
      return "Waiting";
    case 6:
      return "Running";
    case 7:
      return "Other";
    default:
      return "Unknown";
  }
}

export function formatUptime(seconds: number): string {
  if (!Number.isFinite(seconds) || seconds <= 0) return "-";

  const days = Math.floor(seconds / 86_400);
  const hours = Math.floor((seconds % 86_400) / 3_600);
  const minutes = Math.floor((seconds % 3_600) / 60);
  const secs = Math.floor(seconds % 60);

  if (days > 0) return `${days}d ${hours}h`;
  if (hours > 0) return `${hours}h ${minutes}m`;
  if (minutes > 0) return `${minutes}m ${secs}s`;
  return `${secs}s`;
}

const columnHelper = createColumnHelper<DataTableFeatures, ProcessItem>();

const sortableHeader =
  (text: string) =>
  ({
    column,
  }: {
    column: {
      getToggleSortingHandler: () => undefined | ((event: unknown) => void);
    };
  }) =>
    renderComponent(SortableTableHeader, {
      text,
      onclick: column.getToggleSortingHandler(),
    });

export const getColumns = () =>
  columnHelper.columns([
    columnHelper.accessor("pid", {
      header: sortableHeader("Pid"),
    }),
    columnHelper.accessor("name", {
      header: sortableHeader("Name"),
      filterFn: (row, _columnId, filterValue) => {
        const query = String(filterValue ?? "").trim().toLowerCase();
        if (!query) return true;
        return `${row.original.name} ${row.original.command}`
          .toLowerCase()
          .includes(query);
      },
    }),
    columnHelper.accessor("command", {
      id: "cmd",
      header: sortableHeader("Command"),
      cell: (info) => info.getValue() || "-",
    }),
    columnHelper.accessor("username", {
      id: "user",
      header: sortableHeader("User"),
      cell: (info) => info.getValue() || "-",
    }),
    columnHelper.accessor("cpu_usage", {
      id: "cpu",
      header: sortableHeader("CPU"),
      cell: (info) => `${info.getValue().toFixed(1)}%`,
    }),
    columnHelper.accessor("memory_usage", {
      id: "memory",
      header: sortableHeader("Memory"),
      cell: (info) => formatBytes(info.getValue()),
    }),
    columnHelper.accessor("status", {
      id: "state",
      header: sortableHeader("State"),
      cell: (info) => formatProcessStatus(info.getValue()),
    }),
    columnHelper.accessor("uptime", {
      header: sortableHeader("Uptime"),
      cell: (info) => formatUptime(info.getValue()),
    }),
  ]);
