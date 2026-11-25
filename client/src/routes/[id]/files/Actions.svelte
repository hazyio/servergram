<script lang="ts">
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import {
    EllipsisVertical,
    FilePlus,
    FolderPlus,
    Upload,
  } from "@lucide/svelte";
  import type { NewEntryType } from "./NewEntry.svelte";
  let {
    entryType = $bindable(),
    onUploadClicked,
  }: {
    entryType: NewEntryType;
    onUploadClicked: () => void;
  } = $props();
</script>

<Popover.Root>
  <Popover.Trigger>
    <Button class="rounded-md" size="icon">
      <EllipsisVertical />
    </Button>
  </Popover.Trigger>

  <Popover.Content align="end" class="min-w-40">
    <Popover.Header>Actions</Popover.Header>

    <div class="flex flex-col gap-1">
      <Popover.Close>
        <Button
          variant="ghost"
          class="justify-start py-5 cursor-pointer gap-2 w-full"
          onclick={() => {
            entryType = "file";
          }}
        >
          <FilePlus class="size-4" />
          New File
        </Button>
      </Popover.Close>
      <Popover.Close>
        <Button
          variant="ghost"
          class="justify-start py-5 cursor-pointer gap-2 w-full"
          onclick={() => {
            entryType = "directory";
          }}
        >
          <FolderPlus class="size-4" />
          New Directory
        </Button>
      </Popover.Close>
      <Popover.Close>
        <Button
          variant="ghost"
          class="justify-start py-5 cursor-pointer gap-2 w-full"
          onclick={onUploadClicked}
        >
          <Upload class="size-4" />
          Upload
        </Button>
      </Popover.Close>
    </div>
  </Popover.Content>
</Popover.Root>
