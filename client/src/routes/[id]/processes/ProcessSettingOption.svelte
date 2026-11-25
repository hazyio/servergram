<script lang="ts" module>
  export type ProcessColumns =
    "cmd" | "user" | "cpu" | "memory" | "state" | "uptime";
  export type ProcessGroupBy = "parent" | "name" | "command";
  export type ProcessGrouping = "none" | ProcessGroupBy;
  export interface ProcessSettingsOptionProps {
    columns: Record<ProcessColumns, boolean>;
    groupBy: ProcessGrouping;
  }
</script>

<script lang="ts">
  import Button from "#lib/components/ui/button/button.svelte";
  import { Checkbox } from "#lib/components/ui/checkbox/index.js";
  import * as Field from "#lib/components/ui/field/index.js";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import { useLocalStorage } from "#lib/local-storage.ts";
  import { Settings } from "@lucide/svelte";
  import { onMount, untrack } from "svelte";
  import { Separator } from "#lib/components/ui/separator/index.js";
  interface ProcessSettingsProps {
    onSave: (settings: ProcessSettingsOptionProps) => void;
  }
  let { onSave }: ProcessSettingsProps = $props();
  let columns = $state<Record<ProcessColumns, boolean>>({
    cmd: true,
    cpu: true,
    memory: true,
    state: true,
    uptime: true,
    user: true,
  });
  let groupBy = $state<ProcessGrouping>("parent");

  const groupOptions: { value: ProcessGrouping; label: string }[] = [
    { value: "none", label: "None" },
    { value: "parent", label: "Parent id" },
    { value: "name", label: "Name" },
    { value: "command", label: "Command" },
  ];

  $effect(() => {
    // Grouping by command needs the command column to be fetched.
    if (!columns.cmd && groupBy === "command") {
      groupBy = "parent";
    }
  });

  function save() {
    if (!columns.cmd && groupBy === "command") {
      groupBy = "parent";
    }
    const settings: ProcessSettingsOptionProps = {
      columns,
      groupBy,
    };
    useLocalStorage().save<ProcessSettingsOptionProps>(
      "process-explorer-settings",
      settings,
    );
    onSave(JSON.parse(JSON.stringify(settings)));
  }
  onMount(() => {
    const saved = useLocalStorage().get<
      ProcessSettingsOptionProps & { groupByParent?: boolean }
    >("process-explorer-settings", {
      columns,
      groupBy,
    });
    columns = { ...columns, ...saved.columns };
    groupBy =
      saved.groupBy ?? (saved.groupByParent === false ? "none" : "parent");
    if (!columns.cmd && groupBy === "command") {
      groupBy = "parent";
    }
    untrack(() =>
      onSave({
        columns: JSON.parse(JSON.stringify(columns)),
        groupBy,
      }),
    );
  });
</script>

<Popover.Root>
  <Popover.Trigger>
    <Button variant="outline"><Settings /></Button>
  </Popover.Trigger>

  <Popover.Content align="end" class="min-h-50">
    <Popover.Header>
      <h3 class="text-lg font-medium">Settings</h3>
    </Popover.Header>
    <div class="mt-3">
      <Field.Group class="max-w-sm">
        <h4 class="text-sm font-medium">Show as column</h4>
        {#each ["cmd", "user", "cpu", "memory", "state", "uptime"] as option (option)}
          <Field.Field orientation="horizontal">
            <Checkbox
              id={`process-column-${option}`}
              name={`process-column-${option}`}
              bind:checked={columns[option as ProcessColumns]}
            />
            <Field.Label for={`process-column-${option}`} class="capitalize"
              >{option}</Field.Label
            >
          </Field.Field>
        {/each}
      </Field.Group>
    </div>
    <Separator
      orientation="horizontal"
      class="mt-4 me-2 data-vertical:h-4 data-vertical:self-auto"
    />
    <div class="mt-4">
      <Field.Group class="max-w-sm">
        <h4 class="text-sm font-medium">Group by</h4>
        {#each groupOptions as option (option.value)}
          <Field.Field orientation="horizontal">
            <input
              type="radio"
              id={`process-group-by-${option.value}`}
              name="process-group-by"
              value={option.value}
              bind:group={groupBy}
              disabled={option.value === "command" && !columns.cmd}
              class="size-4 shrink-0 accent-primary disabled:cursor-not-allowed disabled:opacity-50"
            />
            <Field.Label
              for={`process-group-by-${option.value}`}
              class={option.value === "command" && !columns.cmd
                ? "opacity-50"
                : ""}
            >
              {option.label}
            </Field.Label>
          </Field.Field>
        {/each}
      </Field.Group>
    </div>
    <Separator
      orientation="horizontal"
      class="mt-4 me-2 data-vertical:h-4 data-vertical:self-auto"
    />
    <div class="flex justify-end">
      <Popover.Close
        ><Button
          onclick={() => {
            save();
          }}>Save</Button
        ></Popover.Close
      >
    </div>
  </Popover.Content>
</Popover.Root>
