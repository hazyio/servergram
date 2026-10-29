<script lang="ts">
  import ChevronDownIcon from "@lucide/svelte/icons/chevron-down";
  import { onMount, tick } from "svelte";
  import * as Command from "#lib/components/ui/command/index.js";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import { Button } from "#lib/components/ui/button/index.js";
  import CreateSHHKey from "./createSHHKey.svelte";
  import { NewRequestClient } from "#lib/api/request.ts";
  import type { components } from "#lib/api/schema.js";
  import Badge from "#lib/components/ui/badge/badge.svelte";

  let availableSSHKeys: components["schemas"]["SSHKeySearchData"][] = $state(
    [],
  );
  let searchQuery = $state("");
  let open = $state(false);
  let triggerRef = $state<HTMLButtonElement>(null!);
  let timeoutId = $state(0);
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
  }: { value: number; onSelect: (value: number, encrypted: boolean) => void } =
    $props();
  async function fetchSSHKeys(query: string) {
    const mySeq = ++seq;
    const { data } = await NewRequestClient().POST("/api/search/ssh", {
      body: { query },
    });
    if (mySeq !== seq) return; // a newer request started
    if (data) {
      if (data.items === null) {
        availableSSHKeys = [];
      }
      if (data?.items)
        availableSSHKeys.splice(0, availableSSHKeys.length, ...data.items);
    }
  }
  async function debounceFetchSSHKeys(query: string) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(async () => {
      await fetchSSHKeys(query);
    }, 300);
  }
  onMount(() => {
    fetchSSHKeys(searchQuery);
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
        {selectedName ? selectedName : "Select an SSH key"}
        <ChevronDownIcon class="text-muted-foreground" />
      </Button>
    {/snippet}
  </Popover.Trigger>
  <Popover.Content class="w-(--bits-popover-anchor-width) p-0">
    <Command.Root shouldFilter={false}>
      <div class="w-full flex gap-2 justify-center items-center">
        <Command.Input
          placeholder="Search SSH key..."
          class="p-0 grow"
          bind:value={searchQuery}
          oninput={() => {
            debounceFetchSSHKeys(searchQuery);
          }}
        />
        <CreateSHHKey text="New" />
      </div>
      <Command.List>
        <Command.Empty>No items found.</Command.Empty>
        <Command.Group>
          {#each availableSSHKeys as sshKey (sshKey.id)}
            <Command.Item
              value={sshKey.id.toString()}
              data-checked={value === sshKey.id}
              onSelect={() => {
                selectedName = sshKey.name;
                onSelect(sshKey.id, sshKey.encrypted);
                closeAndFocusTrigger();
              }}
            >
              {sshKey.name}
              <Badge variant={sshKey.encrypted ? "default" : "outline"}
                >{sshKey.encrypted ? "encrypted" : "not encrypted"}</Badge
              >
            </Command.Item>
          {/each}
        </Command.Group>
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
