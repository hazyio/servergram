<script lang="ts">
  import { onMount } from "svelte";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Sort from "#lib/components/base/sort.svelte";
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import * as Table from "#lib/components/ui/table/index.js";
  import { FlexRender } from "@tanstack/svelte-table";
  import { createTable } from "@tanstack/svelte-table";
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import { getColumns, type TopProcess } from "./topResourcseColumns";

  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import Button from "#lib/components/ui/button/button.svelte";
  import CardFooter from "#lib/components/ui/card/card-footer.svelte";
  import appPaths from "#lib/appPaths.ts";
  import { addEventSourceTimeout } from "#lib/utils.ts";
  type SortBy = "CPU" | "Memory";
  let { id } = $props();
  let sortBy: SortBy = $state("CPU");
  let refreshFrequencySeconds = $state(1);
  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let processes: TopProcess[] = $state([]);
  let eventSource: EventSource | undefined;
  const columns = getColumns();
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return processes;
    },
    columns,
  });

  function getTopProcesses() {
    eventSource?.close();
    loading = true;
    loadError = null;
    processes = [];

    const sortValue = sortBy === "CPU" ? 0 : 1;
    const query = new URLSearchParams({
      interval: String(refreshFrequencySeconds * 1000),
      sortBy: String(sortValue),
      limit: "5",
    });
    const source = new EventSource(
      `/api/server/watch-top-process/${id}?${query}`,
      { withCredentials: true },
    );
    eventSource = source;
    addEventSourceTimeout(source, () => {
      if (source !== eventSource) return;
      eventSource = undefined;
      loading = false;
      loadError = "Timed out connecting to the process stream";
    });

    source.addEventListener("data", (event) => {
      if (source !== eventSource) return;
      try {
        const response = JSON.parse((event as MessageEvent<string>).data) as {
          items?: TopProcess[] | null;
        };
        if (!Array.isArray(response.items)) {
          throw new Error("Invalid process data");
        }
        processes = response.items;
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
      loading = false;
    });

    source.onerror = () => {
      if (source !== eventSource || source.readyState !== EventSource.CLOSED) {
        return;
      }
      loading = false;
      loadError ??= "Process stream closed";
    };
  }

  onMount(() => {
    getTopProcesses();
    return () => eventSource?.close();
  });
</script>

<Card class="flex min-w-0">
  <CardHeader>
    <CardTitle class="flex items-center gap-2"
      ><span class="grow">Top System Process</span>
      <Sort
        items={["CPU", "Memory"]}
        bind:value={sortBy}
        onChange={getTopProcesses}
      />
      <Sort
        items={[1, 5, 10, 15, 30]}
        bind:value={refreshFrequencySeconds}
        capitalize={false}
        caption="Refresh interval"
        onChange={getTopProcesses}
        prepend="s"
      />
    </CardTitle>
  </CardHeader>

  <CardContent class="grow flex flex-col p-0">
    {#if loading}
      <Loading loading class="my-7" />
    {:else if loadError}
      <ErrorSection message={loadError} onRetry={getTopProcesses} />
    {:else}
      <Table.Root>
        <Table.Header>
          {#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
            <Table.Row>
              {#each headerGroup.headers as header (header.id)}
                <Table.Head colspan={header.colSpan}>
                  {#if !header.isPlaceholder}
                    <FlexRender {header} />
                  {/if}
                </Table.Head>
              {/each}
            </Table.Row>
          {/each}
        </Table.Header>
        <Table.Body>
          {#each table.getRowModel().rows as row (row.id)}
            <Table.Row>
              {#each row.getVisibleCells() as cell (cell.id)}
                <Table.Cell>
                  <FlexRender {cell} />
                </Table.Cell>
              {/each}
            </Table.Row>
          {:else}
            <Table.Row>
              <Table.Cell colspan={columns.length} class="h-24 text-center">
                No processes reported.
              </Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    {/if}
  </CardContent>
  <CardFooter class="justify-end">
    <a href={appPaths.servers.processes(id)}>
      <Button variant="outline">See All</Button>
    </a>
  </CardFooter>
</Card>
