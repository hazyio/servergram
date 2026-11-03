<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import { Progress } from "#lib/components/ui/progress/index.js";
  import type { components } from "#lib/types/request-schema.d.ts";

  let {
    information,
  }: {
    information: components["schemas"]["GetServerResourcesResponse"][];
  } = $props();

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
</script>

{#if information === undefined}
  <ErrorSection message={"Failed to load server information"} />
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
        {#if information.memory?.available !== undefined}
          <span>· {formatBytes(information.memory.available)} available</span>
        {/if}
      </div>
    </div>
    <Progress value={progressValue(information.memory?.percent)} />
  </div>
  <div class="space-y-2">
    <div class="flex">
      <span class="grow">Disk</span>
      <span class="text-[12px] text-muted-foreground">
        {formatPercent(information.disk?.percent)}
        ({formatUsedTotal(information.disk?.usage, information.disk?.total)})
        {#if information.disk?.available !== undefined}
          · {formatBytes(information.disk.available)} available
        {/if}
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
          ({formatUsedTotal(information.swap.usage, information.swap.total)})
          {#if information.swap.available !== undefined}
            · {formatBytes(information.swap.available)} available
          {/if}
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
