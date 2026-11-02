<script lang="ts">
  import Card from "#lib/components/ui/card/card.svelte";
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import { Clock } from "@lucide/svelte";
  import InCard from "./InCard.svelte";
  import SeverResources from "./SeverResources.svelte";
</script>

<div class="grid gap-2">
  <div
    class="grid md:grid-cols-2 lg:md:grid-cols-4 gap-2 md:col-start-1 md:col-end-8 lg:col-end-9"
  >
    <Card>
      <CardHeader><CardTitle>System Information</CardTitle></CardHeader>
      <CardContent class="space-y-2">
        <InCard>OS:</InCard>
        <InCard>Kernel:</InCard>
        <InCard>Processor:</InCard>
        <InCard>HostName:</InCard>
        <InCard>Uptime:</InCard>
        <InCard class="mt-5">
          <Clock />
        </InCard>
      </CardContent>
    </Card>
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
