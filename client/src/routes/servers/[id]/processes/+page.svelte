<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Sort from "#lib/components/base/sort.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import Input from "#lib/components/ui/input/input.svelte";
  import { severPageInfo } from "#lib/states/severPageInfo.svelte.ts";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { addEventSourceTimeout } from "#lib/utils.ts";
  import { RefreshCcw } from "@lucide/svelte";
  import { untrack } from "svelte";
  import type { PageProps } from "./$types";
  import ProcessesTable from "./ProcessesTable/ProcessesTable.svelte";
  import {
    toProcessItem,
    type ProcessItem,
  } from "./ProcessesTable/processesColumns";
  import ProcessSettingOption, {
    type ProcessColumns,
    type ProcessGrouping,
  } from "./ProcessSettingOption.svelte";

  type WatchProcessesResponse = components["schemas"]["WatchProcessesResponse"];

  const processColumnValues: Record<ProcessColumns, number> = {
    memory: 1,
    cmd: 2,
    user: 3,
    cpu: 4,
    state: 5,
    uptime: 6,
  };

  let search = $state("");
  let refreshFrequencySeconds = $state(1);
  let processes = $state<ProcessItem[]>([]);
  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let eventSource: EventSource | undefined;

  let { params }: PageProps = $props();

  let columnsToShow = $state<Record<ProcessColumns, boolean>>({
    cmd: true,
    cpu: true,
    memory: true,
    state: true,
    uptime: true,
    user: true,
  });
  let groupBy = $state<ProcessGrouping>("parent");

  function getSelectedColumns(): number[] {
    const selected = (Object.keys(processColumnValues) as ProcessColumns[])
      .filter((column) => columnsToShow[column])
      .map((column) => processColumnValues[column]);

    // The API requires at least one column; pid and name are always returned.
    return selected.length > 0 ? selected : [processColumnValues.cmd];
  }

  function stop() {
    eventSource?.close();
    eventSource = undefined;
  }

  function startProcessWatch() {
    stop();
    loading = true;
    loadError = null;
    processes = [];

    const query = new URLSearchParams({
      interval: String(refreshFrequencySeconds * 1000),
      columns: getSelectedColumns().join(","),
    });
    const source = new EventSource(
      `/api/server/watch-processes/${params.id}?${query}`,
      { withCredentials: true },
    );
    eventSource = source;
    addEventSourceTimeout(source, () => {
      if (source !== eventSource) return;
      stop();
      loading = false;
      loadError = "Timed out connecting to the process stream";
    });

    source.addEventListener("data", (event) => {
      if (source !== eventSource) return;
      try {
        const response = JSON.parse(
          (event as MessageEvent<string>).data,
        ) as WatchProcessesResponse;
        processes = (response.items ?? []).map(toProcessItem);
        loading = false;
        loadError = null;
      } catch {
        loading = false;
        loadError = "Failed to parse process data";
      }
    });

    source.addEventListener("error", (event) => {
      if (source !== eventSource || !(event instanceof MessageEvent)) return;
      try {
        const response = JSON.parse(event.data) as { message?: string };
        loadError = response.message ?? "Failed to load processes";
      } catch {
        loadError = "Failed to load processes";
      }
      stop();
      loading = false;
    });

    source.onerror = () => {
      if (source !== eventSource || source.readyState !== EventSource.CLOSED) {
        return;
      }
      stop();
      loading = false;
      loadError ??= "Process stream closed";
    };
  }

  $effect(() => {
    // Restart the stream when the refresh interval or visible columns change.
    void refreshFrequencySeconds;
    void columnsToShow;
    untrack(startProcessWatch);
    return stop;
  });
</script>

<svelte:head>
  <title>{severPageInfo.name} Processes | ServerGram</title>
  <meta
    name="description"
    content="View and search running processes on a server, including CPU, memory, state, and uptime."
  />
</svelte:head>
<div class="flex gap-2">
  <Input
    name="filter-process"
    placeholder="Search here..."
    autocomplete="off"
    type="search"
    bind:value={search}
  />
  <ProcessSettingOption
    onSave={(settings) => {
      columnsToShow = settings.columns;
      groupBy = settings.groupBy;
    }}
  />
  <Sort
    items={[1, 3, 5, 10, 20, 30, 40, 50]}
    bind:value={refreshFrequencySeconds}
    capitalize={false}
    prepend="s"
  />
  <Button variant="outline" onclick={startProcessWatch}><RefreshCcw /></Button>
</div>

{#if loading}
  <Loading loading class="mt-20 w-12 h-12" />
{:else if loadError}
  <ErrorSection message={loadError} onRetry={startProcessWatch} class="my-7" />
{:else}
  <div class="p-2">
    <ProcessesTable
      {processes}
      {search}
      {groupBy}
      serverId={Number(params.id)}
      visibleColumns={columnsToShow}
    />
  </div>
{/if}
