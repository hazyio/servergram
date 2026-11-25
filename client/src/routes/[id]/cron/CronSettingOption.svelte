<script lang="ts" module>
  export type CronColumns =
    "status" | "job" | "schedule" | "lastRun" | "nextRun";
  export interface CronSettingOptionProps {
    columns: Record<CronColumns, boolean>;
  }
</script>

<script lang="ts">
  import Button from "#lib/components/ui/button/button.svelte";
  import { Checkbox } from "#lib/components/ui/checkbox/index.js";
  import * as Field from "#lib/components/ui/field/index.js";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import { useLocalStorage } from "#lib/local-storage.ts";
  import { Settings } from "@lucide/svelte";
  import { onMount } from "svelte";

  let { columns = $bindable() }: CronSettingOptionProps = $props();
  let mounted = $state(false);

  const options: { value: CronColumns; label: string }[] = [
    { value: "status", label: "Status" },
    { value: "job", label: "Job" },
    { value: "schedule", label: "Schedule" },
    { value: "lastRun", label: "Last run" },
    { value: "nextRun", label: "Next run" },
  ];

  $effect(() => {
    if (!mounted) return;
    useLocalStorage().save<CronSettingOptionProps>("cron-explorer-settings", {
      columns,
    });
  });

  onMount(() => {
    const saved = useLocalStorage().get<CronSettingOptionProps>(
      "cron-explorer-settings",
      { columns },
    );
    columns = { ...columns, ...saved.columns };
    mounted = true;
  });
</script>

<Popover.Root>
  <Popover.Trigger>
    <Button variant="outline" title="Choose visible columns"
      ><Settings /></Button
    >
  </Popover.Trigger>

  <Popover.Content align="end" class="min-h-50">
    <Popover.Header>
      <h3 class="text-lg font-medium">Settings</h3>
    </Popover.Header>
    <div class="mt-3">
      <Field.Group class="max-w-sm">
        <h4 class="text-sm font-medium">Show as column</h4>
        {#each options as option (option.value)}
          <Field.Field orientation="horizontal">
            <Checkbox
              id={`cron-column-${option.value}`}
              name={`cron-column-${option.value}`}
              bind:checked={columns[option.value]}
            />
            <Field.Label for={`cron-column-${option.value}`}
              >{option.label}</Field.Label
            >
          </Field.Field>
        {/each}
      </Field.Group>
    </div>
  </Popover.Content>
</Popover.Root>
