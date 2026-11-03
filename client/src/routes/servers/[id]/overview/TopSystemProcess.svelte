<script lang="ts">
  import Sort from "#lib/components/base/sort.svelte";
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import * as Table from "#lib/components/ui/table/index.js";
  import { FlexRender } from "@tanstack/svelte-table";
  import { createTable } from "@tanstack/svelte-table";
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import { getColumns } from "./topResourcseColumns";

  import { CardContent, CardHeader } from "#lib/components/ui/card/index.ts";
  import Button from "#lib/components/ui/button/button.svelte";
  import CardFooter from "#lib/components/ui/card/card-footer.svelte";
  import appPaths from "#lib/appPaths.ts";
  type SortBy = "CPU" | "Memory" | "Disk";
  let { id } = $props();
  let sortBy: SortBy = $state("CPU");
  const columns = getColumns();
  const processes: {
    pid: number;
    name: string;
    cpu_usage: number;
    memory_usage: number;
    disk_usage: number;
  }[] = [];
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return processes;
    },
    columns,
  });
</script>

<Card class="flex min-w-0 xl:col-span-4">
  <CardHeader>
    <CardTitle class="flex items-center"
      ><span class="grow">Top System Process</span>
      <Sort items={["CPU", "Memory", "Disk"]} bind:value={sortBy} />
    </CardTitle>
  </CardHeader>

  <CardContent class="grow flex flex-col p-0">
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
              {processes.length
                ? "No matching processes."
                : "No processes yet."}
            </Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  </CardContent>
  <CardFooter class="justify-end">
  <a href={appPaths.servers.process(id)}>
    <Button variant="outline">See All</Button>
  </a>
  </CardFooter>
</Card>
