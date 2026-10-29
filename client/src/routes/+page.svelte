<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import appPaths from "#lib/appPaths.ts";
  import Sort, { type SortDirection } from "#lib/components/base/sort.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import * as DropdownMenu from "#lib/components/ui/dropdown-menu/index.js";
  import Input from "#lib/components/ui/input/input.svelte";
  import Database from "@lucide/svelte/icons/database";
  import Server from "@lucide/svelte/icons/server";
  let sortDirection: SortDirection = $state("asc");
  let sortValue: "Name" | "Date" | "Size" = $state("Name");
  $effect(() => {
    async function fetchData() {
      const { data, error } = await NewRequestClient().GET("/api/files/health");
    }
    fetchData();
  });
</script>

<svelte:head>
  <title>Home | ServerGram</title>
  <meta
    name="description"
    content="Manage your server infrastructure from one place with ServerGram."
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
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        <Button>Add Service</Button>
      </DropdownMenu.Trigger>
      <DropdownMenu.Content>
        <DropdownMenu.Group>
          <DropdownMenu.Item>
            <a href={appPaths.servers.add} class="gap-2 flex items-center">
              <Server />
              Server
            </a>
          </DropdownMenu.Item>
          <DropdownMenu.Item>
            <a href={appPaths.databases.add} class="gap-2 flex items-center">
              <Database />
              Database
            </a>
          </DropdownMenu.Item>
        </DropdownMenu.Group>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  </div>
</div>
