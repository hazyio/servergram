<script lang="ts">
  import Sort from "#lib/components/base/sort.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import Input from "#lib/components/ui/input/input.svelte";
  import { severPageInfo } from "#lib/states/severPageInfo.svelte.ts";
  import { RefreshCcw } from "@lucide/svelte";
  import type { PageProps } from "./$types";
  import ProcessesTable from "./ProcessesTable/ProcessesTable.svelte";
  import ProcessSettingOption, {
    type ProcessColumns
  } from "./ProcessSettingOption.svelte";
  let search = $state("");
  let refreshFrequencySeconds = $state(1);

  let { params }: PageProps = $props();

  let columnsToShow = $state<Record<ProcessColumns, boolean>>({
    cmd: true,
    cpu: true,
    memory: true,
    state: true,
    uptime: true,
    user: true,
  });
  function startProcessWatch() {
    stop();
  }
  function stop() {}
</script>

<svelte:head>
  <title>{severPageInfo.name} Processes | ServerGram</title>
  <meta
    name="description"
    content="View the details of a server, including its hostname, SSH username, port, group, and SSH key."
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
    onSave={(c) => {
      console.log("saved")
      columnsToShow = c;
      startProcessWatch();
    }}
  />
  <Sort
    items={[1, 3, 5, 10, 20, 30, 40, 50]}
    bind:value={refreshFrequencySeconds}
    capitalize={false}
    onChange={startProcessWatch}
    prepend="s"
  />
  <Button variant="outline" onclick={startProcessWatch}><RefreshCcw /></Button>
</div>

<div class="p-2">
      <ProcessesTable {search} visibleColumns={columnsToShow} />
</div>
