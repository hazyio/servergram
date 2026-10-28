<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import CreateServer from "#lib/components/base/servers/createServer.svelte";
  import Sort, { type SortDirection } from "#lib/components/base/sort.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  let sortDirection: SortDirection = $state("asc");
  let sortValue: "Name" | "Date" | "Size" = $state("Name");
  $effect(() => {
    async function fetchData() {
      const { data, error } = await NewRequestClient().GET("/api/files/health");
    }
    fetchData();
  });
</script>
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
      <CreateServer />
    </div>
  </div>
