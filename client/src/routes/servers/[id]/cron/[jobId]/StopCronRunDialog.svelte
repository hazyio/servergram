<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import { toast } from "svelte-sonner";
  import { stopCronJobRun, type CronRun } from "../cronJobs.ts";

  let {
    open = $bindable(false),
    serverId,
    run,
    onStopped,
  }: {
    open: boolean;
    serverId: number;
    run: CronRun | null;
    onStopped?: (runId: number) => void;
  } = $props();

  let stopping = $state(false);
  let stopError = $state("");

  const runLabel = $derived(run ? `#${run.id}` : "this run");

  async function stopRun() {
    if (!run || stopping) return;

    stopping = true;
    stopError = "";

    const { error } = await stopCronJobRun(serverId, run.id);

    stopping = false;

    if (error) {
      stopError = unWrapErrorMessage(error, `Failed to stop run ${runLabel}`);
      return;
    }

    open = false;
    onStopped?.(run.id);
    toast.success(`Run ${runLabel} stopped`);
  }

  $effect.pre(() => {
    if (!open) {
      stopError = "";
    }
  });
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Content>
    {#if stopping}
      <AlertDialog.Header>
        <AlertDialog.Title>Stopping run</AlertDialog.Title>
      </AlertDialog.Header>
      <Loading
        loading
        containerClass="h-40"
        class="h-8 w-8"
        text={`Stopping run ${runLabel}...`}
      />
    {:else if stopError}
      <AlertDialog.Header>
        <AlertDialog.Title>Could not stop run</AlertDialog.Title>
      </AlertDialog.Header>
      <ErrorSection
        title="Failed to stop cron run"
        message={stopError}
        onRetry={stopRun}
      />
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Close</AlertDialog.Cancel>
      </AlertDialog.Footer>
    {:else}
      <AlertDialog.Header>
        <AlertDialog.Title>Stop run {runLabel}?</AlertDialog.Title>
        <AlertDialog.Description>
          This cancels the running process for this cron job. The run is marked
          as canceled and cannot be resumed.
        </AlertDialog.Description>
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <Button variant="destructive" onclick={stopRun}>Stop run</Button>
      </AlertDialog.Footer>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
