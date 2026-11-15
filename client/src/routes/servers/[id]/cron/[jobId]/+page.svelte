<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import appPaths from "#lib/appPaths.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";
  import {
    Alert,
    AlertDescription,
    AlertTitle,
  } from "#lib/components/ui/alert/index.js";
  import { Button } from "#lib/components/ui/button/index.ts";
  import Card from "#lib/components/ui/card/card.svelte";
  import { updateCronJobPageInfo } from "#lib/states/severPageInfo.svelte.ts";
  import { ArrowLeft, CircleAlert } from "@lucide/svelte";
  import type { PageProps } from "./$types";
  import CronRunsTable from "./CronRunsTable.svelte";
  import CronSchedule from "../CronTable/CronSchedule.svelte";
  import CronStatusBadge from "../CronTable/CronStatusBadge.svelte";
  import {
    cronJobDisplayName,
    getCronJob,
    listCronRuns,
    type CronJob,
    type CronRun,
  } from "../cronJobs.ts";

  let { params }: PageProps = $props();
  const serverId = $derived(Number(params.id));
  const jobId = $derived(Number(params.jobId));

  let job = $state<CronJob | null>(null);
  let runs = $state<CronRun[]>([]);
  let loading = $state(true);
  let loadError = $state<string | null>(null);

  const displayName = $derived(job ? cronJobDisplayName(job) : "Cron job");

  async function load() {
    loading = true;
    loadError = null;
    job = null;
    runs = [];

    const [jobResponse, runsResponse] = await Promise.all([
      getCronJob(serverId, jobId),
      listCronRuns(serverId, jobId),
    ]);

    loading = false;

    if (jobResponse.error || !jobResponse.data) {
      loadError = jobResponse.error
        ? unWrapErrorMessage(jobResponse.error, "Failed to load cron job")
        : "Failed to load cron job";
      return;
    }

    if (runsResponse.error || !runsResponse.data) {
      loadError = runsResponse.error
        ? unWrapErrorMessage(runsResponse.error, "Failed to load cron runs")
        : "Failed to load cron runs";
      return;
    }

    job = jobResponse.data;
    runs = runsResponse.data.items;
  }

  $effect(() => {
    load();
  });

  $effect(() => {
    updateCronJobPageInfo(jobId, job ? cronJobDisplayName(job) : "Job runs");
  });
</script>

<svelte:head>
  <title>{displayName} · Runs | ServerGram</title>
  <meta
    name="description"
    content="View the execution history and status of a cron job on your server."
  />
</svelte:head>

<div class="p-2">
  {#if loading}
    <Loading loading class="mt-20 w-12 h-12" text="Loading job runs..." />
  {:else if loadError || !job}
    <ErrorSection
      title="Failed to load cron job"
      message={loadError ?? "Cron job not found"}
      onRetry={load}
      class="mt-10"
    />
  {:else}
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div class="min-w-0">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="truncate text-lg font-semibold">{displayName}</h1>
          <CronStatusBadge status={job.status} />
        </div>
        <p class="mt-1 truncate font-mono text-xs text-muted-foreground">
          {job.command}
        </p>
      </div>
      <Button variant="outline" href={appPaths.servers.cron(serverId)}>
        <ArrowLeft />
        Back to cron jobs
      </Button>
    </div>

    <div class="mt-4 grid gap-3 sm:grid-cols-3">
      <Card class="p-4">
        <p class="text-xs text-muted-foreground">Schedule</p>
        <div class="mt-1">
          <CronSchedule schedule={job.schedule} />
        </div>
      </Card>
      <Card class="p-4">
        <p class="text-xs text-muted-foreground">Last run</p>
        <p class="mt-1 text-sm">
          {#if job.last_run}
            <TimeAgo datetime={job.last_run} />
          {:else}
            <span class="text-muted-foreground">Never</span>
          {/if}
        </p>
      </Card>
      <Card class="p-4">
        <p class="text-xs text-muted-foreground">Next run</p>
        <p class="mt-1 text-sm">
          {#if job.next_run}
            <TimeAgo datetime={job.next_run} />
          {:else}
            <span class="text-muted-foreground">Not scheduled</span>
          {/if}
        </p>
      </Card>
    </div>

    {#if job.status === "unmanaged"}
      <Alert class="mt-4 border-amber-500/40 bg-amber-500/5">
        <CircleAlert class="text-amber-600 dark:text-amber-400" />
        <AlertTitle>Runs are not being recorded</AlertTitle>
        <AlertDescription>
          This job is not managed by the agent, so there is no run history. Wrap
          it with the agent from the cron jobs page to collect logging, status,
          and run history.
        </AlertDescription>
      </Alert>
    {/if}

    <Card class="mt-4 p-4">
      <div class="flex items-center justify-between gap-2">
        <h2 class="text-sm font-medium">Runs</h2>
        <span class="text-xs text-muted-foreground">
          {runs.length}
          {runs.length === 1 ? "run" : "runs"}
        </span>
      </div>
      <CronRunsTable {runs} {serverId} {job} />
    </Card>
  {/if}
</div>
