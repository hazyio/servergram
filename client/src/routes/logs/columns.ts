import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import TimeAgo from "#lib/components/base/TimeAgo.svelte";
import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import type { components } from "#lib/types/request-schema.js";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";

const columnHelper = createColumnHelper<
  DataTableFeatures,
  components["schemas"]["LogListItem"]
>();

export const getColumns = () =>
  columnHelper.columns([
    columnHelper.accessor("action", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Action",
          onclick: column.getToggleSortingHandler(),
        }),
    }),
    columnHelper.accessor("message", {
      header: "Message",
    }),

    columnHelper.accessor("created_at", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Time",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: ({ row }) =>
        renderComponent(TimeAgo, {
          datetime: row.original.created_at,
        }),
    }),
  ]);
