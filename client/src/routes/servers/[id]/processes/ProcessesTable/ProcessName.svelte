<script lang="ts">
  import * as ContextMenu from "#lib/components/ui/context-menu/index.js";

  import { FolderOpen } from "@lucide/svelte";
  let {
    name,
  }: {
    name: string;
  } = $props();
</script>

<ContextMenu.Root>
  <ContextMenu.Trigger>
    <div class="flex items-center gap-3 h-15 max-w-50">
      <span class="min-h-0 wrap-break-word"> {name}</span>
    </div>
  </ContextMenu.Trigger>
  <ContextMenu.Content>
    <ContextMenu.Group>
      <ContextMenu.Item class="py-2 my-1">
        <FolderOpen />
        Kill
      </ContextMenu.Item>
    </ContextMenu.Group>
  </ContextMenu.Content>
</ContextMenu.Root>
