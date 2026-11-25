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
  import { Bot, CircleCheck } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import {
    cronJobDisplayName,
    wrapCronJobsWithAgent,
    type CronJob,
  } from "./cronJobs.ts";

  let {
    open = $bindable(false),
    serverId,
    jobs,
    onWrapped,
  }: {
    open: boolean;
    serverId: number;
    jobs: CronJob[];
    onWrapped?: (fixed: number) => void;
  } = $props();

  let submitting = $state(false);
  let error = $state("");
  let wrappedCount = $state(0);

  const count = $derived(jobs.length);

  async function wrap() {
    if (submitting) return;

    submitting = true;
    error = "";

    const { data, error: requestError } = await wrapCronJobsWithAgent(serverId);
    submitting = false;

    if (requestError || !data) {
      error = requestError
        ? unWrapErrorMessage(requestError, "Failed to wrap cron jobs")
        : "Failed to wrap cron jobs";
      return;
    }

    if (data.fixed === 0) {
      error = "No unmanaged cron jobs were found";
      return;
    }

    wrappedCount = data.fixed;
    onWrapped?.(data.fixed);
    toast.success(
      `${wrappedCount} cron ${wrappedCount === 1 ? "job" : "jobs"} wrapped with the agent`,
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
      wrappedCount = 0;
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content showCloseButton={false}>
    {#if submitting}
      <Dialog.Header>
        <Dialog.Title>Wrapping cron jobs with the agent</Dialog.Title>
        <Dialog.Description>
          Enabling logging, status tracking, and run history for {count} cron
          {count === 1 ? "job" : "jobs"}.
        </Dialog.Description>
      </Dialog.Header>
      <Loading
        loading
        class="h-50"
        text={`Wrapping ${count} cron ${count === 1 ? "job" : "jobs"}...`}
      />
    {:else if error}
      <Dialog.Header>
        <Dialog.Title>Could not wrap cron jobs</Dialog.Title>
      </Dialog.Header>
      <ErrorSection
        title="Failed to wrap cron jobs"
        message={error}
        onRetry={wrap}
      />
      <Dialog.Footer>
        <Button onclick={close}>Close</Button>
      </Dialog.Footer>
    {:else if wrappedCount > 0}
      <Dialog.Header>
        <Dialog.Title class="sr-only">Cron jobs wrapped</Dialog.Title>
      </Dialog.Header>
      <Alert class="border-emerald-500/40 bg-emerald-500/5">
        <CircleCheck class="text-emerald-600 dark:text-emerald-400" />
        <AlertTitle>Cron jobs wrapped</AlertTitle>
        <AlertDescription>
          {wrappedCount} cron {wrappedCount === 1 ? "job is" : "jobs are"} now managed
          by the agent. Logging, status tracking, and run history are available for
          {wrappedCount === 1 ? "it" : "them"}.
        </AlertDescription>
      </Alert>
      <Dialog.Footer>
        <Button onclick={close}>Close</Button>
      </Dialog.Footer>
    {:else}
      <Dialog.Header>
        <Dialog.Title>Wrap cron jobs with the agent?</Dialog.Title>
        <Dialog.Description>
          {count} unmanaged cron {count === 1 ? "job is" : "jobs are"} not tracked
          by the ServerGram agent. Wrapping {count === 1 ? "it" : "them"}{" "}
          gives you logging, status tracking, and run history.
        </Dialog.Description>
      </Dialog.Header>

      <ul class="max-h-48 space-y-1 overflow-auto text-sm">
        {#each jobs as job (job.id)}
          <li class="rounded-md bg-muted px-2 py-1.5">
            <span class="font-medium">{cronJobDisplayName(job)}</span>
            <span class="block truncate font-mono text-xs text-muted-foreground"
              >{job.command}</span
            >
          </li>
        {/each}
      </ul>

      <Dialog.Footer>
        <Button disabled={submitting} onclick={close}>Cancel</Button>
        <Button onclick={wrap}>
          <Bot />
          Wrap with agent
        </Button>
      </Dialog.Footer>
    {/if}
  </Dialog.Content>
</Dialog.Root>
