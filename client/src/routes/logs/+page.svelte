<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import Error from "#lib/components/base/error.svelte";
  import CreateGroup from "#lib/components/base/groups/createGroup.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import * as Table from "#lib/components/ui/table/index.js";
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { onMount } from "svelte";
  import { createTable, FlexRender } from "@tanstack/svelte-table";
  import { getColumns } from "./columns.js";

  let loading = $state(true);
  let groups: components["schemas"]["GroupListResponseBody"] | undefined =
    $state();
  let loadError = $state(false);
  let search = $state(page.url.searchParams.get("search") ?? "");
  function updateSearch(event: Event) {
    const value = (event.currentTarget as HTMLInputElement).value;
    table.getColumn("name")?.setFilterValue(value);
    const url = new URL(page.url.href);
    if (value) {
      url.searchParams.set("search", value);
    } else {
      url.searchParams.delete("search");
    }
    goto(url, { state: page.state, shallow: true, replace: true });
  }
  async function loadGroups() {
    loading = true;
    const { data, error } = await NewRequestClient().GET("/api/group/list");
    loading = false;
    if (error) {
      loadError = true;
      return;
    }
    loadError = false;
    groups = data;
  }
  const columns = getColumns();
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return groups?.items ?? [];
    },
    columns,
  });
  onMount(() => {
    table.getColumn("name")?.setFilterValue(search);
    loadGroups();
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
      placeholder="Search groups..."
      bind:value={search}
      oninput={updateSearch}
    />
    <Button variant="secondary" onclick={loadGroups}><RefreshCw /></Button>
    <CreateGroup text="Create Group" onCreated={loadGroups} />
  </div>
  <Loading class="mt-20" {loading} text="Loading..." />
  {#if (loadError || !groups) && !loading}
    <Error
      title="Failed to load groups"
      message="An error occurred while fetching the group list."
      onRetry={loadGroups}
    />
  {:else if groups && !loading}
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
                {groups.items?.length
                  ? "No matching groups."
                  : "No groups yet. Create a group to organize related servers."}
              </Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
      <div class="flex items-center justify-end space-x-2 py-4">
        <Button
          variant="outline"
          size="sm"
          onclick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          onclick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          Next
        </Button>
      </div>
    </div>
  {/if}
</div>
