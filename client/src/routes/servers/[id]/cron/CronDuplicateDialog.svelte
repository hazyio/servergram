<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import {
    Alert,
    AlertDescription,
    AlertTitle,
  } from "#lib/components/ui/alert/index.js";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import { CircleCheck, Trash2 } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import {
    cronJobDisplayName,
    removeDuplicateCronJobs,
    type CronDuplicateGroup,
  } from "./cronJobs.ts";

  let {
    open = $bindable(false),
    serverId,
    duplicates,
    onRemoved,
  }: {
    open: boolean;
    serverId: number;
    duplicates: CronDuplicateGroup[];
    onRemoved?: (removed: number) => void;
  } = $props();

  let submitting = $state(false);
  let error = $state("");
  let removedCount = $state(0);

  const extraCount = $derived(
    duplicates.reduce((total, group) => total + group.count - 1, 0),
  );

  async function remove() {
    if (submitting) return;

    submitting = true;
    error = "";

    const { data, error: requestError } =
      await removeDuplicateCronJobs(serverId);
    submitting = false;

    if (requestError || !data) {
      error = requestError
        ? unWrapErrorMessage(requestError, "Failed to remove duplicate jobs")
        : "Failed to remove duplicate jobs";
      return;
    }

    removedCount = data.removed;
    onRemoved?.(data.removed);
    toast.success(
      `${data.removed} duplicate cron ${data.removed === 1 ? "job" : "jobs"} removed`,
    );
  }

  function close() {
    if (submitting) return;
    open = false;
    error = "";
  }

  $effect(() => {
    if (!open) {
      submitting = false;
      error = "";
      removedCount = 0;
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content showCloseButton={false}>
    {#if submitting}
      <Dialog.Header>
        <Dialog.Title>Removing duplicate cron jobs</Dialog.Title>
        <Dialog.Description>
          Cleaning up {extraCount} duplicate cron
          {extraCount === 1 ? "entry" : "entries"} so each job runs only once.
        </Dialog.Description>
      </Dialog.Header>
      <Loading
        loading
        class="h-50"
        text={`Removing ${extraCount} duplicate cron ${extraCount === 1 ? "job" : "jobs"}...`}
      />
    {:else if error}
      <Dialog.Header>
        <Dialog.Title>Could not remove duplicate jobs</Dialog.Title>
      </Dialog.Header>
      <ErrorSection
        title="Failed to remove duplicate jobs"
        message={error}
        onRetry={remove}
      />
      <Dialog.Footer>
        <Button onclick={close}>Close</Button>
      </Dialog.Footer>
    {:else if removedCount > 0}
      <Dialog.Header>
        <Dialog.Title class="sr-only">Duplicate cron jobs removed</Dialog.Title>
      </Dialog.Header>
      <Alert class="border-emerald-500/40 bg-emerald-500/5">
        <CircleCheck class="text-emerald-600 dark:text-emerald-400" />
        <AlertTitle>Duplicate cron jobs removed</AlertTitle>
        <AlertDescription>
          {removedCount} duplicate cron {removedCount === 1
            ? "entry was"
            : "entries were"}
          removed. Each job will now run only once per schedule.
        </AlertDescription>
      </Alert>
      <Dialog.Footer>
        <Button onclick={close}>Close</Button>
      </Dialog.Footer>
    {:else}
      <Dialog.Header>
        <Dialog.Title>Remove duplicate cron jobs?</Dialog.Title>
        <Dialog.Description>
          {extraCount} duplicate cron {extraCount === 1 ? "entry" : "entries"} share
          the same schedule and command as another job, so they can run more than
          once. Removing them keeps the first entry for each job.
        </Dialog.Description>
      </Dialog.Header>

      <ul class="max-h-48 space-y-1 overflow-auto text-sm">
        {#each duplicates as group (group.hash)}
          <li class="rounded-md bg-muted px-2 py-1.5">
            <div class="flex items-center justify-between gap-2">
              <span class="truncate font-medium">
                {cronJobDisplayName(group.job)}
              </span>
              <span class="shrink-0 text-xs text-muted-foreground">
                {group.count} copies
              </span>
            </div>
            <span
              class="block truncate font-mono text-xs text-muted-foreground"
            >
              {group.job.command}
            </span>
          </li>
        {/each}
      </ul>

      <Dialog.Footer>
        <Button disabled={submitting} onclick={close}>Cancel</Button>
        <Button variant="destructive" onclick={remove}>
          <Trash2 />
          Remove duplicates
        </Button>
      </Dialog.Footer>
    {/if}
  </Dialog.Content>
</Dialog.Root>
