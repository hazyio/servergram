<script lang="ts">
  import Sort from "#lib/components/base/sort.svelte";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { scaleUtc } from "d3-scale";
  import { curveNatural } from "d3-shape";
  import { AreaChart } from "layerchart";

  let {
    information,
  }: {
    information: {
      timestamp: Date;
      resources: components["schemas"]["GetServerResourcesResponse"];
    }[];
  } = $props();

  type Range = "1m" | "5m" | "15m";
  let timeRange: Range = $state("5m");
  const rangeMilliseconds: Record<Range, number> = {
    "1m": 60_000,
    "5m": 300_000,
    "15m": 900_000,
  };
  const dateFormatter = new Intl.DateTimeFormat(undefined, {
    hour: "numeric",
    minute: "2-digit",
  });
  const series = [
    { key: "cpu", label: "CPU", color: "var(--chart-1)" },
    { key: "memory", label: "Memory", color: "var(--chart-2)" },
    { key: "disk", label: "Disk", color: "var(--chart-3)" },
    { key: "swap", label: "Swap", color: "var(--chart-4)" },
  ];

  let chartData = $derived.by(() => {
    const cutoff = Date.now() - rangeMilliseconds[timeRange];
    return information
      .filter((sample) => sample.timestamp.getTime() >= cutoff)
      .map(({ timestamp, resources }) => ({
        timestamp,
        cpu: resources.cpu?.usage,
        memory: resources.memory?.percent,
        disk: resources.disk?.percent,
        swap:
          resources.swap && (resources.swap.total ?? 0) > 0
            ? resources.swap.percent
            : undefined,
      }));
  });
</script>

<div class="space-y-3">
  <div class="flex items-center justify-between gap-2">
    <span class="grow"></span>
    <Sort
      items={["1m", "5m", "15m"]}
      bind:value={timeRange}
      capitalize={false}
      caption="Chart time range"
      
    />
  </div>
  {#if chartData.length < 2}
    <div
      class="flex h-60 items-center justify-center rounded-md border border-dashed text-sm text-muted-foreground"
    >
      Collecting resource history...
    </div>
  {:else}
    <AreaChart
      data={chartData}
      x="timestamp"
      xScale={scaleUtc()}
      yDomain={[0, 100]}
      {series}
      seriesLayout="overlap"
      height={240}
      legend
      props={{
        area: {
          curve: curveNatural,
          fillOpacity: 0.3,
          line: { class: "stroke-1" },
        },
        xAxis: { format: (value: Date) => dateFormatter.format(value) },
        yAxis: {
          ticks: [0, 25, 50, 75, 100],
          format: (value: number) => `${value}%`,
        },
      }}
    />
  {/if}
</div>
