<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Sort from "#lib/components/base/sort.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { addEventSourceTimeout } from "#lib/utils.ts";
  import { ChartColumnDecreasing, Menu } from "@lucide/svelte";
  import { onMount } from "svelte";
  import ServerResourcesBar from "./ServerResourcesBar.svelte";
  import ServerResourcesChart from "./ServerResourcesChart.svelte";

  let { id } = $props();
  type ResourceSample = {
    timestamp: Date;
    resources: components["schemas"]["GetServerResourcesResponse"];
  };
  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let information: ResourceSample[] = $state([]);
  let lastInformation:
    components["schemas"]["GetServerResourcesResponse"] | undefined =
    $state(undefined);
  let refreshFrequencySeconds = $state(1);
  let eventSource: EventSource | undefined = undefined;
  let viewMode = $state<"bar" | "chart">("bar");

  async function getServerResources() {
    if (eventSource) {
      eventSource.close();
    }
    loading = true;
    loadError = null;
    const query = new URLSearchParams({
      interval: String(refreshFrequencySeconds * 1000),
    });
    const source = new EventSource(
      `/api/server/watch-resources/${id}?${query}`,
      {
        withCredentials: true, // send cookies; omit if same-origin or no auth
      },
    );
    eventSource = source;
    addEventSourceTimeout(source, () => {
      if (source !== eventSource) return;
      eventSource = undefined;
      loading = false;
      loadError = "Timed out connecting to the resource stream";
    });
    source.onopen = () => {
      loading = false;
      console.log("connected");
    };
    source.addEventListener("data", (event) => {
      try {
        const newData = JSON.parse(
          (event as MessageEvent<string>).data,
        ) as components["schemas"]["GetServerResourcesResponse"];
        information = [
          ...information,
          { timestamp: new Date(), resources: newData },
        ].slice(-900);
        lastInformation = newData;
        loading = false;
        loadError = null;
      } catch {
        loading = false;
        loadError = "Failed to parse server resource data";
      }
    });

    source.onerror = () => console.log("disconnected, browser will retry");
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

<Card class="min-w-0">
  <CardHeader>
    <CardTitle class="flex items-center gap-1"
      ><span class="grow">Resources Usage</span>
      <Button
        variant="outline"
        onclick={() => (viewMode = viewMode === "chart" ? "bar" : "chart")}
      >
        {#if viewMode === "chart"}
          <Menu />
        {:else}
          <ChartColumnDecreasing />
        {/if}
      </Button>
      <Sort
        items={[1, 3, 5, 10, 20, 30, 40, 50]}
        bind:value={refreshFrequencySeconds}
        capitalize={false}
        onChange={getServerResources}
        prepend="s"
      />
    </CardTitle>
  </CardHeader>
  <CardContent class="space-y-6">
    {#if loading}
      <Loading loading class="my-7" />
    {:else if loadError || information === undefined}
      <ErrorSection
        message={loadError ?? "Failed to load server information"}
        onRetry={getServerResources}
      />
    {:else if viewMode === "bar"}
      <ServerResourcesBar information={lastInformation} />
    {:else if viewMode === "chart"}
      <ServerResourcesChart {information} />
    {/if}
  </CardContent>
</Card>
