<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import Loading from "#lib/components/base/loading.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import * as Field from "#lib/components/ui/field/index.ts";
  import Input from "#lib/components/ui/input/input.svelte";
  import { toast } from "svelte-sonner";
  import { renameCronJob, type CronJob } from "../cronJobs.ts";

  let {
    open = $bindable(false),
    serverId,
    job,
    onRenamed,
  }: {
    open: boolean;
    serverId: number;
    job: CronJob;
    onRenamed?: (job: CronJob) => void;
  } = $props();

  let name = $state("");
  let submitting = $state(false);
  let formError = $state("");

  const isNaming = $derived(!job.name.trim());
  const title = $derived(isNaming ? "Name cron job" : "Rename cron job");

  $effect(() => {
    if (open) {
      name = job.name;
      formError = "";
    }
  });

  async function submit() {
    const value = name.trim();
    if (!value) {
      formError = "Name is required";
      return;
    }

    submitting = true;
    formError = "";

    const { data, error } = await renameCronJob(serverId, job.id, value);
    submitting = false;

    if (error || !data) {
      formError = error
        ? unWrapErrorMessage(
            error,
            `Failed to ${isNaming ? "name" : "rename"} job`,
          )
        : `Failed to ${isNaming ? "name" : "rename"} job`;
      return;
    }

    onRenamed?.(data);
    toast.success(
      isNaming ? `Job named "${data.name}"` : `Renamed to "${data.name}"`,
    );
    open = false;
  }

  function close() {
    if (submitting) return;
    open = false;
    formError = "";
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content showCloseButton={false}>
    <form
      onsubmit={(event) => {
        event.preventDefault();
        submit();
      }}
      class="flex flex-col gap-4"
    >
      <Dialog.Header>
        <Dialog.Title>{title}</Dialog.Title>
        <Dialog.Description>
          Choose a name for <span class="font-mono">{job.command}</span>.
        </Dialog.Description>
      </Dialog.Header>

      <Field.Field data-invalid={!!formError}>
        <Field.Label for="cron-job-name">Name</Field.Label>
        <Input
          id="cron-job-name"
          placeholder="e.g. Database backup"
          disabled={submitting}
          type="text"
          bind:value={name}
          aria-invalid={formError ? "true" : "false"}
          oninput={() => {
            formError = "";
          }}
        />
        {#if formError}
          <Field.Error>{formError}</Field.Error>
        {/if}
      </Field.Field>

      <Dialog.Footer>
        <Button disabled={submitting} onclick={close}>Cancel</Button>
        <Button type="submit" disabled={submitting}>
          {#if submitting}
            <Loading loading />
            Saving...
          {:else}
            {isNaming ? "Save name" : "Rename"}
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
