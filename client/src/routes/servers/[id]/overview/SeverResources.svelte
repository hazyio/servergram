<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Sort from "#lib/components/base/sort.svelte";
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import { Progress } from "#lib/components/ui/progress/index.js";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { onMount } from "svelte";

  let { id } = $props();
  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let information:
    components["schemas"]["GetServerResourcesResponse"] | undefined = $state();
  let refreshFrequencySeconds = $state(1);
  let eventSource: EventSource | undefined = undefined;
  function formatPercent(percent?: number) {
    if (typeof percent !== "number" || !Number.isFinite(percent)) return "-";
    return `${new Intl.NumberFormat(undefined, { maximumFractionDigits: 1 }).format(percent)}%`;
  }

  function progressValue(percent?: number) {
    if (typeof percent !== "number" || !Number.isFinite(percent)) return 0;
    return Math.max(0, Math.min(100, percent));
  }

  function formatBytes(bytes?: number) {
    if (typeof bytes !== "number" || !Number.isFinite(bytes) || bytes < 0) {
      return "-";
    }
    if (bytes === 0) return "0 B";

    const units = ["B", "KB", "MB", "GB", "TB", "PB"];
    const unitIndex = Math.min(
      Math.floor(Math.log(bytes) / Math.log(1024)),
      units.length - 1,
    );
    const value = bytes / 1024 ** unitIndex;
    return `${new Intl.NumberFormat(undefined, { maximumFractionDigits: 1 }).format(value)} ${units[unitIndex]}`;
  }

  function formatUsedTotal(usage?: number, total?: number) {
    if (typeof usage !== "number" || typeof total !== "number") return "-";
    return `${formatBytes(usage)} / ${formatBytes(total)}`;
  }

  async function getServerResources() {
    if (eventSource) {
      eventSource.close();
    }
    loading = true;
    loadError = null;
    eventSource = new EventSource(
      `/api/server/watch-resources/${id}?interval=${refreshFrequencySeconds * 1000}`,
      {
        withCredentials: true, // send cookies; omit if same-origin or no auth
      },
    );
    eventSource.onopen = () => {
      loading = false;
      console.log("connected");
    };
    eventSource.addEventListener("data", (event) => {
      try {
        information = JSON.parse(
          (event as MessageEvent<string>).data,
        ) as components["schemas"]["GetServerResourcesResponse"];
        loading = false;
        loadError = null;
      } catch {
        loading = false;
        loadError = "Failed to parse server resource data";
      }
    });

    eventSource.onerror = () => console.log("disconnected, browser will retry");
  }
  onMount(() => {
    getServerResources();
    return () => {
      if (eventSource) {
        eventSource.close();
      }
    };
  });
</script>

<Card class="md:col-start-1 md:col-end-3">
  <CardHeader>
    <CardTitle class="flex items-center"
      ><span class="grow">Top System Process</span>
      <Sort
        items={[1, 3, 5, 10, 20, 30, 40, 50]}
        bind:value={refreshFrequencySeconds}
        capitalize={false}
        onChange={getServerResources}
        prepend="s"
      />
    </CardTitle>
    <CardTitle>Resources</CardTitle>
  </CardHeader>
  <CardContent class="space-y-6">
    {#if loading}
      <Loading loading class="my-7" />
    {:else if loadError || information === undefined}
      <ErrorSection
        message={loadError ?? "Failed to load server information"}
        onRetry={getServerResources}
      />
    {:else}
      <div class="space-y-2">
        <div class="flex">
          <span class="grow">CPU</span>
          <span class="text-[12px] text-muted-foreground">
            {formatPercent(information.cpu?.usage)}
            {#if information.cpu?.cores !== undefined}
              / {information.cpu.cores} cores
            {/if}
            {#if information.cpu?.threads !== undefined}
              / {information.cpu.threads} threads
            {/if}
          </span>
        </div>
        <Progress value={progressValue(information.cpu?.usage)} />
      </div>
      <div class="space-y-2">
        <div class="flex">
          <span class="grow">Memory</span>
          <div class="text-[12px] text-muted-foreground">
            <span>{formatPercent(information.memory?.percent)}</span>
            <span
              >({formatUsedTotal(
                information.memory?.usage,
                information.memory?.total,
              )})</span
            >
          </div>
        </div>
        <Progress value={progressValue(information.memory?.percent)} />
      </div>
      <div class="space-y-2">
        <div class="flex">
          <span class="grow">Disk</span>
          <span class="text-[12px] text-muted-foreground">
            {formatPercent(information.disk?.percent)}
            ({formatUsedTotal(
              information.disk?.usage,
              information.disk?.total,
            )})
          </span>
        </div>
        <Progress value={progressValue(information.disk?.percent)} />
      </div>
      <div class="space-y-2">
        {#if information.swap && (information.swap.total ?? 0) > 0}
          <div class="flex">
            <span class="grow">Swap</span>
            <span class="text-[12px] text-muted-foreground">
              {formatPercent(information.swap.percent)}
              ({formatUsedTotal(
                information.swap.usage,
                information.swap.total,
              )})
            </span>
          </div>
          <Progress value={progressValue(information.swap.percent)} />
        {:else}
          <div class="flex text-muted-foreground">
            <span class="grow">Swap</span>
            <span class="text-[12px]">Not available</span>
          </div>
        {/if}
      </div>
    {/if}
  </CardContent>
</Card>
