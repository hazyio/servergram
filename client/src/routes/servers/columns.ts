import type { components } from "#lib/types/request-schema.js";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";
import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import ServerActions from "./ServerActions.svelte";
import ServerName from "./ServerName.svelte";
import TimeAgo from "#lib/components/base/TimeAgo.svelte";

const columnHelper = createColumnHelper<
  DataTableFeatures,
  components["schemas"]["ServerItem"]
>();

export const getColumns = (onDeleted: () => void) =>
  columnHelper.columns([
    columnHelper.accessor("name", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Name",
          onClick: column.getToggleSortingHandler(),
        }),
      cell: ({ row }) => {
        return renderComponent(ServerName, {
          id: row.original.id,
          name: row.original.name,
        });
      },
    }),
    columnHelper.accessor("host", {
      header: "Host",
    }),
    columnHelper.accessor("trusted", {
      header: "Trusted",
      cell: ({ row }) => (row.original.trusted ? "Yes" : "No"),
    }),
    columnHelper.accessor("created_at", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Created",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: ({ row }) => {
        return renderComponent(TimeAgo, {
          datetime: row.original.created_at,
        });
      },
    }),
    columnHelper.display({
      id: "actions",
      cell: ({ row }) =>
        renderComponent(ServerActions, {
          id: row.original.id,
          name: row.original.name,
          onDeleted,
        }),
    }),
  ]);
