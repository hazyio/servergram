<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import * as Table from "#lib/components/ui/table/index.js";
  import type { components } from "#lib/types/request-schema.d.ts";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { createTable, FlexRender } from "@tanstack/svelte-table";
  import { onMount } from "svelte";
  import { getColumns } from "./columns.js";

  let loading = $state(true);
  let sessions: components["schemas"]["Session"][] = $state([]);
  let loadError = $state(false);

  async function loadSessions() {
    loading = true;
    const { data, error } = await NewRequestClient().GET(
      "/api/account/sessions/list",
    );
    loading = false;
    if (error) {
      loadError = true;
      return;
    }
    loadError = false;
    sessions = data.sessions ?? [];
  }

  const columns = getColumns(loadSessions);
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return sessions;
    },
    columns,
  });

  onMount(loadSessions);
</script>

<svelte:head>
  <title>Sessions | ServerGram</title>
  <meta
    name="description"
    content="View and remove the devices signed in to your ServerGram account."
  />
</svelte:head>

<div class="p-2">
  <div
    class="flex gap-2 justify-center items-center sticky z-10 top-0 pt-2 bg-background h-15"
  >
    <Input
      placeholder="Search sessions by user agent..."
      onchange={(e) => {
        table.getColumn("user_agent")?.setFilterValue(e.currentTarget.value);
      }}
      oninput={(e) => {
        table.getColumn("user_agent")?.setFilterValue(e.currentTarget.value);
      }}
    />
    <Button variant="secondary" onclick={loadSessions}><RefreshCw /></Button>
  </div>

  <Loading class="mt-20" {loading} text="Loading sessions..." />
  {#if loadError && !loading}
    <Error
      title="Failed to load sessions"
      message="An error occurred while fetching the session list."
      onRetry={loadSessions}
    />
  {:else if !loading}
    <div class="rounded-md border mt-12 mb-20">
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
            <Table.Row data-state={row.getIsSelected() && "selected"}>
              {#each row.getVisibleCells() as cell (cell.id)}
                <Table.Cell>
                  <FlexRender {cell} />
                </Table.Cell>
              {/each}
            </Table.Row>
          {:else}
            <Table.Row>
              <Table.Cell colspan={columns.length} class="h-24 text-center">
                No sessions found.
              </Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    </div>
  {/if}
</div>
