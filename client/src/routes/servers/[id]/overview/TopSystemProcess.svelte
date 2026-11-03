<script lang="ts">
  import Sort from "#lib/components/base/sort.svelte";
  import CardTitle from "#lib/components/ui/card/card-title.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CardHeader } from "#lib/components/ui/card/index.ts";
  type SortBy = "CPU" | "Memory" | "Disk";
  let { id } = $props();
  let sortBy: SortBy = $state("CPU");
</script>

<Card class="md:col-start-3 md:col-end-6">
  <CardHeader>
    <CardTitle class="flex items-center"
      ><span class="grow">Top System Process</span>
      <Sort items={["CPU", "Memory", "Disk"]} bind:value={sortBy} />
    </CardTitle>
  </CardHeader>
</Card>
