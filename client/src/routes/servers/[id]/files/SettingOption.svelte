<script lang="ts" module>
  export type SettingColumns = "size" | "modified" | "mode" | "owner";
  export interface SettingOptionProps {
    showHidden: boolean;
    dirsFirst: boolean;
    columns: Record<SettingColumns, boolean>;
    askBeforeDelete: boolean;
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


  let {
    showHidden = $bindable(false),
    dirsFirst = $bindable(true),
    columns = $bindable(),
    askBeforeDelete = $bindable(true),
  }: SettingOptionProps = $props();
  let mounted = $state(false);
  $effect(() => {
    if (!mounted) return;
    useLocalStorage().save<SettingOptionProps>("file-explorer-settings", {
      showHidden,
      dirsFirst,
      columns,
      askBeforeDelete,
    });
  });
  onMount(() => {
    const saved = useLocalStorage().get<SettingOptionProps>(
      "file-explorer-settings",
      {
        showHidden: false,
        dirsFirst: true,
        columns: {
          size: true,
          modified: true,
          mode: true,
          owner: true,
        },
        askBeforeDelete: true,
      },
    );
    showHidden = saved.showHidden;
    dirsFirst = saved.dirsFirst;
    columns = saved.columns;
    askBeforeDelete = saved.askBeforeDelete;
    mounted = true;
  });
</script>

<Popover.Root>
  <Popover.Trigger>
    <Button variant="outline"
      ><Settings /></Button
    >
  </Popover.Trigger>

  <Popover.Content align="end" class="min-h-50">
    <Popover.Header>
      <h3 class="text-lg font-medium">Settings</h3>
    </Popover.Header>
    <div class="mt-3">
      <Field.Group class="max-w-sm">
        <Field.Label>
          <Field.Field orientation="horizontal">
            <Checkbox
              id="show-hidden"
              name="show-hidden"
              bind:checked={showHidden}
            />
            <Field.Content>
              <Field.Title>Show Hidden</Field.Title>
              <Field.Description>
                Toggle hidden files and directories.
              </Field.Description>
            </Field.Content>
          </Field.Field></Field.Label
        >
        <Field.Label>
          <Field.Field orientation="horizontal">
            <Checkbox
              id="dirs-first"
              name="dirs-first"
              bind:checked={dirsFirst}
            />
            <Field.Content>
              <Field.Title>Show Directories First</Field.Title>
            </Field.Content>
          </Field.Field></Field.Label
        >
        <Field.Label>
          <Field.Field orientation="horizontal">
            <Checkbox
              id="ask-before-delete"
              name="ask-before-delete"
              bind:checked={askBeforeDelete}
            />
            <Field.Content>
              <Field.Title>Ask Before Delete</Field.Title>
              <Field.Description>
                Show prompt when deleting one file or directory.</Field.Description
              >
            </Field.Content>
          </Field.Field></Field.Label
        >

        <h4 class="text-sm font-medium">Show as column</h4>
        {#each ["size", "modified", "mode", "owner"] as SettingColumns[] as option (option)}
          <Field.Field orientation="horizontal">
            <Checkbox
              id={`column-${option}`}
              name={`column-${option}`}
              bind:checked={columns[option]}
            />
            <Field.Label for={`column-${option}`} class="capitalize">{option}</Field.Label>
          </Field.Field>
        {/each}
      </Field.Group>
    </div>
  </Popover.Content>
</Popover.Root>
