<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import type { components } from "#lib/api/schema.js";
  import appPaths from "#lib/appPaths.ts";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import Error from "#lib/components/base/error.svelte";
  import CreateGroup from "#lib/components/base/groups/createGroup.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import NoItemsFound from "#lib/components/base/noItemsFound.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { onMount } from "svelte";
  import DeleteRequestDialog from "#lib/components/base/DeleteRequestDialog.svelte";
  let loading = $state(true);
  let groups: components["schemas"]["GroupListResponseBody"] | undefined =
    $state();
  let loadError = $state(false);
  let search = $state(page.url.searchParams.get("search") ?? "");
  let filteredGroups = $derived(
    groups?.items?.filter((group) =>
      group.name.toLowerCase().includes(search.trim().toLowerCase()),
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
  onMount(loadGroups);
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
    {#if groups.items?.length && filteredGroups.length}
      <ul class="mt-3 grid md:grid-cols-2 lg:grid-cols-3 gap-2">
        {#each filteredGroups as group (group.id)}
          <li class="bg-secondary rounded-md mt-2 py-3 px-2 flex flex-row">
            <div class="flex gap-2 flex-col grow">
              <span class="font-bold">{group.name}</span>
              <small
                ><a href={appPaths.groups.resources(group.id)}
                  >Resource: {group.resources_attached}</a
                ></small
              >
              <small
                >Created at: {new Date(
                  group.created_at,
                ).toLocaleString()}</small
              >
            </div>
            <DeleteRequestDialog
              id={group.id}
              name={group.name}
              onDeleted={loadGroups}
              type="group"
            />
          </li>
        {/each}
      </ul>
    {:else if groups.items?.length}
      <NoItemsFound
        title="No matching groups"
        message="Try a different name."
      />
    {:else}
      <NoItemsFound
        title="No groups yet"
        message="Create a group to keep related servers together."
      />
    {/if}
  {/if}
</div>
