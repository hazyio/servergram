<script lang="ts">
  import { goto } from "$app/navigation";
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import appPaths from "#lib/appPaths.ts";
  import CronRetentionSelect from "#lib/components/base/cron/CronRetentionSelect.svelte";
  import CronScheduleDialog from "#lib/components/base/cron/CronScheduleDialog.svelte";
  import { isCronExpression } from "#lib/components/base/cron/cronSchedule.ts";
  import {
    Alert,
    AlertDescription,
    AlertTitle,
  } from "#lib/components/ui/alert/index.js";
  import { Button } from "#lib/components/ui/button/index.ts";
  import Card from "#lib/components/ui/card/card.svelte";
  import { Checkbox } from "#lib/components/ui/checkbox/index.js";
  import * as Field from "#lib/components/ui/field/index.ts";
  import { Input } from "#lib/components/ui/input/index.ts";
  import { Separator } from "#lib/components/ui/separator/index.js";
  import { CircleAlert, CalendarClock } from "@lucide/svelte";
  import Loader from "@lucide/svelte/icons/loader";
  import { toast } from "svelte-sonner";
  import * as z from "zod";
  import type { PageProps } from "./$types";
  import { createCronJob } from "../cronJobs.ts";

  let { params }: PageProps = $props();
  const serverId = $derived(Number(params.id));

  const runCountOptions = [5, 10, 20, 40, 80, 160, 240, 480, 960];
  const errorLogDayOptions = [1, 7, 14, 30, 60, 90, 180, 365];

  let name = $state("");
  let command = $state("");
  let schedule = $state("");
  let scheduleDescription = $state("");
  let scheduleDialogOpen = $state(false);
  let wrapWithAgent = $state(false);
  let metadataRuns = $state("20");
  let runLogs = $state("20");
  let errorLogDays = $state("14");
  let submitting = $state(false);
  let formError = $state({
    name: "",
    command: "",
    schedule: "",
  });

  const formSchema = z.object({
    name: z
      .string()
      .trim()
      .min(1, "Job name is required")
      .max(100, "Job name is too long"),
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

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    formError.name = "";
    formError.command = "";
    formError.schedule = "";

    const result = formSchema.safeParse({ name, command, schedule });
    if (!result.success) {
      for (const issue of result.error.issues) {
        if (issue.path[0] === "name") formError.name = issue.message;
        if (issue.path[0] === "command") formError.command = issue.message;
        if (issue.path[0] === "schedule") {
          formError.schedule = issue.message;
        }
      }
      return;
    }

    submitting = true;
    const { data, error } = await createCronJob(serverId, {
      name: result.data.name,
      command: result.data.command,
      schedule: result.data.schedule,
      wrap: wrapWithAgent,
      meta_keep_count: Number(metadataRuns),
      run_logs_keep_count: Number(runLogs),
      error_logs_keep_days: Number(errorLogDays),
    });
    submitting = false;

    if (error || !data) {
      toast.error(
        error
          ? unWrapErrorMessage(error, "Failed to add cron job")
          : "Failed to add cron job",
      );
      return;
    }

    toast.success(`Cron job "${data.name}" added`);
    await goto(appPaths.servers.cron(serverId));
  }
</script>

<svelte:head>
  <title>Add Cron Job | ServerGram</title>
  <meta
    name="description"
    content="Create a cron job, configure its schedule, and optionally wrap it with the ServerGram agent for logging and run history."
  />
</svelte:head>

<div class="mx-auto w-full max-w-5xl p-3 sm:p-5">
  <Card class="mt-2 w-full p-4 sm:p-6">
    <div class="mb-5">
      <h1 class="text-lg font-semibold">Add Cron Job</h1>
      <p class="mt-1 text-sm text-muted-foreground">
        Choose when the job runs and whether ServerGram should manage its
        logging and run history.
      </p>
    </div>

    <form
      id="create-cron-job-form"
      onsubmit={handleSubmit}
      class="flex flex-col gap-5"
    >
      <div class="grid gap-x-4 gap-y-5 sm:grid-cols-2">
        <Field.Field data-invalid={!!formError.name}>
          <Field.Label for="cron-job-name">Name</Field.Label>
          <Input
            id="cron-job-name"
            type="text"
            bind:value={name}
            oninput={() => (formError.name = "")}
            placeholder="e.g. Database backup"
            disabled={submitting}
            required
            aria-invalid={formError.name ? "true" : "false"}
          />
          {#if formError.name}
            <Field.Error>{formError.name}</Field.Error>
          {/if}
        </Field.Field>

        <Field.Field data-invalid={!!formError.schedule}>
          <Field.Label for="cron-schedule">Run schedule</Field.Label>
          <Button
            id="cron-schedule"
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
          {:else}
            <Field.Description>
              The cron expression is generated from the schedule you pick.
            </Field.Description>
          {/if}
        </Field.Field>
      </div>

      <Field.Field data-invalid={!!formError.command}>
        <Field.Label for="cron-command">Command</Field.Label>
        <Input
          id="cron-command"
          type="text"
          bind:value={command}
          oninput={() => (formError.command = "")}
          placeholder="/usr/local/bin/backup-db.sh"
          class="font-mono"
          disabled={submitting}
          required
          aria-invalid={formError.command ? "true" : "false"}
        />
        {#if formError.command}
          <Field.Error>{formError.command}</Field.Error>
        {:else}
          <Field.Description>
            The command or script the job runs on the server, for example
            <code>/usr/bin/rsync -a /var/www /backup/www</code>.
          </Field.Description>
        {/if}
      </Field.Field>

      <Separator />

      <Field.Label>
        <Field.Field orientation="horizontal">
          <Checkbox
            id="wrap-with-agent"
            name="wrap-with-agent"
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

      {#if wrapWithAgent}
        <Alert class="border-amber-500/40 bg-amber-500/5">
          <CircleAlert class="text-amber-600 dark:text-amber-400" />
          <AlertTitle>Logging uses storage</AlertTitle>
          <AlertDescription>
            Run logs can use a lot of storage. If the server has 200 MB or less
            free storage, older run logs will be deleted automatically.
          </AlertDescription>
        </Alert>

        <div class="grid gap-x-4 gap-y-5 sm:grid-cols-2 lg:grid-cols-3">
          <CronRetentionSelect
            id="metadata-runs"
            label="Runs metadata to keep"
            description="Stores status, start time, end time, and exit code for recent runs."
            bind:value={metadataRuns}
            options={runCountOptions}
            unit="run"
          />
          <CronRetentionSelect
            id="run-logs"
            label="Run logs to keep"
            description="Stores stdout and stderr for recent runs."
            bind:value={runLogs}
            options={runCountOptions}
            unit="run"
          />
          <CronRetentionSelect
            id="error-log-days"
            label="Error logs to keep"
            description="Stores error logs for the selected number of days."
            bind:value={errorLogDays}
            options={errorLogDayOptions}
            unit="day"
          />
        </div>
      {/if}

      <div class="flex justify-end gap-2 pt-2">
        <Button
          variant="outline"
          href={appPaths.servers.cron(serverId)}
          disabled={submitting}
        >
          Cancel
        </Button>
        <Button type="submit" form="create-cron-job-form" disabled={submitting}>
          {#if submitting}
            <Loader class="animate-spin" />
            Saving...
          {:else}
            Add job
          {/if}
        </Button>
      </div>
    </form>
  </Card>
</div>

<CronScheduleDialog
  bind:open={scheduleDialogOpen}
  expression={schedule}
  onApply={({ expression, description }) => {
    schedule = expression;
    scheduleDescription = description;
    formError.schedule = "";
  }}
/>
