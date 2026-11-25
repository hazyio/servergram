<script lang="ts">
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Table from "#lib/components/ui/table/index.js";
  import { FileText, Square } from "@lucide/svelte";
  import type { CronRun } from "../cronJobs.ts";
  import CronRunLogDialog from "./CronRunLogDialog.svelte";
  import CronRunStatusBadge from "./CronRunStatusBadge.svelte";
  import StopCronRunDialog from "./StopCronRunDialog.svelte";

  let {
    runs,
    serverId,
    jobName,
  }: { runs: CronRun[]; serverId: number; jobName: string } = $props();

  let selectedRun = $state<CronRun | null>(null);
  let logOpen = $state(false);
  let stopOpen = $state(false);
  let runToStop = $state<CronRun | null>(null);

  function openLog(run: CronRun) {
    selectedRun = run;
    logOpen = true;
  }

  function openStop(run: CronRun) {
    runToStop = run;
    stopOpen = true;
  }

  function onRunStopped(runId: number) {
    runs = runs.map((r) =>
      r.id === runId ? { ...r, status: "canceled" as const } : r,
    );
  }

  function duration(run: CronRun) {
    if (!run.finished_at) return "Running";

    const ms = Date.parse(run.finished_at) - Date.parse(run.started_at);
    if (!Number.isFinite(ms) || ms < 0) return "—";
    if (ms < 1_000) return `${ms} ms`;

    const totalSeconds = Math.round(ms / 1_000);
    if (totalSeconds < 60) return `${totalSeconds}s`;

    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes}m ${seconds}s`;
  }

  const isRunning = (run: CronRun) => run.status === "running";
</script>

<Table.Root>
  <Table.Header>
    <Table.Row>
      <Table.Head class="h-10 px-2 md:px-4">Run</Table.Head>
      <Table.Head class="h-10 px-2 md:px-4">Status</Table.Head>
      <Table.Head class="h-10 px-2 md:px-4">Started</Table.Head>
      <Table.Head class="h-10 px-2 md:px-4">Duration</Table.Head>
      <Table.Head class="h-10 px-2 md:px-4">Trigger</Table.Head>
      <Table.Head class="h-10 px-2 md:px-4">Exit code</Table.Head>
      <Table.Head class="h-10 px-2 text-right md:px-4">Actions</Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#each runs as run (run.id)}
      <Table.Row class="h-12">
        <Table.Cell
          class="px-2 font-mono text-xs text-muted-foreground md:px-4"
        >
          #{run.id}
        </Table.Cell>
        <Table.Cell class="px-2 md:px-4">
          <CronRunStatusBadge status={run.status} />
        </Table.Cell>
        <Table.Cell class="px-2 md:px-4">
          <TimeAgo datetime={run.started_at} />
        </Table.Cell>
        <Table.Cell class="px-2 md:px-4">{duration(run)}</Table.Cell>
        <Table.Cell class="px-2 md:px-4">
          {run.trigger === "manual" ? "Manual" : "Schedule"}
        </Table.Cell>
        <Table.Cell class="px-2 font-mono text-xs md:px-4">
          {#if run.exit_code === null}
            <span class="text-muted-foreground">—</span>
          {:else}
            <span class={run.exit_code === 0 ? "" : "text-destructive"}
              >{run.exit_code}</span
            >
          {/if}
        </Table.Cell>
        <Table.Cell class="px-2 text-right md:px-4">
          <div class="flex items-center justify-end gap-1">
            <Button
              variant="ghost"
              size="sm"
              onclick={() => openLog(run)}
              class="gap-1.5"
            >
              <FileText class="size-3.5" />
              View logs
            </Button>
            {#if isRunning(run)}
              <Button
                variant="ghost"
                size="sm"
                onclick={() => openStop(run)}
                class="gap-1.5 text-destructive hover:text-destructive"
                title="Stop run"
              >
                <Square class="size-3.5" />
                <span class="hidden sm:inline">Stop</span>
              </Button>
            {/if}
          </div>
        </Table.Cell>
      </Table.Row>
    {:else}
      <Table.Row>
        <Table.Cell colspan={7} class="h-24 text-center text-muted-foreground">
          No runs recorded yet.
        </Table.Cell>
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>

<CronRunLogDialog bind:open={logOpen} {serverId} {jobName} run={selectedRun} />
<StopCronRunDialog
  bind:open={stopOpen}
  {serverId}
  run={runToStop}
  onStopped={onRunStopped}
/>
