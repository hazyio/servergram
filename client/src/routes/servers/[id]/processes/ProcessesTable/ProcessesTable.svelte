<script lang="ts">
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import * as Table from "#lib/components/ui/table/index.js";
  import ChevronDownIcon from "@lucide/svelte/icons/chevron-down";
  import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
  import { FlexRender, createTable } from "@tanstack/svelte-table";
  import { SvelteSet } from "svelte/reactivity";
  import { getColumns, type ProcessItem } from "./processesColumns";

  let {
    processes,
    visibleColumns,
    search = "",
    groupByParent = true,
  }: {
    processes: ProcessItem[];
    visibleColumns: Record<string, boolean>;
    search?: string;
    groupByParent?: boolean;
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

  const collapsedGroups = new SvelteSet<number>();

  const groups = $derived.by(() => {
    if (!groupByParent) return null;

    const rows = table.getRowModel().rows;
    const processByPid = new Map(
      rows.map((row) => [row.original.pid, row.original]),
    );
    const groupsByParent = new Map<
      number,
      { ppid: number; parent: ProcessItem | undefined; rows: typeof rows }
    >();

    for (const row of rows) {
      const ppid = row.original.ppid;
      let group = groupsByParent.get(ppid);
      if (!group) {
        group = { ppid, parent: processByPid.get(ppid), rows: [] };
        groupsByParent.set(ppid, group);
      }
      group.rows.push(row);
    }

    // Groups keep the order of their first row in the sorted row model, so
    // sorting reorders groups as well as the rows inside each group.
    return [...groupsByParent.values()];
  });

  const visibleColumnCount = $derived(table.getVisibleLeafColumns().length);

  function toggleGroup(ppid: number) {
    if (collapsedGroups.has(ppid)) {
      collapsedGroups.delete(ppid);
    } else {
      collapsedGroups.add(ppid);
    }
  }

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
    {#if groups}
      {#each groups as group (group.ppid)}
        <Table.Row class="border-t bg-muted/40 hover:bg-muted/40">
          <Table.Cell colspan={visibleColumnCount} class="p-0">
            <button
              type="button"
              class="flex w-full items-center gap-2 px-1 py-3 text-left text-sm md:px-4"
              aria-expanded={!collapsedGroups.has(group.ppid)}
              onclick={() => toggleGroup(group.ppid)}
            >
              {#if collapsedGroups.has(group.ppid)}
                <ChevronRightIcon class="size-4 shrink-0" />
              {:else}
                <ChevronDownIcon class="size-4 shrink-0" />
              {/if}
              <span class="font-medium">
                {#if group.parent}
                  {group.parent.name}
                  <span class="font-normal text-muted-foreground"
                    >(PID {group.ppid})</span
                  >
                {:else}
                  Parent PID {group.ppid}
                {/if}
              </span>
              <span class="ms-auto text-xs text-muted-foreground">
                {group.rows.length}
                {group.rows.length === 1 ? "process" : "processes"}
              </span>
            </button>
          </Table.Cell>
        </Table.Row>
        {#if !collapsedGroups.has(group.ppid)}
          {#each group.rows as row (row.id)}
            <Table.Row>
              {#each row.getVisibleCells() as cell, index (cell.id)}
                <Table.Cell
                  class={`px-4 py-2 ${index === 0 ? "ps-8" : ""}`}
                >
                  <FlexRender {cell} />
                </Table.Cell>
              {/each}
            </Table.Row>
          {/each}
        {/if}
      {:else}
        <Table.Row>
          <Table.Cell colspan={visibleColumnCount} class="h-24 text-center">
            No processes.
          </Table.Cell>
        </Table.Row>
      {/each}
    {:else}
      {#each table.getRowModel().rows as row (row.id)}
        <Table.Row>
          {#each row.getVisibleCells() as cell (cell.id)}
            <Table.Cell class="px-4 py-2">
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
    {/if}
  </Table.Body>
</Table.Root>
