<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import NoItemsFound from "#lib/components/base/noItemsFound.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import type { components } from "#lib/types/request-schema.d.ts";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { onMount } from "svelte";
  import SessionCard from "./SessionCard.svelte";

  let loading = $state(true);
  let sessions: components["schemas"]["Session"][] = $state([]);
  let loadError = $state(false);
  let search = $state("");

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

  let visibleSessions = $derived.by(() => {
    const query = search.trim().toLowerCase();
    const filtered = query
      ? sessions.filter(
          (session) =>
            session.os.toLowerCase().includes(query) ||
            session.user_agent.toLowerCase().includes(query),
        )
      : sessions;

    return [...filtered].sort((a, b) => {
      if (a.this_device) return -1;
      if (b.this_device) return 1;
      return b.id - a.id;
    });
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
      bind:value={search}
      placeholder="Search sessions by OS or user agent..."
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
    {#if visibleSessions.length === 0}
      <div class="rounded-md border mt-12 mb-20">
        <NoItemsFound
          title={sessions.length === 0
            ? "No sessions found"
            : "No matching sessions"}
          message={sessions.length === 0
            ? "There are no active sessions for this account."
            : "Try a different search term."}
        />
      </div>
    {:else}
      <div
        class="grid grid-cols-1 gap-4 mt-12 mb-20 sm:grid-cols-2 xl:grid-cols-3"
      >
        {#each visibleSessions as session (session.id)}
          <SessionCard
            {session}
            current={session.this_device}
            onRemoved={loadSessions}
          />
        {/each}
      </div>
    {/if}
  {/if}
</div>
