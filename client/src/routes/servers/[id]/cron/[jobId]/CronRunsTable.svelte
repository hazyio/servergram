<script lang="ts">
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Table from "#lib/components/ui/table/index.js";
  import { FileText } from "@lucide/svelte";
  import type { CronJobInfo, CronRun } from "../cronJobs.ts";
  import CronRunLogDialog from "./CronRunLogDialog.svelte";
  import CronRunStatusBadge from "./CronRunStatusBadge.svelte";

  let {
    runs,
    serverId,
    job,
  }: { runs: CronRun[]; serverId: number; job: CronJobInfo } = $props();

  let logsOpen = $state(false);
  let selectedRun = $state<CronRun | null>(null);

  function viewLogs(run: CronRun) {
    selectedRun = run;
    logsOpen = true;
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
      <Table.Head class="h-10 px-2 text-right md:px-4">Logs</Table.Head>
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
          <Button variant="ghost" size="sm" onclick={() => viewLogs(run)}>
            <FileText class="size-3.5" />
            View logs
          </Button>
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

<CronRunLogDialog bind:open={logsOpen} {serverId} {job} run={selectedRun} />
