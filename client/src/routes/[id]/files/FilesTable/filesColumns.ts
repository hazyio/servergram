import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import { KIND_DIR, type FileEntry } from "#lib/types/index.ts";
import { formatBytes, formatMode } from "#lib/utils.ts";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";
import type { DeleteEntryProps } from "../DeleteEntriesDialog.svelte";
import FileName from "./FileName.svelte";

export const isDir = (entry: FileEntry) => entry.kind === KIND_DIR;

const columnHelper = createColumnHelper<DataTableFeatures, FileEntry>();

export const getColumns = (
  onDelete: (entry: DeleteEntryProps) => void,
  onRenameClicked: (entry: DeleteEntryProps) => void,
  onOpenPropertiesClicked: (entry: DeleteEntryProps) => void,
  onViewClicked: (entry: DeleteEntryProps) => void,
  onDownloadClicked: (entry: DeleteEntryProps) => void,
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
          onViewClicked: onViewClicked,
          onDownloadClicked: onDownloadClicked,
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
