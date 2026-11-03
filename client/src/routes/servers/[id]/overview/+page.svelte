<script lang="ts">
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardHeader } from "#lib/components/ui/card/index.ts";
  import { severPageInfo } from "#lib/states/severPageInfo.svelte.ts";
  import type { PageProps } from "./$types";
  import ServerInfo from "./ServerInfo.svelte";
  import SeverResources from "./SeverResources.svelte";
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
<div class="grid gap-2">
  <div
    class="grid grid-cols-1 md:grid-cols-2 lg:md:grid-cols-4 gap-2 md:col-start-1 md:col-end-9 lg:col-end-10"
  >
    <SeverResources id={params.id} />
    <TopSystemProcess id={params.id} />
  </div>
  <Card class="col-end-13 md:col-start-9  lg:col-start-10 ">
    <CardHeader><CardTitle>Cron & Scheduled Tasks</CardTitle></CardHeader>
  </Card>
  <Card class="md:col-start-1 md:col-end-13">
    <CardHeader><CardTitle>Logs</CardTitle></CardHeader>
  </Card>
    <ServerInfo id={params.id} />

</div>
