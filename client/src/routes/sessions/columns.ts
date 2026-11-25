import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import TimeAgo from "#lib/components/base/TimeAgo.svelte";
import type { components } from "#lib/types/request-schema.js";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";
import RemoveSession from "./RemoveSession.svelte";

const columnHelper = createColumnHelper<
  DataTableFeatures,
  components["schemas"]["Session"]
>();

// Timestamps may be in seconds or milliseconds.
const toDate = (value: number) => new Date(value < 1e12 ? value * 1000 : value);

export const getColumns = (onDeleted: () => void) =>
  columnHelper.columns([
    columnHelper.accessor("os", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "OS",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: ({ row }) => row.original.os || "Unknown",
    }),
    columnHelper.accessor("user_agent", {
      header: "User Agent",
      cell: ({ row }) => row.original.user_agent || "Unknown",
    }),
    columnHelper.accessor("created_at", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Created",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: ({ row }) =>
        row.original.created_at
          ? renderComponent(TimeAgo, {
              datetime: toDate(row.original.created_at),
            })
          : "—",
    }),
    columnHelper.accessor("expires_at", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Expires",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: ({ row }) =>
        row.original.expires_at
          ? renderComponent(TimeAgo, {
              datetime: toDate(row.original.expires_at),
            })
          : "—",
    }),
    columnHelper.display({
      id: "actions",
      cell: ({ row }) =>
        renderComponent(RemoveSession, {
          id: row.original.id,
          os: row.original.os,
          onDeleted,
        }),
    }),
  ]);
