import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";
import FileName from "./FileName.svelte";

export type FileEntry = {
  name: string;
  size: number;
  modified: string;
  mode: string;
  owner: string;
  is_dir: boolean;
};

const columnHelper = createColumnHelper<DataTableFeatures, FileEntry>();

function formatBytes(bytes: number) {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.min(
    Math.floor(Math.log(bytes) / Math.log(1024)),
    units.length - 1,
  );
  return `${new Intl.NumberFormat(undefined, { maximumFractionDigits: 1 }).format(bytes / 1024 ** i)} ${units[i]}`;
}

export const getColumns = () =>
  columnHelper.columns([
    columnHelper.accessor("name", {
      header: "Name",
      cell: (info) =>
        renderComponent(FileName, {
          name: info.getValue(),
          isDir: info.row.original.is_dir,
        }),
    }),
    columnHelper.accessor("size", {
      header: "Size",
      cell: (info) =>
        info.row.original.is_dir ? "-" : formatBytes(info.getValue()),
    }),
    columnHelper.accessor("modified", {
      header: "Modified",
      cell: (info) => new Date(info.getValue()).toLocaleString(),
    }),
    columnHelper.accessor("mode", {
      header: "Mode",
    }),
    columnHelper.accessor("owner", {
      header: "Owner",
    }),
  ]);

export const stubFiles: FileEntry[] = [
  {
    name: "assets",
    size: 4096,
    modified: "2026-10-28T09:12:00Z",
    mode: "drwxr-xr-x",
    owner: "www-data",
    is_dir: true,
  },
  {
    name: "logs",
    size: 4096,
    modified: "2026-11-03T06:30:00Z",
    mode: "drwxr-x---",
    owner: "root",
    is_dir: true,
  },
  {
    name: "index.html",
    size: 2048,
    modified: "2026-10-30T14:05:00Z",
    mode: "-rw-r--r--",
    owner: "www-data",
    is_dir: false,
  },
  {
    name: "app.js",
    size: 183500,
    modified: "2026-11-01T18:44:00Z",
    mode: "-rw-r--r--",
    owner: "www-data",
    is_dir: false,
  },
  {
    name: ".env",
    size: 312,
    modified: "2026-09-12T11:20:00Z",
    mode: "-rw-------",
    owner: "root",
    is_dir: false,
  },
  {
    name: "backup.tar.gz",
    size: 52428800,
    modified: "2026-10-15T02:00:00Z",
    mode: "-rw-r--r--",
    owner: "daniel",
    is_dir: false,
  },
];
