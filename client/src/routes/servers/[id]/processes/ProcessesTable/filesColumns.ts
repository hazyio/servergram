import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import { formatBytes } from "#lib/utils.ts";
import { createColumnHelper } from "@tanstack/svelte-table";

interface p {
  name: string;
  pid: number;
  cmd: string;
  user: string;
  cpu: number;
  memory: number;
  state: string;
  uptime: number;
}
const columnHelper = createColumnHelper<DataTableFeatures, p>();

export const getColumns = () =>
  columnHelper.columns([
    columnHelper.accessor("pid", {
      header: "Pid",
    }),
    columnHelper.accessor("name", {
      header: "Name",
    }),
    columnHelper.accessor("cmd", {
      header: "Command",
    }),
    columnHelper.accessor("user", {
      header: "User",
    }),
    columnHelper.accessor("cpu", {
      header: "Cpu",
    }),
    columnHelper.accessor("memory", {
      header: "Memory",
      cell: (info) => formatBytes(info.getValue()),
    }),
    columnHelper.accessor("state", {
      header: "State",
    }),
    columnHelper.accessor("uptime", {
      header: "Uptime",
    }),
  ]);
