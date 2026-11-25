<script lang="ts">
  import ServerIcon from "@lucide/svelte/icons/server";
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

<Card class="min-w-0 md:col-span-2 xl:col-span-3">
  <CardHeader class="pb-3">
    <div class="flex items-center gap-3">
      <span
        class="grid size-9 place-items-center rounded-md bg-primary/10 text-primary"
      >
        <ServerIcon class="size-5" />
      </span>
      <CardTitle>System Information</CardTitle>
    </div>
  </CardHeader>
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
        class="grid grid-cols-1 gap-x-4 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-6"
      >
        <div
          class="flex min-w-0 flex-col gap-1 border-t py-3 border-border/70 pt-2 sm:col-span-2 md:col-span-4 lg:col-span-6"
        >
          <dt class="text-xs text-muted-foreground">Host name</dt>
          <dd class="wrap-break-word text-sm font-medium">
            {information.hostname || "-"}
          </dd>
        </div>
        <div class="flex min-w-0 flex-col py-2 gap-1 border-t border-border/70 pt-2">
          <dt class="text-xs text-muted-foreground">Architecture</dt>
          <dd class="wrap-break-word text-sm font-medium">
            {information.architecture || "-"}
          </dd>
        </div>
        <div class="flex min-w-0 flex-col py-2 gap-1 border-t border-border/70 pt-2">
          <dt class="text-xs text-muted-foreground">Kernel version</dt>
          <dd class="wrap-break-word text-sm font-medium">
            {information.kernel_version || "-"}
          </dd>
        </div>
        <div class="flex min-w-0 flex-col py-2 gap-1 border-t border-border/70 pt-2">
          <dt class="text-xs text-muted-foreground">OS</dt>
          <dd class="wrap-break-word text-sm font-medium">
            {information.os || "-"}
          </dd>
        </div>
        <div class="flex min-w-0 flex-col py-2 gap-1 border-t border-border/70 pt-2">
          <dt class="text-xs text-muted-foreground">Platform</dt>
          <dd class="wrap-break-word text-sm font-medium">
            {information.platform || "-"}
          </dd>
        </div>
        <div class="flex min-w-0 flex-col py-2 gap-1 border-t border-border/70 pt-2">
          <dt class="text-xs text-muted-foreground">Processor</dt>
          <dd class="wrap-break-word text-sm font-medium">
            {information.cpu_model || "-"}
          </dd>
        </div>
        <div class="flex min-w-0 flex-col py-2 gap-1 border-t border-border/70 pt-2">
          <dt class="text-xs text-muted-foreground">System time</dt>
          <dd class="flex flex-wrap items-center gap-2 text-sm font-medium">
            {formattedSystemTime ?? "-"}
            {#if formattedSystemTime}
              <span class="text-muted-foreground text-xs">UTC</span>
            {/if}
          </dd>
        </div>
        <div class="flex min-w-0 flex-col py-2 gap-1 border-t border-border/70 pt-2">
          <dt class="text-xs text-muted-foreground">Uptime</dt>
          <dd class="text-sm font-medium">
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
