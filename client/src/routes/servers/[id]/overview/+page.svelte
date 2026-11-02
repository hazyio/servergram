<script lang="ts">
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import { severPageInfo } from "#lib/states/severPageInfo.svelte.ts";
  import { Clock } from "@lucide/svelte";
  import type { PageProps } from "./$types";
  import InCard from "./InCard.svelte";
  import SeverResources from "./SeverResources.svelte";
  import ServerInfo from "./ServerInfo.svelte";
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
    class="grid md:grid-cols-2 lg:md:grid-cols-4 gap-2 md:col-start-1 md:col-end-8 lg:col-end-9"
  >
    <ServerInfo id={params.id} />
    <SeverResources />
    <Card>
      <CardHeader><CardTitle>Top System Process</CardTitle></CardHeader>
    </Card>
    <Card>
      <CardHeader><CardTitle>Systemd Units</CardTitle></CardHeader>
    </Card>
  </div>

  <Card class="md:col-start-8 md:col-end-13 lg:col-start-9 ">
    <CardHeader><CardTitle>Cron & Scheduled Tasks</CardTitle></CardHeader>
  </Card>
  <Card class="md:col-start-1 md:col-end-13">
    <CardHeader><CardTitle>Logs</CardTitle></CardHeader>
  </Card>
</div>
