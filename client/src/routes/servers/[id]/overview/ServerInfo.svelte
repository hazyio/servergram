<script lang="ts">
  import { onMount } from "svelte";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";
  import type { components } from "#lib/types/request-schema.d.ts";

  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  let { id } = $props();
  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let information: components["schemas"]["GetServerInfoResponse"] | undefined =
    $state();
  let uptimeStartedAt = $derived(
    typeof information?.uptime === "number"
      ? new Date(Date.now() - information.uptime * 1000)
      : undefined,
  );
  let formattedSystemTime = $derived.by(() => {
    const systemTime = information?.system_time;
    if (!systemTime) return undefined;

    const date = new Date(systemTime);
    if (Number.isNaN(date.getTime())) return systemTime;

    return new Intl.DateTimeFormat(undefined, {
      dateStyle: "medium",
      timeStyle: "medium",
      timeZone: "UTC",
    }).format(date);
  });
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

    information = data;
  }
  onMount(getServerInfo);
</script>

<Card class="col-start-1 col-end-13">
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
      <dl
        class="grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-6"
      >
        <div class="bg-muted rounded-md px-2 py-1 flex flex-col gap-1">
          <dt class="text-muted-foreground">Architecture</dt>
          <dd class="grow flex items-center wrap-break-word">
            {information.architecture || "-"}
          </dd>
        </div>
        <div class="bg-muted rounded-md px-2 py-1 flex flex-col gap-1">
          <dt class="text-muted-foreground">Host name</dt>
          <dd class="grow flex items-center wrap-break-word">
            {information.hostname || "-"}
          </dd>
        </div>
        <div class="bg-muted rounded-md px-2 py-1 flex flex-col gap-1">
          <dt class="text-muted-foreground">Kernel version</dt>
          <dd class="grow flex items-center wrap-break-word">
            {information.kernel_version || "-"}
          </dd>
        </div>
        <div class="bg-muted rounded-md px-2 py-1 flex flex-col gap-1">
          <dt class="text-muted-foreground">OS</dt>
          <dd class="grow flex items-center wrap-break-word">
            {information.os || "-"}
          </dd>
        </div>
        <div class="bg-muted rounded-md px-2 py-1 flex flex-col gap-1">
          <dt class="text-muted-foreground">Platform</dt>
          <dd class="grow flex items-center wrap-break-word">
            {information.platform || "-"}
          </dd>
        </div>
        <div class="bg-muted rounded-md px-2 py-1 flex flex-col gap-1">
          <dt class="text-muted-foreground">Processor</dt>
          <dd class="grow flex items-center wrap-break-word">
            {information.cpu_model || "-"}
          </dd>
        </div>
        <div class="bg-muted rounded-md px-2 py-1 flex flex-col gap-1">
          <dt class="text-muted-foreground">System time</dt>
          <dd class="grow flex flex-wrap items-center gap-2">
            {formattedSystemTime ?? "-"}
            {#if formattedSystemTime}
              <span class="text-muted-foreground text-xs">UTC</span>
            {/if}
          </dd>
        </div>
        <div class="bg-muted rounded-md px-2 py-1 flex flex-col gap-1">
          <dt class="text-muted-foreground">Uptime</dt>
          <dd class="grow flex items-center">
            {#if uptimeStartedAt}
              <TimeAgo datetime={uptimeStartedAt} />
            {:else}
              -
            {/if}
          </dd>
        </div>
      </dl>
    {/if}
  </CardContent>
</Card>
