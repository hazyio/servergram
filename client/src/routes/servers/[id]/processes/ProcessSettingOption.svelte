<script lang="ts" module>
  export type ProcessColumns =
    "cmd" | "user" | "cpu" | "memory" | "state" | "uptime";
  export interface ProcessSettingsOptionProps {
    columns: Record<ProcessColumns, boolean>;
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
    onSave: (columns: Record<ProcessColumns, boolean>) => void;
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

  function save() {
    useLocalStorage().save<ProcessSettingsOptionProps>(
      "process-explorer-settings",
      {
        columns,
      },
    );
    onSave(JSON.parse(JSON.stringify(columns)));
  }
  onMount(() => {
    const saved = useLocalStorage().get<ProcessSettingsOptionProps>(
      "process-explorer-settings",
      {
        columns: columns,
      },
    );
    columns = saved.columns;
    untrack(() => onSave(saved.columns));
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
