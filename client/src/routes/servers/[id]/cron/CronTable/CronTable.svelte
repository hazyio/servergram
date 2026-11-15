<script lang="ts">
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import * as Table from "#lib/components/ui/table/index.js";
  import { FlexRender, createTable } from "@tanstack/svelte-table";
  import { untrack } from "svelte";
  import type { CronJob } from "../cronJobs.ts";
  import { getColumns } from "./cronColumns";

  let {
    jobs,
    serverId,
    search = "",
    statusFilter = "all",
    visibleColumns,
    onUpdated,
    onCreated,
    onDeleted,
  }: {
    jobs: CronJob[];
    serverId: number;
    search?: string;
    statusFilter?: string;
    visibleColumns: Record<string, boolean>;
    onUpdated: (job: CronJob) => void;
    onCreated: (job: CronJob) => void;
    onDeleted: (id: number) => void;
  } = $props();

  const columns = untrack(() =>
    getColumns(serverId, onUpdated, onCreated, onDeleted),
  );
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return jobs;
    },
    getRowId: (job) => String(job.id),
    columns,
  });

  const visibleColumnCount = $derived(table.getVisibleLeafColumns().length);

  $effect(() => {
    table.getColumn("job")?.setFilterValue(search);
  });

  $effect(() => {
    table.getColumn("status")?.setFilterValue(statusFilter);
  });

  $effect(() => {
    for (const [id, visible] of Object.entries(visibleColumns)) {
      table.getColumn(id)?.toggleVisibility(visible);
    }
  });
</script>

<Table.Root class="mt-2">
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
      <Table.Row class="h-15">
        {#each row.getVisibleCells() as cell (cell.id)}
          <Table.Cell class="p-0 px-4">
            <FlexRender {cell} />
          </Table.Cell>
        {/each}
      </Table.Row>
    {:else}
      <Table.Row>
        <Table.Cell colspan={visibleColumnCount} class="h-24 text-center">
          {jobs.length ? "No matching cron jobs." : "No cron jobs yet."}
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
