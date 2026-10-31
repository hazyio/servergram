<script lang="ts">
  import * as DropdownMenu from "#lib/components/ui/dropdown-menu/index.js";
  import Check from "@lucide/svelte/icons/check";
  import SlidersHorizontal from "@lucide/svelte/icons/sliders-horizontal";
  import Button from "../ui/button/button.svelte";
  export type SortDirection = "asc" | "desc";
  interface SortOption<T> {
    value: T;
    direction?: SortDirection;
    items: T[];
    capitalize?: boolean;
    caption?: string;
    onChange?: (value: T, direction?: SortDirection) => void;
  }
  let {
    value = $bindable(),
    items,
    direction = $bindable(undefined),
    capitalize = true,
    caption = "Sort",
    onChange,
  }: SortOption<any> = $props();
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger title={caption}>
    <Button variant="outline"
      ><SlidersHorizontal />
      {capitalize
        ? value.charAt(0).toUpperCase() + value.slice(1)
        : value}</Button
    >
  </DropdownMenu.Trigger>
  <DropdownMenu.Content>
    <DropdownMenu.Group>
      {#each items as item}
        <DropdownMenu.Item
          class="flex justify-between items-center"
          onclick={() => {
            value = item;
            onChange?.(value, direction);
          }}
          ><span
            >{capitalize
              ? item.charAt(0).toUpperCase() + item.slice(1)
              : item}</span
          >{#if item === value}
            <Check />
          {/if}</DropdownMenu.Item
        >
      {/each}
    </DropdownMenu.Group>
    {#if direction !== undefined}
      <DropdownMenu.Separator />

      <DropdownMenu.Group>
        <DropdownMenu.Item
          class="flex justify-between items-center"
          onclick={() => {
            direction = "asc";
            onChange?.(value, direction);
          }}
        >
          <span>{capitalize ? "Ascending" : "Ascending"}</span>
          {#if direction === "asc"}
            <Check />
          {/if}
        </DropdownMenu.Item>
        <DropdownMenu.Item
          class="flex justify-between items-center"
          onclick={() => {
            direction = "desc";
            onChange?.(value, direction);
          }}
        >
          <span>{capitalize ? "Descending" : "Descending"}</span>

          {#if direction === "desc"}
            <Check />
          {/if}
        </DropdownMenu.Item>
      </DropdownMenu.Group>
    {/if}
  </DropdownMenu.Content>
</DropdownMenu.Root>
