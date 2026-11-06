import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import type { components } from "#lib/types/request-schema.d.ts";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";
import FileName from "./FileName.svelte";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import type { DeleteEntryProps } from "../DeleteEntriesDialog.svelte";
import { formatBytes } from "#lib/utils.ts";

export type FileEntry = components["schemas"]["FileEntry"];

const KIND_DIR = 1;
const KIND_SYMLINK = 2;

export const isDir = (entry: FileEntry) => entry.kind === KIND_DIR;

const columnHelper = createColumnHelper<DataTableFeatures, FileEntry>();


function formatMode(entry: FileEntry) {
  const mode = entry.mode ?? 0;
  const type =
    entry.kind === KIND_DIR ? "d" : entry.kind === KIND_SYMLINK ? "l" : "-";
  const bits = "rwxrwxrwx";
  let out = type;
  for (let i = 0; i < 9; i++) {
    out += mode & (1 << (8 - i)) ? bits[i] : "-";
  }
  return out;
}

export const getColumns = (
  onDelete: (entry: DeleteEntryProps) => void,
  onRenameClicked: (entry: DeleteEntryProps) => void,
  onOpenPropertiesClicked: (entry: DeleteEntryProps) => void,
) =>
  columnHelper.columns([
    columnHelper.accessor((row) => row.name ?? "", {
      id: "name",
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Name",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: (info) =>
        renderComponent(FileName, {
          name: info.getValue(),
          isDir: isDir(info.row.original),
          onDelete: onDelete,
          onRenameClicked: onRenameClicked,
          onOpenPropertiesClicked: onOpenPropertiesClicked,
        }),
    }),
    columnHelper.accessor((row) => row.size ?? 0, {
      id: "size",
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Size",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: (info) =>
        isDir(info.row.original) ? "-" : formatBytes(info.getValue()),
    }),
    columnHelper.accessor((row) => row.modified_unix_ns ?? 0, {
      id: "modified",
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Modified",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: (info) => new Date(info.getValue() / 1e6).toLocaleString(),
    }),
    columnHelper.accessor((row) => formatMode(row), {
      id: "mode",
      header: "Mode",
    }),
    columnHelper.accessor((row) => `${row.uid ?? 0}:${row.gid ?? 0}`, {
      id: "owner",
      header: "Owner",
    }),
  ]);
