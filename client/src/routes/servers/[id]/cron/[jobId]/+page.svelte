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
  import { ArrowLeft, CircleAlert, Settings2 } from "@lucide/svelte";
  import type { PageProps } from "./$types";
  import CronRunsTable from "./CronRunsTable.svelte";
  import CronEditDialog from "../CronEditDialog.svelte";
  import CronSchedule from "../CronTable/CronSchedule.svelte";
  import CronStatusBadge from "../CronTable/CronStatusBadge.svelte";
  import {
    cronJobDisplayName,
    getCronJobInfo,
    listCronRuns,
    type CronJobInfo,
    type CronRun,
  } from "../cronJobs.ts";

  let { params }: PageProps = $props();
  const serverId = $derived(Number(params.id));
  const jobId = $derived(Number(params.jobId));

  let info = $state<CronJobInfo | null>(null);
  let infoLoading = $state(true);
  let infoError = $state<string | null>(null);

  let runs = $state<CronRun[]>([]);
  let runsLoading = $state(false);
  let runsError = $state<string | null>(null);

  let editOpen = $state(false);

  let infoToken = 0;
  let runsToken = 0;

  const displayName = $derived(info ? cronJobDisplayName(info) : "Cron job");

  async function loadRuns() {
    const token = ++runsToken;

    runsLoading = true;
    runsError = null;
    runs = [];

    const { data, error } = await listCronRuns(serverId, jobId);

    if (token !== runsToken) return;

    runsLoading = false;

    if (error || !data) {
      runsError = error
        ? unWrapErrorMessage(error, "Failed to load cron runs")
        : "Failed to load cron runs";
      return;
    }

    runs = data.items;
  }

  async function loadInfo() {
    const token = ++infoToken;
    runsToken++;

    infoLoading = true;
    infoError = null;
    info = null;

    runsLoading = false;
    runsError = null;
    runs = [];

    const { data, error } = await getCronJobInfo(serverId, jobId);

    if (token !== infoToken) return;

    infoLoading = false;

    if (error || !data) {
      infoError = error
        ? unWrapErrorMessage(error, "Failed to load cron job")
        : "Failed to load cron job";
      return;
    }

    info = data;

    if (data.managed) {
      void loadRuns();
    }
  }

  $effect(() => {
    loadInfo();
  });

  $effect(() => {
    updateCronJobPageInfo(jobId, info ? cronJobDisplayName(info) : "Job runs");
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
  {#if infoLoading}
    <Loading loading class="mt-20 w-12 h-12" text="Loading cron job..." />
  {:else if infoError || !info}
    <ErrorSection
      title="Failed to load cron job"
      message={infoError ?? "Cron job not found"}
      onRetry={loadInfo}
      class="mt-10"
    />
  {:else}
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div class="min-w-0">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="truncate text-lg font-semibold">{displayName}</h1>
          <CronStatusBadge status={info.status} />
        </div>
        <p class="mt-1 truncate font-mono text-xs text-muted-foreground">
          {info.command}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <Button variant="outline" onclick={() => (editOpen = true)}>
          <Settings2 />
          Settings
        </Button>
        <Button variant="outline" href={appPaths.servers.cron(serverId)}>
          <ArrowLeft />
          Back to cron jobs
        </Button>
      </div>
    </div>

    <div class="mt-4 grid gap-3 sm:grid-cols-3">
      <Card class="p-4">
        <p class="text-xs text-muted-foreground">Schedule</p>
        <div class="mt-1">
          <CronSchedule schedule={info.schedule} />
        </div>
      </Card>
      <Card class="p-4">
        <p class="text-xs text-muted-foreground">Last run</p>
        <p class="mt-1 text-sm">
          {#if info.last_run}
            <TimeAgo datetime={info.last_run} />
          {:else}
            <span class="text-muted-foreground">Never</span>
          {/if}
        </p>
      </Card>
      <Card class="p-4">
        <p class="text-xs text-muted-foreground">Next run</p>
        <p class="mt-1 text-sm">
          {#if info.next_run}
            <TimeAgo datetime={info.next_run} />
          {:else}
            <span class="text-muted-foreground">Not scheduled</span>
          {/if}
        </p>
      </Card>
    </div>

    {#if info.managed}
      <Card class="mt-4 p-4">
        <h2 class="text-sm font-medium">Retention</h2>
        <div class="mt-3 grid gap-3 sm:grid-cols-3">
          <div>
            <p class="text-xs text-muted-foreground">Runs metadata to keep</p>
            <p class="mt-1 text-sm">
              {info.retention.metadata_runs}
              {info.retention.metadata_runs === 1 ? "run" : "runs"}
            </p>
          </div>
          <div>
            <p class="text-xs text-muted-foreground">Run logs to keep</p>
            <p class="mt-1 text-sm">
              {info.retention.run_logs}
              {info.retention.run_logs === 1 ? "run" : "runs"}
            </p>
          </div>
          <div>
            <p class="text-xs text-muted-foreground">Error logs to keep</p>
            <p class="mt-1 text-sm">
              {info.retention.error_log_days}
              {info.retention.error_log_days === 1 ? "day" : "days"}
            </p>
          </div>
        </div>
      </Card>
    {/if}

    {#if info.status === "unmanaged"}
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

    {#if info.managed}
      <Card class="mt-4 p-4">
        <div class="flex items-center justify-between gap-2">
          <h2 class="text-sm font-medium">Runs</h2>
          {#if !runsLoading && !runsError}
            <span class="text-xs text-muted-foreground">
              {runs.length}
              {runs.length === 1 ? "run" : "runs"}
            </span>
          {/if}
        </div>

        {#if runsLoading}
          <Loading
            loading
            containerClass="h-32"
            class="w-8 h-8"
            text="Loading runs..."
          />
        {:else if runsError}
          <ErrorSection
            title="Failed to load run history"
            message={runsError}
            onRetry={loadRuns}
            class="mt-3"
          />
        {:else}
          <CronRunsTable {runs} {serverId} job={info} />
        {/if}
      </Card>
    {/if}
  {/if}
</div>

{#if info}
  <CronEditDialog
    bind:open={editOpen}
    {serverId}
    job={info}
    onUpdated={(updated) => {
      if (!info) return;

      info = {
        ...info,
        name: updated.name,
        command: updated.command,
        schedule: updated.schedule,
        retention: updated.retention ?? info.retention,
      };
    }}
  />
{/if}
