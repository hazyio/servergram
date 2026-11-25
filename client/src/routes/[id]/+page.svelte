<script lang="ts">
  import { ConnectedServices } from "#lib/components/base/services/services.svelte.ts";
  import type { PageProps } from "./$types";
  import OverView from "./overview/ServerInfo.svelte";

  let { params }: PageProps = $props();
  let server = $derived.by(() => {
    return ConnectedServices.findServer(Number(params.id));
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
  <OverView id={Number(params.id)} />
</div>
