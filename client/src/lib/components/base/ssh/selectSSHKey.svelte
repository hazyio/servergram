<script lang="ts">
  import ChevronDownIcon from "@lucide/svelte/icons/chevron-down";
  import { tick } from "svelte";
  import * as Command from "#lib/components/ui/command/index.js";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import { Button } from "#lib/components/ui/button/index.js";

  const frameworks = ["Next.js", "SvelteKit", "Nuxt.js", "Remix", "Astro"];

  let open = $state(false);
  let triggerRef = $state<HTMLButtonElement>(null!);

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
  }: { value: string; onSelect: (value: string) => void } = $props();
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
        {value || "Select an SSH key"}
        <ChevronDownIcon class="text-muted-foreground" />
      </Button>
    {/snippet}
  </Popover.Trigger>
  <Popover.Content class="w-(--bits-popover-anchor-width) p-0">
    <Command.Root>
      <div class="w-full flex gap-2 justify-center items-center">
        <Command.Input placeholder="Search SSH key..." class="p-0 grow" />
        <Button>New</Button>
      </div>
      <Command.List>
        <Command.Empty>No items found.</Command.Empty>
        <Command.Group>
          {#each frameworks as framework (framework)}
            <Command.Item
              value={framework}
              data-checked={value === framework}
              onSelect={() => {
                onSelect(framework);
                closeAndFocusTrigger();
              }}
            >
              {framework}
            </Command.Item>
          {/each}
        </Command.Group>
      </Command.List>
    </Command.Root>
  </Popover.Content>
</Popover.Root>
