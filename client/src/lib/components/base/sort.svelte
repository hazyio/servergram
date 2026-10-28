<script lang="ts">
  import * as DropdownMenu from "#lib/components/ui/dropdown-menu/index.js";
  import Check from "@lucide/svelte/icons/check";
  import SlidersHorizontal from "@lucide/svelte/icons/sliders-horizontal";
  import ListSortAscending from "@lucide/svelte/icons/list-sort-ascending";
  import ListSortDescending from "@lucide/svelte/icons/list-sort-descending";
  import Button from "../ui/button/button.svelte";
  export type SortDirection = "asc" | "desc";
  interface SortOption<T> {
    value: T;
    direction: SortDirection;
    items: T[];
    onSelect?: (value: T, direction: SortDirection) => void;
  }
  let { value, items, direction, onSelect }: SortOption<any> = $props();
</script>

<div>
  <DropdownMenu.Root>
    <DropdownMenu.Trigger>
      <Button variant="outline"><SlidersHorizontal /> {value}</Button>
    </DropdownMenu.Trigger>
    <DropdownMenu.Content>
      <DropdownMenu.Group>
        {#each items as item}
          <DropdownMenu.Item
            class="flex justify-between items-center"
            onclick={() => onSelect?.(item, direction)}
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
          onclick={() => onSelect?.(value, "asc")}
        >
          <span>Ascending</span>
          {#if direction === "asc"}
            <Check />
          {/if}
        </DropdownMenu.Item>
        <DropdownMenu.Item
          class="flex justify-between items-center"
          onclick={() => onSelect?.(value, "desc")}
        >
          <span>Descending</span>

          {#if direction === "desc"}
            <Check />
          {/if}
        </DropdownMenu.Item>
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>
</div>
