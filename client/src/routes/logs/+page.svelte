<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import Error from "#lib/components/base/error.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Sort from "#lib/components/base/sort.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import * as Table from "#lib/components/ui/table/index.js";
  import type { components } from "#lib/types/request-schema.d.ts";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { createTable, FlexRender } from "@tanstack/svelte-table";
  import { onMount } from "svelte";
  import { getColumns } from "./columns.js";

  let loading = $state(true);
  let logs: components["schemas"]["LogListResponseBody"] | undefined = $state();
  let sortOption = $state<"all" | "delete" | "create" | "add" | "update">(
    "all",
  );
  let loadError = $state(false);
  let startIndex = $state(0);
  let previousIndex = $state(0);
  let loadLimit = $state(10);
  async function loadLogs() {
    loading = true;
    const { data, error } = await NewRequestClient().GET("/api/logs/list", {
      params: {
        query: {
          "start-index": startIndex,
          limit: loadLimit,
        },
      },
    });
    loading = false;
    if (error) {
      loadError = true;
      return;
    }
    loadError = false;
    logs = data;
    startIndex = data?.next_index ?? startIndex;
    previousIndex = data?.previous_index ?? previousIndex;
  }
  const columns = getColumns();
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return logs?.items ?? [];
    },
    columns,
  });
  onMount(() => {
    loadLogs();
  });
</script>

<svelte:head>
  <title>Groups | ServerGram</title>
  <meta
    name="description"
    content="Create and manage groups to organize your ServerGram servers."
  />
</svelte:head>

<div class="p-2">
  <div class="flex gap-2 justify-center items-center">
    <Input
      placeholder="Search logs..."
      onchange={(e) => {
        table.getColumn("message")?.setFilterValue(e.currentTarget.value);
      }}
      oninput={(e) => {
        table.getColumn("message")?.setFilterValue(e.currentTarget.value);
      }}
    />
    <Button variant="secondary" onclick={loadLogs}><RefreshCw /></Button>
    <Sort
      items={["all", "delete", "create", "add", "update"]}
      bind:value={sortOption}
      caption="Sort logs"
      onChange={(value) => {
        table
          .getColumn("action")
          ?.setFilterValue(value === "all" ? undefined : value);
      }}
    />
    <Button
      variant="outline"
      onclick={() => {
        startIndex = previousIndex;
        loadLogs();
      }}
      disabled={previousIndex === 0}
    >
      Previous
    </Button>
    <Button
      variant="outline"
      onclick={() => {
        loadLogs();
      }}
      disabled={logs?.items?.length === 0 ||
        (logs?.items?.length ?? 0) < loadLimit}
    >
      Next
    </Button>
  </div>
  <Loading class="mt-20" {loading} text="Loading..." />
  {#if (loadError || !logs) && !loading}
    <Error
      title="Failed to load logs"
      message="An error occurred while fetching the log list."
      onRetry={loadLogs}
    />
  {:else if logs && !loading}
    <div class="mt-7 rounded-md border">
      <Table.Root>
        <Table.Header>
          {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
            <Table.Row>
              {#each headerGroup.headers as header (header.id)}
                <Table.Head colspan={header.colSpan}>
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
                <Table.Cell>
                  <FlexRender {cell} />
                </Table.Cell>
              {/each}
            </Table.Row>
          {:else}
            <Table.Row>
              <Table.Cell colspan={columns.length} class="h-24 text-center">
                {logs.items?.length ? "No matching logs." : "No logs yet."}
              </Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    </div>
  {/if}
</div>
