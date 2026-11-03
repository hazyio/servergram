<script lang="ts">
  import { severPageInfo } from "#lib/states/severPageInfo.svelte.ts";
  import type { PageProps } from "./$types";
  import ServerInfo from "./ServerInfo.svelte";
  import SeverResources from "./ServerResources/ServerResources.svelte";
  import TopSystemProcess from "./TopSystemProcess.svelte";
  let { params }: PageProps = $props();
</script>

<svelte:head>
  <title>{severPageInfo.name} Overview | ServerGram</title>
  <meta
    name="description"
    content="View the details of a server, including its hostname, SSH username, port, group, and SSH key."
  />
</svelte:head>
<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
  <SeverResources id={params.id} />
  <TopSystemProcess id={params.id} />
  <ServerInfo id={params.id} />
</div>
