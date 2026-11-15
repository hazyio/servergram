<script lang="ts">
  import { ConnectedServices } from "#lib/components/base/services/services.svelte.ts";
  import Sort from "#lib/components/base/sort.svelte";
  import { Input } from "#lib/components/ui/input/index.ts";
  import type { PageProps } from "./$types";
  import AddNewJob from "./AddNewJob.svelte";
  import * as Card from "#lib/components/ui/card/index.js";

  let { params }: PageProps = $props();
  let server = $derived.by(() => {
    return ConnectedServices.findServer(Number(params.id));
  });
  type SortStatueType = "all" | "active" | "inactive";
  let sortStatus = $state<SortStatueType>("all");
  let search = $state("");
</script>

<svelte:head>
  <title>{server?.name ?? "Cron Jobs"} | ServerGram</title>
  <meta
    name="description"
    content="Manage cron jobs on your server, including scheduling, execution history, and configuration settings."
  />
</svelte:head>
<div class="p-2">
  <div class="flex flex-row flex-wrap justify-between gap-4">
    <Card.Root class="grow">
      <Card.Header>
        <Card.Title class="text-foreground/30">Total</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-lg">6</p>
      </Card.Content>
    </Card.Root>
    <Card.Root class="grow">
      <Card.Header>
        <Card.Title class="text-foreground/30">Active</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-lg">4</p>
      </Card.Content>
    </Card.Root>
    <Card.Root class="grow">
      <Card.Header>
        <Card.Title class="text-foreground/30">Failed</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-destructive text-lg">2</p>
      </Card.Content>
    </Card.Root>
    <Card.Root class="grow">
      <Card.Header>
        <Card.Title class="text-foreground/30">Next Run</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-lg">Card Content</p>
      </Card.Content>
    </Card.Root>
  </div>
  <div class="flex gap-2 mt-3">
    <Input
      name="filter-cron-jobs"
      placeholder="Search here..."
      autocomplete="off"
      type="search"
      bind:value={search}
    />

    <Sort
      items={["all", "active", "inactive"]}
      bind:value={sortStatus}
      capitalize={true}
    />
    <AddNewJob />
  </div>
</div>
