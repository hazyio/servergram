<script lang="ts">
  import * as ContextMenu from "#lib/components/ui/context-menu/index.js";
  import { getFileIcon } from "#lib/fileIcons.ts";

  import DeletePathDialog from "./DeletePathDialog.svelte";
  import { ArrowDownToLine, Eye, FolderOpen, Info } from "@lucide/svelte";
  import PencilIcon from "@lucide/svelte/icons/pencil";
  import TrashIcon from "@lucide/svelte/icons/trash";
  let { name, isDir }: { name: string; isDir: boolean } = $props();
  const icon = $derived(getFileIcon(name, isDir));
</script>

<ContextMenu.Root>
  <ContextMenu.Trigger>
    <span class="flex items-center gap-3 h-15">
      <img src={icon} alt="" class="size-6 shrink-0" />
      {name}
    </span>
  </ContextMenu.Trigger>
  <ContextMenu.Content>
    <ContextMenu.Group>
      {#if isDir}
        <ContextMenu.Item>
          <FolderOpen />
          Open
        </ContextMenu.Item>
      {:else}
        <ContextMenu.Item>
          <Eye />
          View
        </ContextMenu.Item>
      {/if}
      <ContextMenu.Item>
        <PencilIcon />
        Rename
      </ContextMenu.Item>

      {#if !isDir}
        <ContextMenu.Item>
          <ArrowDownToLine />
          Download
        </ContextMenu.Item>
      {/if}
      <ContextMenu.Item>
        <Info />
        Properties
      </ContextMenu.Item>
    </ContextMenu.Group>
    <ContextMenu.Separator />
    <ContextMenu.Group>
      <DeletePathDialog id={10000} {name} />
    </ContextMenu.Group>
  </ContextMenu.Content>
</ContextMenu.Root>
