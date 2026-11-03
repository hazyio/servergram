<script lang="ts">
  import Button from "#lib/components/ui/button/button.svelte";
  import { Settings } from "@lucide/svelte";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import type { PageProps } from "./$types";
  import { Checkbox } from "#lib/components/ui/checkbox/index.js";
  import * as Field from "#lib/components/ui/field/index.js";
  import { Label } from "#lib/components/ui/label/index.js";
  const columnOptions = [
    { id: "size", label: "Size" },
    { id: "modified", label: "Modified" },
    { id: "mode", label: "Mode" },
    { id: "owner", label: "Owner" },
  ];
  let {
    showHidden = $bindable(false),
    dirsFirst = $bindable(true),
    columns = $bindable({} as Record<string, boolean>),
  } = $props();
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
        <Field.Field orientation="horizontal">
          <Checkbox
            id="show-hidden"
            name="show-hidden"
            bind:checked={showHidden}
          />
          <Label for="show-hidden">Show Hidden</Label>
        </Field.Field>
        <Field.Field orientation="horizontal">
          <Checkbox
            id="dirs-first"
            name="dirs-first"
            bind:checked={dirsFirst}
          />
          <Label for="dirs-first">Show Directories First</Label>
        </Field.Field>
        <h4 class="text-sm font-medium">Show as column</h4>
        {#each columnOptions as option (option.id)}
          <Field.Field orientation="horizontal">
            <Checkbox
              id={`column-${option.id}`}
              name={`column-${option.id}`}
              bind:checked={columns[option.id]}
            />
            <Label for={`column-${option.id}`}>{option.label}</Label>
          </Field.Field>
        {/each}
      </Field.Group>
    </div>
  </Popover.Content>
</Popover.Root>
