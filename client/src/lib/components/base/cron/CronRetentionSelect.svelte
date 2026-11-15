<script lang="ts">
  import * as Field from "#lib/components/ui/field/index.ts";
  import * as Select from "#lib/components/ui/select/index.js";

  let {
    id,
    label,
    description,
    value = $bindable(),
    options,
    unit,
  }: {
    id: string;
    label: string;
    description?: string;
    value: string;
    options: number[];
    unit: string;
  } = $props();

  const optionLabel = (option: number) =>
    `${option} ${unit}${option === 1 ? "" : "s"}`;
</script>

<Field.Field>
  <Field.Label for={id}>{label}</Field.Label>
  <Select.Root type="single" bind:value>
    <Select.Trigger {id} class="w-full">
      <Select.Value placeholder="Select a value" />
    </Select.Trigger>
    <Select.Content>
      {#each options as option (option)}
        <Select.Item value={String(option)} label={optionLabel(option)}>
          {optionLabel(option)}
        </Select.Item>
      {/each}
    </Select.Content>
  </Select.Root>
  {#if description}
    <Field.Description>{description}</Field.Description>
  {/if}
</Field.Field>
