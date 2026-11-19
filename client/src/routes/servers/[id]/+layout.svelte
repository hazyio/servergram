<script lang="ts">
  import appPaths from "#lib/appPaths.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import {
    ConnectedServices,
    servicesState,
  } from "#lib/components/base/services/services.svelte.ts";
  import UrlTabsController from "#lib/components/base/UrlTabsController.svelte";
  import { updatePageServerInfo } from "#lib/states/severPageInfo.svelte.ts";
  import type { LayoutProps } from "./$types";

  let { params, children }: LayoutProps = $props();
  let serverId = $derived(Number(params.id));
  let server = $derived(ConnectedServices.findServer(serverId));
  let tabs = $derived([
    { label: "Overview", href: appPaths.servers.overview(serverId) },
    { label: "Files", href: appPaths.servers.files(serverId) },
    { label: "Processes", href: appPaths.servers.processes(serverId) },
    { label: "Docker", href: appPaths.servers.docker(serverId) },
    { label: "Cron", href: appPaths.servers.cron(serverId) },
  ]);
  $effect(() => {
    if (server) {
            updatePageServerInfo(server.id, server.name);
    }
  });
</script>

<UrlTabsController items={tabs} />
<div class="p-2">
  {#if servicesState.connecting || servicesState.gettingConnected || (server && server.status === "connecting")}
    <Loading loading containerClass="mt-20" />
  {:else if !server}
    <ErrorSection
      title="Server not connected"
      message="This server is not connected, you can connect it by clicking the button below."
      class="mt-20"
      retryText="Connect"
      onRetry={() => {
        ConnectedServices.startServerConnection(serverId, `Server ${serverId}`);
      }}
    />
  {:else}
    {@render children()}
  {/if}
</div>
