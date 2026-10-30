<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import type { components } from "#lib/types/request-schema.d.ts";
  import appPaths from "#lib/appPaths.ts";
  import DeleteRequestDialog from "#lib/components/base/DeleteRequestDialog.svelte";
  import Error from "#lib/components/base/error.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import NoItemsFound from "#lib/components/base/noItemsFound.svelte";
  import Sort, { type SortDirection } from "#lib/components/base/sort.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { onMount } from "svelte";
  import { startServerConnection } from "#lib/components/base/services/services.svelte.ts";
  let sortDirection: SortDirection = $state("asc");
  let sortValue: "Name" | "Date" | "Size" = $state("Name");
  let loading = $state(true);
  let servers: components["schemas"]["ListServersResponseBody"] | undefined =
    $state();
  let loadError = $state(false);
  let search = $state(page.url.searchParams.get("search") ?? "");
  let filteredServers = $derived(
    servers?.items?.filter((server) =>
      server.name.toLowerCase().includes(search.trim().toLowerCase()),
    ) ?? [],
  );
  function updateSearch(event: Event) {
    const value = (event.currentTarget as HTMLInputElement).value;
    const url = new URL(page.url.href);
    if (value) {
      url.searchParams.set("search", value);
    } else {
      url.searchParams.delete("search");
    }
    goto(url, { state: page.state, shallow: true, replace: true });
  }
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
  onMount(loadServers);
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
      oninput={updateSearch}
      bind:value={search}
    />
    <Sort
      bind:value={sortValue}
      items={["Name", "Date", "Size"]}
      bind:direction={sortDirection}
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
    {#if servers.items?.length && filteredServers.length}
      <ul class="mt-3 grid md:grid-cols-2 lg:grid-cols-3 gap-2">
        {#each filteredServers as server (server.id)}
          <li class="bg-secondary rounded-md mt-2 py-3 px-2 flex flex-row">
            <div class="flex gap-2 flex-col grow">
              <span class="font-bold">{server.name}</span>
              <small>Host: {server.host}</small>
              <small>Trusted: {server.trusted ? "Yes" : "No"}</small>
              <small
                >Created at: {new Date(
                  server.created_at,
                ).toLocaleString()}</small
              >
            </div>
            <Button
              onclick={() => {
                startServerConnection(server.id);
              }}>Connect</Button
            >
            <DeleteRequestDialog
              id={server.id}
              name={server.name}
              onDeleted={loadServers}
              type="server"
            />
          </li>
        {/each}
      </ul>
    {:else if servers.items?.length}
      <NoItemsFound
        title="No matching servers"
        message="Try a different name."
      />
    {:else}
      <NoItemsFound title="No servers yet" message="Create a server." />
    {/if}
  {/if}
</div>
