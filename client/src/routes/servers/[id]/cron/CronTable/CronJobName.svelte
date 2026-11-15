<script lang="ts">
  import { cronJobDisplayName, type CronJob } from "../cronJobs.ts";

  let { job }: { job: CronJob } = $props();

  const displayName = $derived(cronJobDisplayName(job));
  const unnamed = $derived(!job.name.trim());
</script>

<div class="flex min-w-0 flex-col gap-0.5 py-1" title={job.command}>
  <span
    class="truncate font-medium"
    class:text-muted-foreground={unnamed}
    class:italic={unnamed}
  >
    {displayName}
  </span>
  {#if job.command}
    <span class="truncate font-mono text-xs text-muted-foreground"
      >{job.command}</span
    >
  {/if}
</div>
