<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import { ConnectedServices } from "#lib/components/base/services/services.svelte.ts";
  import Sort from "#lib/components/base/sort.svelte";
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";
  import {
    Alert,
    AlertAction,
    AlertDescription,
    AlertTitle,
  } from "#lib/components/ui/alert/index.js";
  import { Button } from "#lib/components/ui/button/index.ts";
  import { Input } from "#lib/components/ui/input/index.ts";
  import * as Card from "#lib/components/ui/card/index.js";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { addEventSourceTimeout } from "#lib/utils.ts";
  import { CircleAlert, Loader, RefreshCcw } from "@lucide/svelte";
  import { onDestroy, onMount } from "svelte";
  import type { PageProps } from "./$types";
  import AddNewJob from "./AddNewJob.svelte";
  import CronSettingOption, {
    type CronColumns,
  } from "./CronSettingOption.svelte";
  import CronTable from "./CronTable/CronTable.svelte";
  import CronDuplicateDialog from "./CronDuplicateDialog.svelte";
  import CronWrapUnmanagedDialog from "./CronWrapUnmanagedDialog.svelte";
  import {
    cronJobFromApi,
    type CronDuplicateGroup,
    type CronJob,
  } from "./cronJobs.ts";

  let { params }: PageProps = $props();
  let server = $derived.by(() => {
    return ConnectedServices.findServer(Number(params.id));
  });
  let serverId = $derived(Number(params.id));
  type SortStatueType = "all" | "active" | "inactive" | "unmanaged";
  let sortStatus = $state<SortStatueType>("all");
  let search = $state("");
  let jobs = $state<CronJob[]>([]);
  let loading = $state(true);
  let loadError = $state<string | null>(null);
  let updating = $state(false);
  let eventSource: EventSource | undefined;
  let wrapDialogOpen = $state(false);
  let duplicateDialogOpen = $state(false);
  let duplicateGroups = $state<CronDuplicateGroup[]>([]);
  let columns = $state<Record<CronColumns, boolean>>({
    status: true,
    job: true,
    schedule: true,
    lastRun: true,
    nextRun: true,
  });

  const totalJobs = $derived(jobs.length);
  const activeJobs = $derived(
    jobs.filter((job) => job.status === "active").length,
  );
  const inactiveJobs = $derived(
    jobs.filter((job) => job.status === "inactive").length,
  );
  const unmanagedJobs = $derived(
    jobs.filter((job) => job.status === "unmanaged"),
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

  type CronChunk = components["schemas"]["ListCronJobsResponse"];

  function stop() {
    eventSource?.close();
    eventSource = undefined;
    updating = false;
  }

  function loadJobs() {
    stop();
    loading = true;
    loadError = null;
    jobs = [];
    duplicateGroups = [];

    const hashCounts = new Map<string, CronDuplicateGroup>();
    let nextKey = 0;

    const source = new EventSource(`/api/cron/${serverId}/list`, {
      withCredentials: true,
    });
    eventSource = source;
    addEventSourceTimeout(source, () => {
      if (source !== eventSource) return;
      stop();
      loading = false;
      loadError = "Timed out connecting to the cron job stream";
    });

    source.addEventListener("data", (event) => {
      if (source !== eventSource) return;

      try {
        const chunk = JSON.parse(
          (event as MessageEvent<string>).data,
        ) as CronChunk;

        if (chunk.jobs?.length) {
          const incoming: CronJob[] = [];

          for (const cron of chunk.jobs) {
            const job = cronJobFromApi(cron);
            job.key = `${job.id}:${job.hash || "no-hash"}:${nextKey++}`;
            incoming.push(job);

            if (job.hash) {
              const group = hashCounts.get(job.hash);
              if (group) {
                group.count += 1;
              } else {
                hashCounts.set(job.hash, {
                  hash: job.hash,
                  job,
                  count: 1,
                });
              }
            }
          }

          const duplicateHashes = new Set(
            [...hashCounts.values()]
              .filter((group) => group.count > 1)
              .map((group) => group.hash),
          );

          jobs = [...jobs, ...incoming].map((job) => ({
            ...job,
            duplicate: !!job.hash && duplicateHashes.has(job.hash),
          }));

          duplicateGroups = [...hashCounts.values()].filter(
            (group) => group.count > 1,
          );
        }

        loading = false;

        if (chunk.done) {
          stop();
        } else {
          updating = true;
        }
      } catch {
        stop();
        loading = false;
        loadError = "Failed to parse cron jobs";
      }
    });

    // Server-sent "error" events carry data; connection errors do not.
    source.addEventListener("error", (event) => {
      if (source !== eventSource || !(event instanceof MessageEvent)) return;

      try {
        const response = JSON.parse(event.data) as { message?: string };
        loadError = response.message ?? "Failed to load cron jobs";
      } catch {
        loadError = "Failed to load cron jobs";
      }

      stop();
      loading = false;
    });

    source.onerror = () => {
      if (source !== eventSource || source.readyState !== EventSource.CLOSED) {
        return;
      }

      stop();
      loading = false;
      if (!loadError) loadError = "Cron stream closed unexpectedly";
    };
  }

  onMount(loadJobs);
  onDestroy(stop);
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
        <Card.Title class="text-foreground/30">Inactive</Card.Title>
      </Card.Header>
      <Card.Content>
        <p class="text-muted-foreground text-lg">{inactiveJobs}</p>
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
    <Button variant="outline" onclick={loadJobs}><RefreshCcw /></Button>

    <Sort
      items={["all", "active", "inactive", "unmanaged"]}
      bind:value={sortStatus}
      capitalize={true}
    />
    <CronSettingOption bind:columns />
    <AddNewJob {serverId} />
  </div>

  {#if !loading && !loadError && unmanagedJobs.length > 0}
    <Alert class="mt-3 border-amber-500/40 bg-amber-500/5">
      <CircleAlert class="text-amber-600 dark:text-amber-400" />
      <AlertTitle>Unmanaged cron jobs</AlertTitle>
      <AlertDescription>
        {unmanagedJobs.length} cron
        {unmanagedJobs.length === 1 ? "job is" : "jobs are"} not managed by the agent.
        Wrap {unmanagedJobs.length === 1 ? "it" : "them"} to get logging, status tracking,
        and run history.
      </AlertDescription>
      <AlertAction>
        <Button size="sm" onclick={() => (wrapDialogOpen = true)}>
          Fix now
        </Button>
      </AlertAction>
    </Alert>
  {/if}

  {#if !loading && !loadError && duplicateGroups.length > 0}
    <Alert class="mt-3 border-amber-500/40 bg-amber-500/5">
      <CircleAlert class="text-amber-600 dark:text-amber-400" />
      <AlertTitle>Duplicate cron jobs</AlertTitle>
      <AlertDescription>
        {duplicateGroups.length} cron
        {duplicateGroups.length === 1 ? "job has" : "jobs have"} duplicate entries
        with the same hash and may run more than once. They are highlighted below.
      </AlertDescription>
      <AlertAction>
        <Button size="sm" onclick={() => (duplicateDialogOpen = true)}>
          Fix now
        </Button>
      </AlertAction>
    </Alert>
  {/if}

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
    {#if updating}
      <div
        class="mt-2 flex items-center gap-2 px-2 text-sm text-muted-foreground"
      >
        <Loader class="size-4 animate-spin" />
        Updating...
      </div>
    {/if}

    <div class="px-2">
      <CronTable
        {jobs}
        {serverId}
        {search}
        statusFilter={sortStatus}
        visibleColumns={columns}
        onUpdated={(updated) => {
          jobs = jobs.map((job) =>
            job.id === updated.id
              ? {
                  ...updated,
                  key: job.key,
                  hash: job.hash ?? updated.hash,
                  duplicate: job.duplicate,
                }
              : job,
          );
        }}
        onDeleted={(id) => {
          jobs = jobs.filter((job) => job.id !== id);
        }}
      />
    </div>
  {/if}
</div>

<CronWrapUnmanagedDialog
  bind:open={wrapDialogOpen}
  {serverId}
  jobs={unmanagedJobs}
  onWrapped={() => loadJobs()}
/>

<CronDuplicateDialog
  bind:open={duplicateDialogOpen}
  {serverId}
  duplicates={duplicateGroups}
  onRemoved={() => {
    duplicateGroups = [];
    loadJobs();
  }}
/>
