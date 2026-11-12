<script lang="ts" module>
  export type ProcessColumns =
    "cmd" | "user" | "cpu" | "memory" | "state" | "uptime";
  export interface ProcessSettingsOptionProps {
    columns: Record<ProcessColumns, boolean>;
    groupByParent: boolean;
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
  let groupByParent = $state(true);

  function save() {
    const settings: ProcessSettingsOptionProps = {
      columns,
      groupByParent,
    };
    useLocalStorage().save<ProcessSettingsOptionProps>(
      "process-explorer-settings",
      settings,
    );
    onSave(JSON.parse(JSON.stringify(settings)));
  }
  onMount(() => {
    const saved = useLocalStorage().get<ProcessSettingsOptionProps>(
      "process-explorer-settings",
      {
        columns,
        groupByParent,
      },
    );
    columns = { ...columns, ...saved.columns };
    groupByParent = saved.groupByParent ?? true;
    untrack(() =>
      onSave({
        columns: JSON.parse(JSON.stringify(columns)),
        groupByParent,
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
        <h4 class="text-sm font-medium">Grouping</h4>
        <Field.Field orientation="horizontal">
          <Checkbox
            id="process-group-by-parent"
            name="process-group-by-parent"
            bind:checked={groupByParent}
          />
          <Field.Label for="process-group-by-parent"
            >Group by parent id</Field.Label
          >
        </Field.Field>
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
