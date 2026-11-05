<script lang="ts">
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Popover from "#lib/components/ui/popover/index.js";
  import {
    EllipsisVertical,
    FolderPlus,
    FilePlus,
    RefreshCcw,
    Upload,
    Download,
  } from "@lucide/svelte";
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
      <Button
        variant="ghost"
        class="justify-start py-4 gap-2 w-full"
        onclick={() => {}}
      >
        <FolderPlus class="size-4" />
        New Folder
      </Button>
      <Button
        variant="ghost"
        class="justify-start py-4 gap-2 w-full"
        onclick={() => {}}
      >
        <FilePlus class="size-4" />
        New File
      </Button>
      <Button
        variant="ghost"
        class="justify-start py-4 gap-2 w-full"
        onclick={() => {}}
      >
        <Upload class="size-4" />
        Upload
      </Button>

      <Button
        variant="ghost"
        class="justify-start py-4 gap-2 w-full"
        onclick={() => {}}
      >
        <RefreshCcw class="size-4" />
        Refresh
      </Button>
    </div>
  </Popover.Content>
</Popover.Root>
