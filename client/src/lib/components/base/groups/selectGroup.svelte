<script lang="ts">
  import ChevronDownIcon from "@lucide/svelte/icons/chevron-down";
  import { onMount, tick } from "svelte";
  import * as Command from "#lib/components/ui/command/index.js";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import { Button } from "#lib/components/ui/button/index.js";
  import CreateGroup from "./createGroup.svelte";
  import { NewRequestClient } from "#lib/api/request.ts";
  import type { components } from "#lib/api/schema.js";

  let availableGroups: components["schemas"]["GroupSearchData"][] = $state([]);
  let searchQuery = $state("");
  let open = $state(false);
  let triggerRef = $state<HTMLButtonElement>(null!);
  let timeoutId: NodeJS.Timeout = $state(null!);
  let seq = $state(0);
  let selectedName = $state("");

  // We want to refocus the trigger button when the user selects
  // an item from the list so users can continue navigating the
  // rest of the form with the keyboard.
  function closeAndFocusTrigger() {
    open = false;
    tick().then(() => {
      triggerRef.focus();
    });
  }
  let {
    value,
    onSelect,
  }: { value: number; onSelect: (value: number) => void } = $props();
  async function fetchGroups(query: string) {
    const mySeq = ++seq;
    const { data } = await NewRequestClient().GET("/api/search/group", {
      params: { query: { query } },
    });
    if (mySeq !== seq) return; // a newer request started
    if (data) {
      if (data.items === null) {
        availableGroups = [];
      }
      if (data?.items)
        availableGroups.splice(0, availableGroups.length, ...data.items);
    }
  }
  async function debounceFetchGroups(query: string) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(async () => {
      await fetchGroups(query);
    }, 300);
  }
  onMount(() => {
    fetchGroups(searchQuery);
  });
</script>

<Popover.Root bind:open>
  <Popover.Trigger bind:ref={triggerRef}>
    {#snippet child({ props })}
      <Button
        {...props}
        variant="outline"
        class="w-full justify-between font-normal"
        role="combobox"
        aria-expanded={open}
      >
        {selectedName ? selectedName : "Select a group"}
        <ChevronDownIcon class="text-muted-foreground" />
      </Button>
    {/snippet}
  </Popover.Trigger>
  <Popover.Content class="w-(--bits-popover-anchor-width) p-0">
    <Command.Root shouldFilter={false}>
      <div class="w-full flex gap-2 justify-center items-center">
        <Command.Input
          placeholder="Search group..."
          class="p-0 grow"
          bind:value={searchQuery}
          oninput={() => {
            debounceFetchGroups(searchQuery);
          }}
        />
        <CreateGroup text="New" />
      </div>
      <Command.List>
        <Command.Empty>No items found.</Command.Empty>
        <Command.Group>
          {#each availableGroups as group (group.id)}
            <Command.Item
              value={group.id.toString()}
              data-checked={value === group.id}
              onSelect={() => {
                selectedName = group.name;
                onSelect(group.id);
                closeAndFocusTrigger();
              }}
            >
              {group.name}
            </Command.Item>
          {/each}
        </Command.Group>
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
