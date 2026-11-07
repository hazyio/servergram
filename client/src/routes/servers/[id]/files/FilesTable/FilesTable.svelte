<script lang="ts">
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import * as Table from "#lib/components/ui/table/index.js";
  import { FlexRender, createTable } from "@tanstack/svelte-table";
  import { getColumns, isDir } from "./filesColumns";
  import type { DeleteEntryProps } from "../DeleteEntriesDialog.svelte";
  import type { FileEntry } from "#lib/types/index.ts";

  let {
    files,
    search = "",
    dirsFirst = true,
    visibleColumns,
    onOpen,
    onDelete = $bindable(),
    onRenameClicked = $bindable(),
    onOpenPropertiesClicked = $bindable(),
    onViewClicked = $bindable(),
    onDownloadClicked = $bindable(),
  }: {
    files: FileEntry[];
    search?: string;
    dirsFirst?: boolean;
    visibleColumns: Record<string, boolean>;
    onOpen: (entry: FileEntry) => void;
    onDelete: (entry: DeleteEntryProps) => void;
    onRenameClicked: (entry: DeleteEntryProps) => void;
    onOpenPropertiesClicked: (entry: DeleteEntryProps) => void;
    onViewClicked: (entry: DeleteEntryProps) => void;
    onDownloadClicked: (entry: DeleteEntryProps) => void;
  } = $props();
  let selected = $state<string[]>([]);

  function clearSelection() {
    selected = [];
  }

  $effect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") clearSelection();
    }
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  });

  const columns = getColumns(
    onDelete,
    onRenameClicked,
    onOpenPropertiesClicked,
    onViewClicked,
    onDownloadClicked,
  );
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return files;
    },
    columns,
  });

  $effect(() => {
    table.getColumn("name")?.setFilterValue(search);
  });

  $effect(() => {
    for (const [id, visible] of Object.entries(visibleColumns)) {
      table.getColumn(id)?.toggleVisibility(visible);
    }
  });

  const rows = $derived.by(() => {
    const all = table.getRowModel().rows;
    if (!dirsFirst) return all;
    // Stable partition keeps the active column sort within each group.
    return [
      ...all.filter((r) => isDir(r.original)),
      ...all.filter((r) => !isDir(r.original)),
    ];
  });
</script>

<Table.Root class="mt-2 ">
  <Table.Header>
    {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
      <Table.Row>
        {#each headerGroup.headers as header (header.id)}
          <Table.Head colspan={header.colSpan} class="h-12 px-1 md:px-4">
            {#if !header.isPlaceholder}
              <FlexRender {header} />
            {/if}
          </Table.Head>
        {/each}
      </Table.Row>
    {/each}
  </Table.Header>
  <Table.Body>
    {#each rows as row (row.id)}
      <Table.Row
        class={`${selected.includes(row.original.name ?? "") ? "bg-accent hover:bg-accent/50" : ""} ${isDir(row.original) ? "cursor-pointer" : ""}`}
        onclick={() => {
          if (row.original.name) {
            if (!selected.includes(row.original.name)) {
              selected.push(row.original.name);
            } else {
              selected = selected.filter((name) => name !== row.original.name);
            }
          }
        }}
        ondblclick={() => {
          onOpen(row.original);
          //   clear selection
          selected = [];
        }}
      >
        {#each row.getVisibleCells() as cell (cell.id)}
          <Table.Cell class=" p-0 px-4">
            <FlexRender {cell} />
          </Table.Cell>
        {/each}
      </Table.Row>
    {:else}
      <Table.Row>
        <Table.Cell colspan={columns.length} class="h-24 text-center">
          No files.
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
