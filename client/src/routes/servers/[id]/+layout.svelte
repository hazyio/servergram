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

  let { params,children }: LayoutProps = $props();
  let server = $derived.by(() => {
    return ConnectedServices.findServer(Number(params.id));
  });
  $effect(() => {
    if (server) {
      updatePageServerInfo(server.id, server.name);
    }
  });
</script>

<UrlTabsController
  items={[
    { label: "Overview", href: appPaths.servers.overview(Number(params.id)) },
    { label: "Files", href: appPaths.servers.files(Number(params.id)) },
    { label: "Process", href: appPaths.servers.process(Number(params.id)) },
    { label: "Docker", href: appPaths.servers.docker(Number(params.id)) },
    { label: "Cron", href: appPaths.servers.cron(Number(params.id)) },
  ]}
/>
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
  {@render children()}

  {/if}
</div>
