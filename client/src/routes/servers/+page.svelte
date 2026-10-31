<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import appPaths from "#lib/appPaths.ts";
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import Error from "#lib/components/base/error.svelte";
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
  let servers: components["schemas"]["ListServersResponseBody"] | undefined =
    $state();
  let loadError = $state(false);

  async function loadServers() {
    loading = true;
    const { data, error } = await NewRequestClient().GET("/api/server/list");
    loading = false;
    if (error) {
      loadError = true;
      return;
    }
    loadError = false;
    servers = data;
  }
  const columns = getColumns(loadServers);
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return servers?.items ?? [];
    },
    columns,
  });
  onMount(() => {
    loadServers();
  });
</script>

<svelte:head>
  <title>Servers | ServerGram</title>
  <meta
    name="description"
    content="View and manage the servers connected to your ServerGram account."
  />
</svelte:head>

<div class="p-2">
  <div class="flex gap-2">
    <Input
      name="filter-servers"
      placeholder="Search..."
      autocomplete="off"
      type="search"
      onchange={(e) => {
        table.getColumn("name")?.setFilterValue(e.currentTarget.value);
      }}
      oninput={(e) => {
        table.getColumn("name")?.setFilterValue(e.currentTarget.value);
      }}
    />
    <Button variant="secondary" onclick={loadServers}><RefreshCw /></Button>

    <Button
      ><a href={appPaths.servers.add} class="gap-2 flex items-center">
        Add Server
      </a></Button
    >
  </div>
  <Loading class="mt-20" {loading} text="Loading..." />

  {#if (loadError || !servers) && !loading}
    <Error
      title="Failed to load servers"
      message="An error occurred while fetching the server list."
      onRetry={loadServers}
    />
  {:else if servers && !loading}
    <div class="rounded-md border mt-12">
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
                {servers.items?.length
                  ? "No matching servers."
                  : "No servers yet."}
              </Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    </div>
  {/if}
</div>
