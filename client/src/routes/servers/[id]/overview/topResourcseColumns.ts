import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";

const columnHelper = createColumnHelper<
  DataTableFeatures,
  {
    pid: number;
    name: string;
    cpu_usage: number;
    memory_usage: number;
    disk_usage: number;
  }
  // components["schemas"]["ServerItem"]
>();

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
    }),
    columnHelper.accessor("memory_usage", {
      header: "Memory",
    }),
    columnHelper.accessor("disk_usage", {
      header: "Disk",
    }),
  ]);
