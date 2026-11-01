<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import {
    ConnectedServices,
    servicesState,
  } from "#lib/components/base/services/services.svelte.ts";
  import { updateServerInfo } from "#lib/states/severPageInfo.svelte.ts";
  import { onMount } from "svelte";
  import type { PageProps } from "./$types";

  let { params }: PageProps = $props();
  //   let server = $derived.by(() =>
  //     ConnectedServices.findServer(params.id as unknown as number),
  //   );
  let server = $derived.by(() => {
    servicesState.services;
    return ConnectedServices.findServer(params.id as unknown as number);
  });
  onMount(() => {
    updateServerInfo(
      params.id as unknown as number,
      server?.name ?? "Loading server",
    );
  });
</script>

<svelte:head>
  <title>{server?.name ?? "Server Details"} | ServerGram</title>
  <meta
    name="description"
    content="View the details of a server, including its hostname, SSH username, port, group, and SSH key."
  />
</svelte:head>

<div class="p-2">
  {#if !server}
    <ErrorSection
      title="Failed to load server"
      message="An error occurred while fetching the server details."
      class="mt-20"
    />
  {/if}
</div>
