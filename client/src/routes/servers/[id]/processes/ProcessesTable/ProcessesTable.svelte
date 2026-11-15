<script lang="ts">
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import * as Table from "#lib/components/ui/table/index.js";
  import { formatBytes } from "#lib/utils.ts";
  import ChevronDownIcon from "@lucide/svelte/icons/chevron-down";
  import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
  import { FlexRender, createTable } from "@tanstack/svelte-table";
  import { SvelteSet } from "svelte/reactivity";
  import { untrack } from "svelte";
  import {
    formatUptime,
    getColumns,
    type ProcessItem,
  } from "./processesColumns";
  import type { ProcessGrouping } from "../ProcessSettingOption.svelte";

  let {
    processes,
    visibleColumns,
    serverId,
    search = "",
    groupBy = "parent",
  }: {
    processes: ProcessItem[];
    visibleColumns: Record<string, boolean>;
    serverId: number;
    search?: string;
    groupBy?: ProcessGrouping;
  } = $props();
  const columns = untrack(() => getColumns(serverId));
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return processes;
    },
    getRowId: (process) => String(process.pid),
    columns,
  });

  const collapsedGroups = new SvelteSet<string>();

  const groups = $derived.by(() => {
    const grouping = groupBy;
    if (grouping === "none") return null;

    const rows = table.getRowModel().rows;
    const processByPid = new Map(
      rows.map((row) => [row.original.pid, row.original]),
    );
    const groupsByKey = new Map<
      string,
      {
        key: string;
        ppid: number | null;
        parent: ProcessItem | undefined;
        label: string;
        rows: typeof rows;
        cpu: number;
        memory: number;
        uptime: number;
      }
    >();

    for (const row of rows) {
      const process = row.original;
      let key: string;
      let ppid: number | null = null;
      let parent: ProcessItem | undefined;
      let label: string;

      switch (grouping) {
        case "parent": {
          ppid = process.ppid;
          parent = processByPid.get(ppid);
          key = `parent:${ppid}`;
          label = parent?.name || `Parent PID ${ppid}`;
          break;
        }
        case "name": {
          key = `name:${process.name}`;
          label = process.name || "-";
          break;
        }
        case "command": {
          key = `command:${process.command}`;
          label = process.command || "-";
          break;
        }
      }

      let group = groupsByKey.get(key);
      if (!group) {
        group = {
          key,
          ppid,
          parent,
          label,
          rows: [],
          cpu: 0,
          memory: 0,
          uptime: 0,
        };
        groupsByKey.set(key, group);
      }
      group.rows.push(row);
      group.cpu += process.cpu_usage;
      group.memory += process.memory_usage;
      group.uptime = Math.max(group.uptime, process.uptime);
    }

    const groups = [...groupsByKey.values()];

    // Sort groups by their cumulative CPU/memory or longest uptime when one of
    // those columns is the primary sort. Other sorts keep the order of their
    // first sorted row.
    const primarySort = table.atoms.sorting?.get()?.[0];
    if (
      primarySort &&
      (primarySort.id === "cpu" ||
        primarySort.id === "memory" ||
        primarySort.id === "uptime")
    ) {
      const key = primarySort.id;
      const direction = primarySort.desc ? -1 : 1;
      groups.sort((a, b) => (a[key] - b[key]) * direction);
    }

    return groups;
  });

  const groupHeaderLayout = $derived.by(() => {
    const columnIds = table.getVisibleLeafColumns().map((column) => column.id);
    const isAggregate = (id: string): id is "cpu" | "memory" | "uptime" =>
      id === "cpu" || id === "memory" || id === "uptime";
    const firstAggregateIndex = columnIds.findIndex(isAggregate);

    const segments: Array<
      | { kind: "label" | "spacer"; colspan: number }
      | { kind: "cpu" | "memory" | "uptime"; colspan: 1 }
    > = [];

    const labelEnd =
      firstAggregateIndex === -1 ? columnIds.length : firstAggregateIndex;
    if (labelEnd > 0) {
      segments.push({ kind: "label", colspan: labelEnd });
    }

    let spacerColspan = 0;
    for (const id of columnIds.slice(labelEnd)) {
      if (isAggregate(id)) {
        if (spacerColspan > 0) {
          segments.push({ kind: "spacer", colspan: spacerColspan });
          spacerColspan = 0;
        }
        segments.push({ kind: id, colspan: 1 });
      } else {
        spacerColspan++;
      }
    }
    if (spacerColspan > 0) {
      segments.push({ kind: "spacer", colspan: spacerColspan });
    }

    return segments;
  });

  const visibleColumnCount = $derived(table.getVisibleLeafColumns().length);

  function toggleGroup(key: string) {
    if (collapsedGroups.has(key)) {
      collapsedGroups.delete(key);
    } else {
      collapsedGroups.add(key);
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
      {#each groups as group (group.key)}
        <Table.Row class="border-t bg-muted/40 hover:bg-muted/40">
          {#each groupHeaderLayout as segment, index (index)}
            {#if segment.kind === "label"}
              <Table.Cell colspan={segment.colspan} class="p-0">
                <button
                  type="button"
                  class="flex w-full items-center gap-2 px-1 py-3 text-left text-sm md:px-4"
                  aria-expanded={!collapsedGroups.has(group.key)}
                  onclick={() => toggleGroup(group.key)}
                >
                  {#if collapsedGroups.has(group.key)}
                    <ChevronRightIcon class="size-4 shrink-0" />
                  {:else}
                    <ChevronDownIcon class="size-4 shrink-0" />
                  {/if}
                  <span
                    class="min-w-0 truncate font-medium"
                    title={group.label}
                  >
                    {group.label}
                    {#if groupBy === "parent" && group.parent}
                      <span class="font-normal text-muted-foreground"
                        >(PID {group.ppid})</span
                      >
                    {/if}
                  </span>
                  <span class="ms-auto text-xs text-muted-foreground">
                    {group.rows.length}
                    {group.rows.length === 1 ? "process" : "processes"}
                  </span>
                </button>
              </Table.Cell>
            {:else if segment.kind === "cpu"}
              <Table.Cell
                class="px-4 py-3 text-sm font-medium tabular-nums"
                title="Cumulative CPU usage"
              >
                {group.cpu.toFixed(1)}%
              </Table.Cell>
            {:else if segment.kind === "memory"}
              <Table.Cell
                class="px-4 py-3 text-sm font-medium tabular-nums"
                title="Cumulative memory usage"
              >
                {formatBytes(group.memory)}
              </Table.Cell>
            {:else if segment.kind === "uptime"}
              <Table.Cell
                class="px-4 py-3 text-sm font-medium tabular-nums"
                title="Longest uptime in group"
              >
                {formatUptime(group.uptime)}
              </Table.Cell>
            {:else}
              <Table.Cell colspan={segment.colspan}></Table.Cell>
            {/if}
          {/each}
        </Table.Row>
        {#if !collapsedGroups.has(group.key)}
          {#each group.rows as row (row.id)}
            <Table.Row>
              {#each row.getVisibleCells() as cell, index (cell.id)}
                <Table.Cell
                  class={`p-0 px-4 ${index === 0 ? "ps-8" : ""}`}
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
            <Table.Cell class="p-0 px-4">
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
