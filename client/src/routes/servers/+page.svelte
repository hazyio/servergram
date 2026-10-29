<script lang="ts">
  import appPaths from "#lib/appPaths.ts";
  import Sort, { type SortDirection } from "#lib/components/base/sort.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  let sortDirection: SortDirection = $state("asc");
  let sortValue: "Name" | "Date" | "Size" = $state("Name");
</script>

<svelte:head>
  <title>Servers | ServerGram</title>
  <meta
    name="description"
    content="View and manage the servers connected to your ServerGram account."
  />
</svelte:head>

<div class="p-2">
  <div class="flex gap-2">
    <Input placeholder="Search..." />
    <Sort
      value={sortValue}
      items={["Name", "Date", "Size"]}
      direction={sortDirection}
      onSelect={(value, direction) => {
        console.log(`Selected sort: ${value} ${direction}`);
        sortValue = value;
        sortDirection = direction;
      }}
    />
    <Button
      ><a href={appPaths.servers.add} class="gap-2 flex items-center">
        Add Server
      </a></Button
    >
  </div>
</div>
