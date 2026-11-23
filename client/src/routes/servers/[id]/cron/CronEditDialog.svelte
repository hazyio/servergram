<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import CronRetentionSelect from "#lib/components/base/cron/CronRetentionSelect.svelte";
  import CronScheduleDialog from "#lib/components/base/cron/CronScheduleDialog.svelte";
  import {
    describeCronExpression,
    isCronExpression,
  } from "#lib/components/base/cron/cronSchedule.ts";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import * as Field from "#lib/components/ui/field/index.ts";
  import { Input } from "#lib/components/ui/input/index.ts";
  import { Separator } from "#lib/components/ui/separator/index.js";
  import { CalendarClock, CircleAlert, Loader } from "@lucide/svelte";
  import {
    Alert,
    AlertDescription,
    AlertTitle,
  } from "#lib/components/ui/alert/index.js";
  import { Checkbox } from "#lib/components/ui/checkbox/index.ts";
  import { untrack } from "svelte";
  import { toast } from "svelte-sonner";
  import * as z from "zod";
  import {
    cronJobDisplayName,
    getCronJobInfo,
    updateCronJob,
    type CronJob,
    type CronJobInfo,
  } from "./cronJobs.ts";

  let {
    open = $bindable(false),
    serverId,
    job,
    onUpdated,
  }: {
    open: boolean;
    serverId: number;
    job: CronJob;
    onUpdated?: (job: CronJob) => void;
  } = $props();

  const runCountOptions = [0, 5, 10, 20, 40, 80, 160, 240, 480, 960];

  let name = $state("");
  let command = $state("");
  let schedule = $state("");
  let scheduleDescription = $state("");
  let runLogs = $state("0");
  let scheduleDialogOpen = $state(false);
  let submitting = $state(false);
  let loading = $state(false);
  let loadError = $state("");
  let loadedInfo = $state<CronJobInfo | null>(null);
  let requestError = $state("");
  let requestToken = 0;
  let formError = $state({
    name: "",
    command: "",
    schedule: "",
  });

  let wrapWithAgent = $state(false);
  let retentionEnabled = $state(false);

  const isManaged = $derived(
    loadedInfo ? loadedInfo.managed : job.status !== "unmanaged",
  );

  const formSchema = z.object({
    name: z.string().trim().max(100, "Job name is too long"),
    command: z.string().trim().min(1, "Command is required"),
    schedule: z
      .string()
      .trim()
      .min(1, "Choose a run schedule")
      .refine(
        (value) => isCronExpression(value),
        "Use a cron expression with 5 or 6 fields",
      ),
  });

  async function load() {
    const token = ++requestToken;

    loading = true;
    loadError = "";
    requestError = "";

    const { data, error } = await getCronJobInfo(serverId, job.id);

    if (token !== requestToken) return;

    loading = false;

    if (error || !data) {
      loadError = error
        ? unWrapErrorMessage(error, "Failed to load cron job")
        : "Failed to load cron job";
      return;
    }

    loadedInfo = data;
    name = data.name || job.name;
    command = data.command;
    schedule = data.schedule;
    scheduleDescription = data.schedule
      ? describeCronExpression(data.schedule)
      : "";
    runLogs = String(data.run_logs_keep_count);
    wrapWithAgent = data.managed;
    retentionEnabled = data.run_logs_keep_count > 0;
    formError.name = "";
    formError.command = "";
    formError.schedule = "";
  }

  $effect.pre(() => {
    if (!open) {
      requestToken++;
      loading = false;
      loadError = "";
      loadedInfo = null;
      return;
    }

    untrack(() => {
      void load();
    });
  });

  function close() {
    if (submitting) return;
    open = false;
    requestError = "";
    wrapWithAgent = false;
    retentionEnabled = false;
  }

  async function save() {
    if (submitting) return;

    formError.name = "";
    formError.command = "";
    formError.schedule = "";
    requestError = "";

    const result = formSchema.safeParse({ name, command, schedule });
    if (!result.success) {
      for (const issue of result.error.issues) {
        if (issue.path[0] === "name") formError.name = issue.message;
        if (issue.path[0] === "command") formError.command = issue.message;
        if (issue.path[0] === "schedule") formError.schedule = issue.message;
      }
      return;
    }

    submitting = true;
    const { error } = await updateCronJob(serverId, job.id, {
      name: result.data.name,
      command: result.data.command,
      schedule: result.data.schedule,
      wrap: wrapWithAgent,
      run_logs_keep_count: retentionEnabled ? Number(runLogs) : 0,
    });
    submitting = false;

    if (error) {
      requestError = unWrapErrorMessage(error, "Failed to update cron job");
      return;
    }

    const updated: CronJob = {
      ...job,
      name: result.data.name,
      command: result.data.command,
      schedule: result.data.schedule,
      run_logs_keep_count: retentionEnabled ? Number(runLogs) : 0,
    };

    onUpdated?.(updated);
    toast.success(`"${cronJobDisplayName(updated)}" updated`);
    open = false;
  }

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    void save();
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-h-[85dvh] overflow-y-auto sm:max-w-lg">
    <Dialog.Header>
      <Dialog.Title>Edit cron job</Dialog.Title>
      <Dialog.Description>
        Update the name, command, or run schedule. If the job is wrapped with
        the agent you can also change how much history it keeps.
      </Dialog.Description>
    </Dialog.Header>

    {#if loading}
      <Loading loading class="h-50" text="Loading cron job..." />
    {:else if loadError}
      <ErrorSection
        title="Failed to load cron job"
        message={loadError}
        onRetry={load}
      />
      <Dialog.Footer>
        <Button variant="outline" onclick={() => (open = false)}>Close</Button>
      </Dialog.Footer>
    {:else}
      <form
        id="edit-cron-job-form"
        onsubmit={handleSubmit}
        class="flex flex-col gap-4"
      >
        <Field.Field data-invalid={!!formError.name}>
          <Field.Label for="edit-cron-name">Name</Field.Label>
          <Input
            id="edit-cron-name"
            type="text"
            bind:value={name}
            oninput={() => (formError.name = "")}
            placeholder="e.g. Database backup"
            disabled={submitting}
            aria-invalid={formError.name ? "true" : "false"}
          />
          {#if formError.name}
            <Field.Error>{formError.name}</Field.Error>
          {/if}
        </Field.Field>

        <Field.Field data-invalid={!!formError.command}>
          <Field.Label for="edit-cron-command">Command</Field.Label>
          <Input
            id="edit-cron-command"
            type="text"
            bind:value={command}
            oninput={() => (formError.command = "")}
            placeholder="/usr/local/bin/backup-db.sh"
            class="font-mono"
            disabled={submitting}
            aria-invalid={formError.command ? "true" : "false"}
          />
          {#if formError.command}
            <Field.Error>{formError.command}</Field.Error>
          {/if}
        </Field.Field>

        <Field.Field data-invalid={!!formError.schedule}>
          <Field.Label for="edit-cron-schedule">Run schedule</Field.Label>
          <Button
            id="edit-cron-schedule"
            variant="outline"
            class="h-auto w-full justify-between gap-2 px-2.5 py-2 text-left font-normal"
            onclick={() => {
              formError.schedule = "";
              scheduleDialogOpen = true;
            }}
            disabled={submitting}
            aria-invalid={formError.schedule ? "true" : "false"}
          >
            {#if schedule}
              <span class="min-w-0 flex-1">
                <span class="block truncate font-mono text-xs">{schedule}</span>
                <span class="block truncate text-xs text-muted-foreground">
                  {scheduleDescription}
                </span>
              </span>
            {:else}
              <span class="flex-1 text-muted-foreground">Set a schedule</span>
            {/if}
            <CalendarClock class="shrink-0 text-muted-foreground" />
          </Button>
          {#if formError.schedule}
            <Field.Error>{formError.schedule}</Field.Error>
          {/if}
        </Field.Field>

        {#if !loading && !loadError}
          <Separator />

          <div class="grid gap-3 md:grid-cols-2">
            <Field.Label>
              <Field.Field orientation="horizontal">
                <Checkbox
                  id="edit-wrap-with-agent"
                  name="edit-wrap-with-agent"
                  bind:checked={wrapWithAgent}
                  disabled={submitting}
                />
                <Field.Content>
                  <Field.Title>Wrap with agent</Field.Title>
                  <Field.Description>
                    Let ServerGram manage this job so you get logging, status
                    tracking, and run history.
                  </Field.Description>
                </Field.Content>
              </Field.Field>
            </Field.Label>

            <Field.Label class={wrapWithAgent ? "" : "opacity-60"}>
              <Field.Field orientation="horizontal">
                <Checkbox
                  id="edit-cron-retention"
                  name="edit-cron-retention"
                  bind:checked={retentionEnabled}
                  disabled={submitting || !wrapWithAgent}
                />
                <Field.Content>
                  <Field.Title>Retention</Field.Title>
                  <Field.Description>
                    Keep run metadata and output logs for this job.
                  </Field.Description>
                </Field.Content>
              </Field.Field>
            </Field.Label>
          </div>

          {#if wrapWithAgent && retentionEnabled}
            <Alert class="border-amber-500/40 bg-amber-500/5">
              <CircleAlert class="text-amber-600 dark:text-amber-400" />
              <AlertTitle>Logging uses storage</AlertTitle>
              <AlertDescription>
                Run logs can use a lot of storage. If the server has 200 MB or
                less free storage, older run logs will be deleted automatically.
              </AlertDescription>
            </Alert>

            <CronRetentionSelect
              id="edit-run-logs"
              label="Runs and logs to keep"
              description="Stores run metadata, stdout, and stderr for recent runs."
              bind:value={runLogs}
              options={runCountOptions}
              unit="run"
            />
          {:else}
            <p class="text-sm text-muted-foreground">
              Retention is disabled. No run history or logs will be kept for
              this job.
            </p>
          {/if}
        {/if}

        {#if requestError}
          <ErrorSection
            title="Failed to update cron job"
            message={requestError}
            onRetry={save}
          />
        {/if}

        <Dialog.Footer>
          <Button
            type="button"
            variant="outline"
            disabled={submitting}
            onclick={close}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={submitting}>
            {#if submitting}
              <Loader class="animate-spin" />
              Saving...
            {:else}
              Save changes
            {/if}
          </Button>
        </Dialog.Footer>
      </form>
    {/if}
  </Dialog.Content>
</Dialog.Root>

<CronScheduleDialog
  bind:open={scheduleDialogOpen}
  expression={schedule}
  onApply={({ expression, description }) => {
    schedule = expression;
    scheduleDescription = description;
    formError.schedule = "";
  }}
/>
