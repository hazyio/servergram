<script lang="ts">
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import { ChevronLeft, ChevronRight, Settings } from "@lucide/svelte";
  import type { PageProps } from "./$types";
  import SettingOption from "./SettingOption.svelte";
  import * as Table from "#lib/components/ui/table/index.js";
  import { FlexRender, createTable } from "@tanstack/svelte-table";
  import { dataTableFeatures } from "#lib/components/base/data-table-features.js";
  import { getColumns, stubFiles } from "./filesColumns";

  let { data }: PageProps = $props();
  let settings = $state({ showHidden: false });

  const columns = getColumns();
  const files = $derived(
    settings.showHidden
      ? stubFiles
      : stubFiles.filter((f) => !f.name.startsWith(".")),
  );
  const table = createTable({
    features: dataTableFeatures,
    get data() {
      return files;
    },
    columns,
  });
</script>

<div class="flex gap-2">
  <div class="flex gap-2">
    <Button variant="outline"><ChevronLeft /></Button>
    <Button variant="outline"><ChevronRight /></Button>
  </div>
  <Input value="/var/www/app" class="grow-2" />
  <Input value="" placeholder="Search files here..." />
  <SettingOption bind:showHidden={settings.showHidden} />
</div>

<Table.Root class="mt-2">
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
          No files.
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
