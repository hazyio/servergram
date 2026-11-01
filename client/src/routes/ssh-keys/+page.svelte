<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import CreateSSHKey from "#lib/components/base/ssh/createSHHKey.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import type { components } from "#lib/types/request-schema.d.ts";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { onMount } from "svelte";

  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import * as Table from "#lib/components/ui/table/index.js";
  import { createTable, FlexRender } from "@tanstack/svelte-table";
  import { getColumns } from "./columns.js";

  let loading = $state(true);
  let sshKeys: components["schemas"]["SSHKeyListResponseBody"] | undefined =
    $state();

  let loadError = $state(false);

  async function loadSSHKeys() {
    loading = true;
    const { data, error } = await NewRequestClient().GET("/api/ssh-key/list");
    loading = false;
    if (error) {
      loadError = true;
      return;
    }
    loadError = false;
    sshKeys = data;
  }
  onMount(loadSSHKeys);
  const columns = getColumns(loadSSHKeys);
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return sshKeys?.items ?? [];
    },
    columns,
  });
</script>

<svelte:head>
  <title>SSH Keys | ServerGram</title>
  <meta
    name="description"
    content="Manage SSH keys used to connect to your ServerGram servers."
  />
</svelte:head>

<div class="p-2">
  <div
    class="flex gap-2 justify-center items-center sticky z-10 top-0 pt-2 bg-background h-15"
  >
    <Input
      placeholder="Search SSH keys by name..."
      onchange={(e) => {
        table.getColumn("name")?.setFilterValue(e.currentTarget.value);
      }}
      oninput={(e) => {
        table.getColumn("name")?.setFilterValue(e.currentTarget.value);
      }}
    />
    <Button variant="secondary" onclick={loadSSHKeys}><RefreshCw /></Button>
    <CreateSSHKey text="Add SSH Key" onCreated={loadSSHKeys} />
  </div>

  <Loading class="mt-20" {loading} text="Loading..." />
  {#if (loadError || !sshKeys) && !loading}
    <Error
      title="Failed to load SSH keys"
      message="An error occurred while fetching the SSH key list."
      onRetry={loadSSHKeys}
    />
  {:else if sshKeys && !loading}
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
                No results.
              </Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    </div>
  {/if}
</div>
