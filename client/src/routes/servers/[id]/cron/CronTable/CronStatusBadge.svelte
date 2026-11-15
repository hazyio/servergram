<script lang="ts">
  import { Badge, type BadgeVariant } from "#lib/components/ui/badge/index.ts";
  import type { CronJobStatus } from "../cronJobs.ts";

  let { status }: { status: CronJobStatus } = $props();

  const statuses: Record<
    CronJobStatus,
    { label: string; variant: BadgeVariant; class: string }
  > = {
    running: {
      label: "Running",
      variant: "secondary",
      class: "text-emerald-600 dark:text-emerald-400",
    },
    scheduled: {
      label: "Scheduled",
      variant: "outline",
      class: "text-sky-600 dark:text-sky-400",
    },
    paused: {
      label: "Paused",
      variant: "outline",
      class: "text-muted-foreground",
    },
    failed: {
      label: "Failed",
      variant: "destructive",
      class: "",
    },
  };

  const current = $derived(statuses[status]);
</script>

<Badge variant={current.variant} class={current.class}>
  <span class="size-1.5 rounded-full bg-current"></span>
  {current.label}
</Badge>
