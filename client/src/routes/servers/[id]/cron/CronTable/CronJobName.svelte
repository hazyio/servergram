<script lang="ts">
  import appPaths from "#lib/appPaths.ts";
  import { cronJobDisplayName, type CronJob } from "../cronJobs.ts";

  let { job, serverId }: { job: CronJob; serverId: number } = $props();

  const displayName = $derived(cronJobDisplayName(job));
  const unnamed = $derived(!job.name.trim());
  const href = $derived(appPaths.servers.cronJob(serverId, job.id));
  const hasHistory = $derived(job.status !== "unmanaged");
</script>

<div class="flex min-w-0 flex-col gap-0.5 py-1" title={job.command}>
  {#if hasHistory}
    <a
      {href}
      class="truncate font-medium underline-offset-4 hover:underline"
      class:text-muted-foreground={unnamed}
      class:italic={unnamed}
      title={`View runs for ${displayName}`}
    >
      {displayName}
    </a>
  {:else}
    <span
      class="truncate font-medium"
      class:text-muted-foreground={unnamed}
      class:italic={unnamed}
      title="Wrap this job with the agent to see its run history"
    >
      {displayName}
    </span>
  {/if}
  {#if job.command}
    <span class="truncate font-mono text-xs text-muted-foreground"
      >{job.command}</span
    >
  {/if}
</div>
