<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import { ConnectedServices } from "#lib/components/base/services/services.svelte.ts";
  import Sort from "#lib/components/base/sort.svelte";
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";
  import { Input } from "#lib/components/ui/input/index.ts";
  import * as Card from "#lib/components/ui/card/index.js";
  import { onMount } from "svelte";
  import type { PageProps } from "./$types";
  import AddNewJob from "./AddNewJob.svelte";
  import CronSettingOption, {
    type CronColumns,
  } from "./CronSettingOption.svelte";
  import CronTable from "./CronTable/CronTable.svelte";
  import { listCronJobs, type CronJob } from "./cronJobs.ts";

  let { params }: PageProps = $props();
  let server = $derived.by(() => {
    return ConnectedServices.findServer(Number(params.id));
  });
  let serverId = $derived(Number(params.id));
  type SortStatueType = "all" | "active" | "inactive";
  let sortStatus = $state<SortStatueType>("all");
  let search = $state("");
  let jobs = $state<CronJob[]>([]);
  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let columns = $state<Record<CronColumns, boolean>>({
    status: true,
    job: true,
    schedule: true,
    lastRun: true,
    nextRun: true,
  });

  const totalJobs = $derived(jobs.length);
  const activeJobs = $derived(
    jobs.filter((job) => job.status !== "paused").length,
  );
  const failedJobs = $derived(
    jobs.filter((job) => job.status === "failed").length,
  );
  const nextJob = $derived.by(() => {
    const upcoming = jobs
      .filter((job) => job.next_run)
      .sort(
        (a, b) =>
          new Date(a.next_run as string).getTime() -
          new Date(b.next_run as string).getTime(),
      );
    return upcoming[0] ?? null;
  });

  async function loadJobs() {
    loading = true;
    loadError = null;

    const { data, error } = await listCronJobs(serverId);
    loading = false;

    if (error) {
      loadError = unWrapErrorMessage(error, "Failed to load cron jobs");
      return;
    }

    jobs = data?.items ?? [];
  }

  onMount(loadJobs);
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
        <p class="text-lg">{totalJobs}</p>
      </Card.Content>
    </Card.Root>
    <Card.Root class="grow">
      <Card.Header>
        <Card.Title class="text-foreground/30">Active</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-lg">{activeJobs}</p>
      </Card.Content>
    </Card.Root>
    <Card.Root class="grow">
      <Card.Header>
        <Card.Title class="text-foreground/30">Failed</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-destructive text-lg">{failedJobs}</p>
      </Card.Content>
    </Card.Root>
    <Card.Root class="grow">
      <Card.Header>
        <Card.Title class="text-foreground/30">Next Run</Card.Title>
      </Card.Header>
      <Card.Content>
        {#if nextJob?.next_run}
          <p class="text-lg"><TimeAgo datetime={nextJob.next_run} /></p>
          <p class="truncate text-xs text-muted-foreground">{nextJob.name}</p>
        {:else}
          <p class="text-lg text-muted-foreground">Not scheduled</p>
        {/if}
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
    <CronSettingOption bind:columns />
    <AddNewJob />
  </div>

  {#if loading}
    <Loading loading class="mt-20 w-12 h-12" text="Loading cron jobs..." />
  {:else if loadError}
    <ErrorSection
      title="Failed to load cron jobs"
      message={loadError}
      onRetry={loadJobs}
      class="my-7"
    />
  {:else}
    <div class="px-2">
      <CronTable
        {jobs}
        {serverId}
        {search}
        statusFilter={sortStatus}
        visibleColumns={columns}
        onUpdated={(updated) => {
          jobs = jobs.map((job) => (job.id === updated.id ? updated : job));
        }}
        onCreated={(created) => {
          jobs = [...jobs, created];
        }}
        onDeleted={(id) => {
          jobs = jobs.filter((job) => job.id !== id);
        }}
      />
    </div>
  {/if}
</div>
