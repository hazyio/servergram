<script lang="ts">
  import { onMount } from "svelte";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import { Clock } from "@lucide/svelte";
  import type { components } from "#lib/types/request-schema.d.ts";

  import InCard from "./InCard.svelte";
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  let { id } = $props();
  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let information: components["schemas"]["GetServerInfoResponse"] | undefined =
    $state();
  async function getServerInfo() {
    loading = true;
    loadError = null;
    const { data, error } = await NewRequestClient().GET(
      "/api/server/info/{id}",
      {
        params: {
          path: {
            id: id,
          },
        },
      },
    );
    loading = false;
    if (error) {
      loadError = unWrapErrorMessage(error);
      return;
    }

    // Implement the logic to fetch server info here
    information = data;
  }
  onMount(getServerInfo);
</script>

<Card>
  <CardHeader><CardTitle>System Information</CardTitle></CardHeader>
  <CardContent class="space-y-2">
    {#if loading}
      <Loading loading class="mt-12" />
    {:else if loadError || information === undefined}
      <ErrorSection
        message={loadError ?? "Failed to load server information"}
        onRetry={getServerInfo}
      />
    {:else}
      <InCard>OS: {information.os}</InCard>
      <InCard>Kernel: {information.kernel_version}</InCard>
      <InCard>Processor: {information.cpu_model}</InCard>
      <InCard>HostName: {information.hostname}</InCard>
      <InCard>Uptime: {information.uptime}</InCard>
      <InCard>Platform: {information.platform}</InCard>
      <InCard>Architecture: {information.architecture}</InCard>
      <InCard class="mt-5">
        <Clock /> {information.system_time}
      </InCard>
    {/if}
  </CardContent>
</Card>
