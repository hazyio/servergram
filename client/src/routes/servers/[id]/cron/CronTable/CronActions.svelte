<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import {
    Copy,
    EllipsisVertical,
    Loader,
    Pencil,
    Play,
    Trash2,
  } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import {
    cronJobDisplayName,
    deleteCronJob,
    runCronJob,
    type CronJob,
  } from "../cronJobs.ts";
  import CronEditDialog from "../CronEditDialog.svelte";

  let {
    serverId,
    job,
    onUpdated,
    onDeleted,
  }: {
    serverId: number;
    job: CronJob;
    onUpdated?: (job: CronJob) => void;
    onDeleted?: (id: number) => void;
  } = $props();

  let pending = $state<"run" | "delete" | null>(null);
  let confirmDelete = $state(false);
  let editOpen = $state(false);
  let runErrorOpen = $state(false);
  let runError = $state("");

  const busy = $derived(pending !== null);
  const displayName = $derived(cronJobDisplayName(job));
  const isUnmanaged = $derived(job.status === "unmanaged");

  async function runJob() {
    if (busy) return;
    pending = "run";
    runError = "";

    const { error } = await runCronJob(serverId, job.id);
    pending = null;

    if (error) {
      runError = unWrapErrorMessage(error, `Failed to run ${displayName}`);
      runErrorOpen = true;
      return;
    }

    runErrorOpen = false;
    onUpdated?.({
      ...job,
      last_run: new Date().toISOString(),
    });
    toast.success(`${displayName} started`);
  }

  async function remove() {
    if (busy) return;
    pending = "delete";

    const { error } = await deleteCronJob(serverId, job.id);
    pending = null;

    if (error) {
      toast.error(unWrapErrorMessage(error, `Failed to delete ${displayName}`));
      return;
    }

    confirmDelete = false;
    onDeleted?.(job.id);
    toast.success(`${displayName} deleted`);
  }

  async function copyCommand() {
    try {
      await navigator.clipboard.writeText(job.command);
      toast.success("Command copied to clipboard");
    } catch {
      toast.error("Could not copy the command");
    }
  }
</script>

<div class="flex items-center justify-end gap-1">
  <Button
    variant="outline"
    size="sm"
    class="gap-1.5 "
    disabled={busy}
    title={`Run ${displayName}`}
    onclick={runJob}
  >
    {#if pending === "run"}
      <Loader class="size-3.5 animate-spin" />
    {:else}
      <Play class="size-3.5" />
    {/if}
    <span class="hidden md:inline">Run</span>
  </Button>

  <Popover.Root>
    <Popover.Trigger>
      <Button
        variant="ghost"
        size="icon-sm"
        disabled={busy}
        title="More actions"
      >
        <EllipsisVertical />
      </Button>
    </Popover.Trigger>

    <Popover.Content align="end" class="w-52">
      <Popover.Header>More actions</Popover.Header>
      <div class="flex flex-col gap-1">
        <Popover.Close>
          <Button
            variant="ghost"
            class="w-full cursor-pointer justify-start gap-2 py-5"
            onclick={copyCommand}
          >
            <Copy class="size-4" />
            Copy command
          </Button>
        </Popover.Close>
        <Popover.Close>
          <Button
            variant="ghost"
            class="w-full cursor-pointer justify-start gap-2 py-5"
            disabled={busy}
            onclick={() => {
              editOpen = true;
            }}
          >
            <Pencil class="size-4" />
            Edit
          </Button>
        </Popover.Close>
        <Popover.Close>
          <Button
            variant="ghost"
            class="w-full cursor-pointer justify-start gap-2 py-5 text-destructive hover:text-destructive"
            disabled={busy}
            onclick={() => {
              confirmDelete = true;
            }}
          >
            <Trash2 class="size-4" />
            Delete
          </Button>
        </Popover.Close>
      </div>
    </Popover.Content>
  </Popover.Root>
</div>

<AlertDialog.Root bind:open={confirmDelete}>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Delete cron job?</AlertDialog.Title>
      <AlertDialog.Description>
        This will permanently delete "{displayName}" and run its schedule no
        more. This action cannot be undone.
      </AlertDialog.Description>
    </AlertDialog.Header>
    <AlertDialog.Footer>
      <AlertDialog.Cancel disabled={pending === "delete"}>
        Cancel
      </AlertDialog.Cancel>
      <Button
        variant="destructive"
        disabled={pending === "delete"}
        onclick={remove}
      >
        {#if pending === "delete"}
          <Loader class="size-4 animate-spin" />
          Deleting...
        {:else}
          Delete
        {/if}
      </Button>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>

<AlertDialog.Root bind:open={runErrorOpen}>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Could not run cron job</AlertDialog.Title>
    </AlertDialog.Header>
    <ErrorSection
      title="Failed to run cron job"
      message={runError}
      onRetry={runJob}
    />
    <AlertDialog.Footer>
      <AlertDialog.Cancel disabled={pending === "run"}>Close</AlertDialog.Cancel
      >
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>

<CronEditDialog
  bind:open={editOpen}
  {serverId}
  {job}
  onUpdated={(updated) => onUpdated?.(updated)}
/>
