<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import type { components } from "#lib/api/schema.js";
  import Error from "#lib/components/base/error.svelte";
  import CreateGroup from "#lib/components/base/groups/createGroup.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import NoItemsFound from "#lib/components/base/noItemsFound.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { onMount } from "svelte";
  let loading = $state(true);
  let groups: components["schemas"]["GroupListResponseBody"] | undefined =
    $state();
  let loadError = $state(false);
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

<div class="p-2">
  <div class="flex gap-2 justify-center items-center">
    <Input />
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
    {#if groups.items?.length}
      <ul>
        {#each groups.items as group (group.id)}
          <li>{group.name}</li>
        {/each}
      </ul>
    {:else}
      <NoItemsFound
        title="No groups yet"
        message="Create a group to keep related servers together."
      />
    {/if}
  {/if}
</div>
