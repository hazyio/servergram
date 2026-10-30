<script lang="ts">
  import * as DropdownMenu from "#lib/components/ui/dropdown-menu/index.js";
  import Check from "@lucide/svelte/icons/check";
  import SlidersHorizontal from "@lucide/svelte/icons/sliders-horizontal";
  import Button from "../ui/button/button.svelte";
  export type SortDirection = "asc" | "desc";
  interface SortOption<T> {
    value: T;
    direction: SortDirection;
    items: T[];
  }
  let {
    value = $bindable(),
    items,
    direction = $bindable("asc"),
  }: SortOption<any> = $props();
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger>
    <Button variant="outline"><SlidersHorizontal /> {value}</Button>
  </DropdownMenu.Trigger>
  <DropdownMenu.Content>
    <DropdownMenu.Group>
      {#each items as item}
        <DropdownMenu.Item
          class="flex justify-between items-center"
          onclick={() => {
            value = item;
          }}
          ><span>{item}</span>{#if item === value}
            <Check />
          {/if}</DropdownMenu.Item
        >
      {/each}
    </DropdownMenu.Group>
    <DropdownMenu.Separator />
    <DropdownMenu.Group>
      <DropdownMenu.Item
        class="flex justify-between items-center"
        onclick={() => {
          direction = "asc";
        }}
      >
        <span>Ascending</span>
        {#if direction === "asc"}
          <Check />
        {/if}
      </DropdownMenu.Item>
      <DropdownMenu.Item
        class="flex justify-between items-center"
        onclick={() => {
          direction = "desc";
        }}
      >
        <span>Descending</span>

        {#if direction === "desc"}
          <Check />
        {/if}
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>
