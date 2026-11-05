<script lang="ts" module>
  export interface DeleteEntryProps {
    name: string;
    isDir: boolean;
  }
</script>

<script lang="ts">
  import * as ContextMenu from "#lib/components/ui/context-menu/index.js";

  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import { buttonVariants } from "#lib/components/ui/button/index.js";
  import { TrashIcon } from "@lucide/svelte";
  let deleting = $state(false);
  let deleteError = $state("");
  let {
    dir,
    onDeleted,
    items = $bindable(),
  }: {
    dir: string;
    onDeleted?: () => void;
    items: DeleteEntryProps[];
  } = $props();

  // Compute the full paths for display
  const paths = $derived(
    items.map((item) =>
      `${dir.replace(/\/+$/, "")}/${item.name}`.replace(/\/+/g, "/"),
    ),
  );
  const isMultiple = $derived(items.length > 1);
  const itemType = $derived(items[0]?.isDir ? "folder" : "file");

  async function deleteServer() {
    deleting = true;
  }
</script>

<AlertDialog.Root open={items.length > 0}>
  <AlertDialog.Trigger class="w-full">
    <ContextMenu.Item
      variant="destructive"
      class={`gap-2 w-full justify-start ${buttonVariants({ variant: "destructive" })}`}
    >
      <TrashIcon class="text-destructive! hover:text-destructive!" />
      <span class="text-destructive!">Delete</span>
    </ContextMenu.Item>
  </AlertDialog.Trigger>
  <AlertDialog.Content>
    {#if deleting}
      <Loading
        class="h-50"
        text={`Deleting ${isMultiple ? `${items.length} items` : paths[0]} ...`}
        loading={deleting}
      />
    {:else if deleteError !== ""}
      <Error message={deleteError} />
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={deleteServer} variant="destructive"
          >Try Again</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {:else}
      <AlertDialog.Header>
        <AlertDialog.Title>
          {#if isMultiple}
            Are you absolutely sure you want to delete these {items.length}
            {itemType}{items.length > 1 ? "s" : ""}?
          {:else}
            Are you absolutely sure you want to delete this {itemType}?
          {/if}
        </AlertDialog.Title>
        <AlertDialog.Description>
          {#if isMultiple}
            This action cannot be undone. This will permanently delete the
            following {items.length}
            {itemType}{items.length > 1 ? "s" : ""}:
            <ul class="mt-2 ml-4 list-disc text-sm max-h-40 overflow-auto">
              {#each paths as path}
                <li class="truncate" title={path}>{path}</li>
              {/each}
            </ul>
          {:else}
            This action cannot be undone. This will permanently delete:
            <div class="mt-2 ml-4 font-mono text-sm break-all">{paths[0]}</div>
          {/if}
        </AlertDialog.Description>
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel
          onclick={() => {
            items = [];
          }}>Cancel</AlertDialog.Cancel
        >
        <AlertDialog.Action onclick={deleteServer} variant="destructive"
          >Continue</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
