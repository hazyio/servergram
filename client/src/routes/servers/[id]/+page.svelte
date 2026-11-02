<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import {
    ConnectedServices,
    servicesState,
  } from "#lib/components/base/services/services.svelte.ts";
  import { updatePageServerInfo } from "#lib/states/severPageInfo.svelte.ts";
  import type { PageProps } from "./$types";
  import * as Tabs from "#lib/components/ui/tabs/index.js";
  import OverView from "./overview/ServerOverView.svelte";

  let { params }: PageProps = $props();
  let server = $derived.by(() => {
    return ConnectedServices.findServer(Number(params.id));
  });
  $effect(() => {
    if (server) {
      updatePageServerInfo(server.id, server.name);
    }
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
  {#if servicesState.connecting || servicesState.gettingConnected}
    <Loading loading containerClass="mt-20" />
  {:else if !server}
    <ErrorSection
      title="Server not connected"
      message="This server is not connected, you can connect it by clicking the button below."
      class="mt-20"
      retryText="Connect"
      onRetry={() => {
        ConnectedServices.startServerConnection(
          Number(params.id),
          `Server ${params.id}`,
        );
      }}
    />
  {:else}
    <Tabs.Root value="overview" class="w-full">
      <Tabs.List class="w-full">
        <Tabs.Trigger value="overview">Overview</Tabs.Trigger>
        <Tabs.Trigger value="files">Files</Tabs.Trigger>
        <Tabs.Trigger value="docker">Docker</Tabs.Trigger>
        <Tabs.Trigger value="cron">Cron</Tabs.Trigger>
      </Tabs.List>
      <Tabs.Content value="overview"><OverView /></Tabs.Content>
      <Tabs.Content value="docker">Docker content</Tabs.Content>
      <Tabs.Content value="cron">Cron content</Tabs.Content>
    </Tabs.Root>
  {/if}
</div>
