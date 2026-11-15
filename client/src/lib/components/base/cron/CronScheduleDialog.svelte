<script lang="ts">
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import * as Field from "#lib/components/ui/field/index.ts";
  import { Input } from "#lib/components/ui/input/index.ts";
  import * as Select from "#lib/components/ui/select/index.js";
  import {
    describeSchedule,
    isCronExpression,
    parseCronExpression,
    scheduleModes,
    toCronExpression,
    weekdayNames,
    type CronSchedule,
  } from "./cronSchedule.ts";

  let {
    open = $bindable(false),
    expression = "",
    onApply,
  }: {
    open: boolean;
    expression?: string;
    onApply: (schedule: { expression: string; description: string }) => void;
  } = $props();

  const everyOptions = ["1", "2", "5", "10", "15", "20", "30"];

  let mode = $state("daily");
  let every = $state("5");
  let minute = $state(0);
  let hour = $state(2);
  let weekday = $state("1");
  let day = $state(1);
  let custom = $state("");

  function hydrate(value: string) {
    const schedule = value ? parseCronExpression(value) : null;

    if (!schedule) {
      mode = "daily";
      every = "5";
      minute = 0;
      hour = 2;
      weekday = "1";
      day = 1;
      custom = value;
      return;
    }

    mode = schedule.mode;

    switch (schedule.mode) {
      case "minutes":
        every = String(schedule.every);
        break;
      case "hourly":
        minute = schedule.minute;
        break;
      case "daily":
        hour = schedule.hour;
        minute = schedule.minute;
        break;
      case "weekly":
        weekday = String(schedule.weekday);
        hour = schedule.hour;
        minute = schedule.minute;
        break;
      case "monthly":
        day = schedule.day;
        hour = schedule.hour;
        minute = schedule.minute;
        break;
      case "custom":
        custom = schedule.expression;
        break;
    }
  }

  $effect(() => {
    if (open) hydrate(expression);
  });

  const toNumber = (value: number | string) => {
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : 0;
  };

  const schedule = $derived.by((): CronSchedule => {
    switch (mode) {
      case "minutes":
        return { mode: "minutes", every: toNumber(every) };
      case "hourly":
        return { mode: "hourly", minute: toNumber(minute) };
      case "daily":
        return {
          mode: "daily",
          hour: toNumber(hour),
          minute: toNumber(minute),
        };
      case "weekly":
        return {
          mode: "weekly",
          weekday: toNumber(weekday),
          hour: toNumber(hour),
          minute: toNumber(minute),
        };
      case "monthly":
        return {
          mode: "monthly",
          day: toNumber(day),
          hour: toNumber(hour),
          minute: toNumber(minute),
        };
      case "custom":
        return { mode: "custom", expression: custom };
      default:
        return {
          mode: "daily",
          hour: toNumber(hour),
          minute: toNumber(minute),
        };
    }
  });

  const builtExpression = $derived(toCronExpression(schedule));
  const description = $derived(describeSchedule(schedule));
  const customError = $derived(
    schedule.mode === "custom" && !isCronExpression(schedule.expression)
      ? "Enter a cron expression with 5 or 6 fields, for example 0 2 * * *"
      : "",
  );
  const canApply = $derived(!customError);

  function apply() {
    if (!canApply) return;
    onApply({ expression: builtExpression, description });
    open = false;
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="sm:max-w-md">
    <Dialog.Header>
      <Dialog.Title>Set run schedule</Dialog.Title>
      <Dialog.Description>
        Pick how often this job runs. The cron expression is generated for you.
      </Dialog.Description>
    </Dialog.Header>

    <div class="grid gap-4">
      <Field.Field>
        <Field.Label for="schedule-mode">Frequency</Field.Label>
        <Select.Root type="single" bind:value={mode}>
          <Select.Trigger id="schedule-mode" class="w-full">
            <Select.Value />
          </Select.Trigger>
          <Select.Content>
            {#each scheduleModes as option (option.value)}
              <Select.Item value={option.value} label={option.label}>
                {option.label}
              </Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      </Field.Field>

      {#if mode === "minutes"}
        <Field.Field>
          <Field.Label for="schedule-every">Run every</Field.Label>
          <Select.Root type="single" bind:value={every}>
            <Select.Trigger id="schedule-every" class="w-full">
              <Select.Value />
            </Select.Trigger>
            <Select.Content>
              {#each everyOptions as option (option)}
                <Select.Item
                  value={option}
                  label={`${option} minute${option === "1" ? "" : "s"}`}
                >
                  {option} minute{option === "1" ? "" : "s"}
                </Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
        </Field.Field>
      {:else if mode === "hourly"}
        <Field.Field>
          <Field.Label for="schedule-minute">At minute</Field.Label>
          <Input
            id="schedule-minute"
            type="number"
            min="0"
            max="59"
            bind:value={minute}
          />
          <Field.Description>Runs once every hour.</Field.Description>
        </Field.Field>
      {:else if mode === "daily" || mode === "weekly" || mode === "monthly"}
        <div class="grid gap-4 sm:grid-cols-2">
          {#if mode === "weekly"}
            <Field.Field class="sm:col-span-2">
              <Field.Label for="schedule-weekday">Day of the week</Field.Label>
              <Select.Root type="single" bind:value={weekday}>
                <Select.Trigger id="schedule-weekday" class="w-full">
                  <Select.Value />
                </Select.Trigger>
                <Select.Content>
                  {#each weekdayNames as name, index (name)}
                    <Select.Item value={String(index)} label={name}>
                      {name}
                    </Select.Item>
                  {/each}
                </Select.Content>
              </Select.Root>
            </Field.Field>
          {:else if mode === "monthly"}
            <Field.Field class="sm:col-span-2">
              <Field.Label for="schedule-day">Day of the month</Field.Label>
              <Input
                id="schedule-day"
                type="number"
                min="1"
                max="31"
                bind:value={day}
              />
              <Field.Description>
                Use 1–28 to run on every month, or 29–31 for months that have
                that day.
              </Field.Description>
            </Field.Field>
          {/if}
          <Field.Field>
            <Field.Label for="schedule-hour">Hour</Field.Label>
            <Input
              id="schedule-hour"
              type="number"
              min="0"
              max="23"
              bind:value={hour}
            />
          </Field.Field>
          <Field.Field>
            <Field.Label for="schedule-daily-minute">Minute</Field.Label>
            <Input
              id="schedule-daily-minute"
              type="number"
              min="0"
              max="59"
              bind:value={minute}
            />
          </Field.Field>
        </div>
        <p class="-mt-2 text-xs text-muted-foreground">
          Uses the server's local time. Hour is 0–23, minute is 0–59.
        </p>
      {:else}
        <Field.Field data-invalid={!!customError}>
          <Field.Label for="custom-cron">Cron expression</Field.Label>
          <Input
            id="custom-cron"
            bind:value={custom}
            placeholder="0 2 * * *"
            class="font-mono"
            aria-invalid={customError ? "true" : "false"}
          />
          <Field.Description>
            Five fields: minute, hour, day of month, month, day of week.
          </Field.Description>
          {#if customError}
            <Field.Error>{customError}</Field.Error>
          {/if}
        </Field.Field>
      {/if}

      <div class="rounded-lg border border-border bg-muted/40 p-3">
        <p
          class="text-xs font-medium tracking-wide text-muted-foreground uppercase"
        >
          Cron expression
        </p>
        <code class="mt-1 block font-mono text-sm wrap-break-word"
          >{builtExpression || "—"}</code
        >
        <p class="mt-1 text-xs text-muted-foreground">{description}</p>
      </div>
    </div>

    <Dialog.Footer>
      <Button variant="outline" onclick={() => (open = false)}>Cancel</Button>
      <Button disabled={!canApply} onclick={apply}>Use this schedule</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
