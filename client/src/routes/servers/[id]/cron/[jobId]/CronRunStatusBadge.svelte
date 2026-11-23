<script lang="ts">
  import { Badge, type BadgeVariant } from "#lib/components/ui/badge/index.ts";
  import type { CronRunStatus } from "../cronJobs.ts";

  let { status }: { status: CronRunStatus } = $props();

  const statuses: Record<
    CronRunStatus,
    { label: string; variant: BadgeVariant; class: string }
  > = {
    running: {
      label: "Running",
      variant: "secondary",
      class: "text-sky-600 dark:text-sky-400",
    },
    succeeded: {
      label: "Succeeded",
      variant: "outline",
      class: "text-emerald-600 dark:text-emerald-400",
    },
    failed: {
      label: "Failed",
      variant: "destructive",
      class: "",
    },
    canceled: {
      label: "Canceled",
      variant: "outline",
      class: "text-muted-foreground",
    },
    timed_out: {
      label: "Timed out",
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
