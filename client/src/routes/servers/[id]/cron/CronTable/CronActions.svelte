<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import {
    Copy,
    CopyPlus,
    EllipsisVertical,
    Loader,
    Play,
    Square,
    Trash2,
  } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import {
    deleteCronJob,
    duplicateCronJob,
    getPrimaryCronAction,
    performCronAction,
    type CronAction,
    type CronJob,
  } from "../cronJobs.ts";

  let {
    serverId,
    job,
    onUpdated,
    onCreated,
    onDeleted,
  }: {
    serverId: number;
    job: CronJob;
    onUpdated?: (job: CronJob) => void;
    onCreated?: (job: CronJob) => void;
    onDeleted?: (id: number) => void;
  } = $props();

  let pending = $state<CronAction | "duplicate" | "delete" | null>(null);
  let confirmDelete = $state(false);

  const primary = $derived(getPrimaryCronAction(job.status));
  const busy = $derived(pending !== null);

  const actionMessage: Record<CronAction, string> = {
    run: "started",
    resume: "resumed",
    stop: "stopped",
  };

  async function applyAction(action: CronAction) {
    if (busy) return;
    pending = action;

    const { data, error } = await performCronAction(serverId, job.id, action);
    pending = null;

    if (error || !data) {
      toast.error(
        error
          ? unWrapErrorMessage(error, `Failed to ${action} ${job.name}`)
          : `Failed to ${action} ${job.name}`,
      );
      return;
    }

    onUpdated?.(data);
    toast.success(`${job.name} ${actionMessage[action]}`);
  }

  async function duplicate() {
    if (busy) return;
    pending = "duplicate";

    const { data, error } = await duplicateCronJob(serverId, job.id);
    pending = null;

    if (error || !data) {
      toast.error(
        error
          ? unWrapErrorMessage(error, `Failed to duplicate ${job.name}`)
          : `Failed to duplicate ${job.name}`,
      );
      return;
    }

    onCreated?.(data);
    toast.success(`${job.name} duplicated`);
  }

  async function remove() {
    if (busy) return;
    pending = "delete";

    const { error } = await deleteCronJob(serverId, job.id);
    pending = null;

    if (error) {
      toast.error(unWrapErrorMessage(error, `Failed to delete ${job.name}`));
      return;
    }

    confirmDelete = false;
    onDeleted?.(job.id);
    toast.success(`${job.name} deleted`);
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
    class="gap-1.5 w-20"
    disabled={busy}
    title={`${primary.label} ${job.name}`}
    onclick={() => applyAction(primary.kind)}
  >
    {#if pending === primary.kind}
      <Loader class="size-3.5 animate-spin" />
    {:else if primary.kind === "stop"}
      <Square class="size-3.5" />
    {:else}
      <Play class="size-3.5" />
    {/if}
    <span class="hidden md:inline">{primary.label}</span>
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
            onclick={duplicate}
          >
            {#if pending === "duplicate"}
              <Loader class="size-4 animate-spin" />
            {:else}
              <CopyPlus class="size-4" />
            {/if}
            Duplicate
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
        This will permanently delete "{job.name}" and run its schedule no more.
        This action cannot be undone.
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
