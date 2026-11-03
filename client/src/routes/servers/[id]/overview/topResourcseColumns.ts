import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";

export type TopProcess = {
  pid: number;
  name: string;
  cpu_usage: number;
  memory_usage: number;
};

const columnHelper = createColumnHelper<DataTableFeatures, TopProcess>();

function formatBytes(bytes: number) {
  if (!Number.isFinite(bytes) || bytes < 0) return "-";
  if (bytes === 0) return "0 B";

  const units = ["B", "KB", "MB", "GB", "TB", "PB"];
  const unitIndex = Math.min(
    Math.floor(Math.log(bytes) / Math.log(1024)),
    units.length - 1,
  );
  const value = bytes / 1024 ** unitIndex;
  return `${new Intl.NumberFormat(undefined, { maximumFractionDigits: 1 }).format(value)} ${units[unitIndex]}`;
}

export const getColumns = () =>
  columnHelper.columns([
    columnHelper.accessor("pid", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Pid",
          onClick: column.getToggleSortingHandler(),
        }),
    }),
    columnHelper.accessor("name", {
      header: "Name",
    }),
    columnHelper.accessor("cpu_usage", {
      header: "CPU",
      cell: (info) => `${info.getValue().toFixed(1)}%`,
    }),
    columnHelper.accessor("memory_usage", {
      header: "Memory",
      cell: (info) => formatBytes(info.getValue()),
    }),
  ]);
