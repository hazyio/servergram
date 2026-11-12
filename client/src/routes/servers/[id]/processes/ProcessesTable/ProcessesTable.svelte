<script lang="ts">
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import * as Table from "#lib/components/ui/table/index.js";
  import { FlexRender, createTable } from "@tanstack/svelte-table";
  import { getColumns, type ProcessItem } from "./processesColumns";

  let {
    processes,
    visibleColumns,
    search = "",
  }: {
    processes: ProcessItem[];
    visibleColumns: Record<string, boolean>;
    search?: string;
  } = $props();
  const columns = getColumns();
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return processes;
    },
    getRowId: (process) => String(process.pid),
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
    {#each table.getRowModel().rows as row (row.id)}
      <Table.Row>
        {#each row.getVisibleCells() as cell (cell.id)}
          <Table.Cell class=" p-0 px-4">
            <FlexRender {cell} />
          </Table.Cell>
        {/each}
      </Table.Row>
    {:else}
      <Table.Row>
        <Table.Cell colspan={columns.length} class="h-24 text-center">
          No processes.
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
