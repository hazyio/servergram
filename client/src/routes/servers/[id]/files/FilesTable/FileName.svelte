<script lang="ts">
  import { buttonVariants } from "#lib/components/ui/button/index.js";
  import * as ContextMenu from "#lib/components/ui/context-menu/index.js";
  import { getFileIcon } from "#lib/fileIcons.ts";
  import type { DeleteEntryProps } from "../DeleteEntriesDialog.svelte";

  import { ArrowDownToLine, Eye, FolderOpen, Info } from "@lucide/svelte";
  import PencilIcon from "@lucide/svelte/icons/pencil";
  import TrashIcon from "@lucide/svelte/icons/trash";
  let {
    name,
    isDir,
    onDelete,
    onRenameClicked,
    onOpenPropertiesClicked,
    onViewClicked,
    onDownloadClicked,
  }: {
    name: string;
    isDir: boolean;
    onDelete: (entry: DeleteEntryProps) => void;
    onRenameClicked: (entry: DeleteEntryProps) => void;
    onOpenPropertiesClicked: (entry: DeleteEntryProps) => void;
    onViewClicked: (entry: DeleteEntryProps) => void;
    onDownloadClicked: (entry: DeleteEntryProps) => void;
  } = $props();
  const icon = $derived(getFileIcon(name, isDir));
</script>

{#if name === ".."}
  <div class="flex items-center gap-3 h-15 max-w-50">
    <img src={icon} alt="" class="size-6 shrink-0" />
    <span class="min-h-0 wrap-break-word"> {name}</span>
  </div>
{:else}
  <ContextMenu.Root>
    <ContextMenu.Trigger>
      <div class="flex items-center gap-3 h-15 max-w-50">
        <img src={icon} alt="" class="size-6 shrink-0" />
        <span class="min-h-0 wrap-break-word"> {name}</span>
      </div>
    </ContextMenu.Trigger>
    <ContextMenu.Content>
      <ContextMenu.Group>
        {#if isDir}
          <ContextMenu.Item class="py-2 my-1">
            <FolderOpen />
            Open
          </ContextMenu.Item>
        {:else}
          <ContextMenu.Item
            class="py-2 my-1"
            onclick={() => onViewClicked({ name, isDir })}
          >
            <Eye />
            View
          </ContextMenu.Item>
        {/if}
        <ContextMenu.Item
          class="py-2 my-1"
          onclick={() => onRenameClicked({ name, isDir })}
        >
          <PencilIcon />
          Rename
        </ContextMenu.Item>

        {#if !isDir}
          <ContextMenu.Item
            class="py-2 my-1"
            onclick={() => onDownloadClicked({ name, isDir })}
          >
            <ArrowDownToLine />
            Download
          </ContextMenu.Item>
        {/if}
        <ContextMenu.Item
          class="py-2 my-1"
          onclick={() => onOpenPropertiesClicked({ name, isDir })}
        >
          <Info />
          Properties
        </ContextMenu.Item>
      </ContextMenu.Group>
      <ContextMenu.Separator />
      <ContextMenu.Group>
        <ContextMenu.Item
          onclick={() => onDelete({ name, isDir })}
          variant="destructive"
          class={`gap-2 w-full justify-start ${buttonVariants({ variant: "destructive" })}`}
        >
          <TrashIcon class="text-destructive! hover:text-destructive!" />
          <span class="text-destructive!">Delete</span>
        </ContextMenu.Item>
      </ContextMenu.Group>
    </ContextMenu.Content>
  </ContextMenu.Root>
{/if}
