import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";
import TimeAgo from "#lib/components/base/TimeAgo.svelte";
import type { components } from "#lib/types/request-schema.js";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";
import ResourcesLink from "./ResourcesLink.svelte";
import SSHKeyActions from "./SSHKeyActions.svelte";

const columnHelper = createColumnHelper<
  DataTableFeatures,
  components["schemas"]["SSHKeyListItem"]
>();
export const getColumns = (onDeleted: () => void) =>
  columnHelper.columns([
    columnHelper.accessor("name", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Name",
          onclick: column.getToggleSortingHandler(),
        }),
    }),
    columnHelper.accessor("encrypted", {
      header: "Encrypted",
      cell: ({ row }) => (row.original.encrypted ? "Yes" : "No"),
    }),
    columnHelper.accessor("resources_attached", {
      header: ({ column }) =>
        renderComponent(SortableTableHeader, {
          text: "Resources",
          onclick: column.getToggleSortingHandler(),
        }),
      cell: ({ row }) => {
        return renderComponent(ResourcesLink, {
          id: row.original.id,
          count: row.original.resources_attached,
        });
      },
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
        renderComponent(SSHKeyActions, {
          id: row.original.id,
          name: row.original.name,
          onDeleted,
        }),
    }),

  ]);
