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
  import { CircleAlert, Loader } from "@lucide/svelte";
  import { onDestroy, onMount } from "svelte";
  import type { PageProps } from "./$types";
  import AddNewJob from "./AddNewJob.svelte";
  import CronSettingOption, {
    type CronColumns,
  } from "./CronSettingOption.svelte";
  import CronTable from "./CronTable/CronTable.svelte";
  import CronWrapUnmanagedDialog from "./CronWrapUnmanagedDialog.svelte";
  import { type CronJob, type CronJobStatus } from "./cronJobs.ts";

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
  let columns = $state<Record<CronColumns, boolean>>({
    status: true,
    job: true,
    schedule: true,
    lastRun: true,
    nextRun: true,
  });

  const totalJobs = $derived(jobs.length);
  const activeJobs = $derived(
    jobs.filter((job) => job.status !== "paused" && job.status !== "unmanaged")
      .length,
  );
  const failedJobs = $derived(
    jobs.filter((job) => job.status === "failed").length,
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

  const knownStatuses = new Set<CronJobStatus>([
    "running",
    "scheduled",
    "paused",
    "failed",
  ]);

  function toCronJob(cron: components["schemas"]["Cron"]): CronJob {
    const status: CronJobStatus = !cron.managed
      ? "unmanaged"
      : knownStatuses.has(cron.status as CronJobStatus)
        ? (cron.status as CronJobStatus)
        : "scheduled";

    return {
      id: cron.id ?? 0,
      name: cron.name ?? "",
      command: cron.command ?? "",
      schedule: cron.schedule ?? "",
      status,
      last_run: cron.last_run_at || null,
      next_run: cron.next_run_at || null,
    };
  }

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

    const source = new EventSource(`/api/cron/${serverId}/list`, {
      withCredentials: true,
    });
    eventSource = source;

    source.addEventListener("data", (event) => {
      if (source !== eventSource) return;

      try {
        const chunk = JSON.parse(
          (event as MessageEvent<string>).data,
        ) as CronChunk;

        if (chunk.jobs?.length) {
          const byId = new Map(jobs.map((job) => [job.id, job]));
          for (const cron of chunk.jobs) {
            const job = toCronJob(cron);
            byId.set(job.id, job);
          }
          jobs = [...byId.values()];
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

<CronWrapUnmanagedDialog
  bind:open={wrapDialogOpen}
  {serverId}
  jobs={unmanagedJobs}
  onWrapped={(wrapped) => {
    const updated = new Map(wrapped.map((job) => [job.id, job]));
    jobs = jobs.map((job) => updated.get(job.id) ?? job);
  }}
/>
