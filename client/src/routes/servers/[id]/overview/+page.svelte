<script lang="ts">
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardHeader } from "#lib/components/ui/card/index.ts";
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
<div class="grid grid-cols-1 gap-2 md:grid-cols-2 xl:grid-cols-12">
  <SeverResources id={params.id} />
  <TopSystemProcess id={params.id} />
  <Card class="min-w-0 md:col-span-2 xl:col-span-4">
    <CardHeader><CardTitle>Cron & Scheduled Tasks</CardTitle></CardHeader>
  </Card>
  <Card class="min-w-0 md:col-span-2 xl:col-span-12">
    <CardHeader><CardTitle>Logs</CardTitle></CardHeader>
  </Card>
  <ServerInfo id={params.id} />
</div>
