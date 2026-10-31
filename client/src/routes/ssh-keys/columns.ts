import DeleteRequestDialog from "#lib/components/base/DeleteRequestDialog.svelte";
import type { components } from "#lib/types/request-schema.js";
import { createColumnHelper, renderComponent } from "@tanstack/svelte-table";
import type { DataTableFeatures } from "#lib/components/base/data-table-features.js";
import ViewKeyInfo from "./viewKeyInfo.svelte";
import ResourcesLink from "./ResourcesLink.svelte";
import SortableTableHeader from "#lib/components/base/SortableTableHeader.svelte";

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
        return new Date(row.original.created_at).toLocaleString();
      },
    }),
    columnHelper.display({
      id: "view-key",
      cell: ({ row }) => {
        return renderComponent(ViewKeyInfo, {
          id: row.original.id,
        });
      },
    }),
    columnHelper.display({
      id: "delete-key",
      cell: ({ row }) => {
        return renderComponent(DeleteRequestDialog, {
          id: row.original.id,
          name: row.original.name,
          type: "ssh-key",
          useIcon: true,
          onDeleted,
        });
      },
    }),
  ]);
