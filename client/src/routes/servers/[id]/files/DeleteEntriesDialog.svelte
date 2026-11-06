<script lang="ts" module>
  export interface DeleteEntryProps {
    name: string;
    isDir: boolean;
  }
</script>

<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";

  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import { Badge } from "#lib/components/ui/badge/index.ts";
  let deleting = $state(false);
  let deleteError = $state("");
  let {
    serverId,

    dir,
    onDeleted,
    items = $bindable(),
  }: {
    dir: string;
    serverId: number;

    onDeleted?: (name: string) => void;
    items: DeleteEntryProps[];
  } = $props();

  // Compute the full paths for display
  const paths = $derived(
    items.map((item) =>
      `${dir.replace(/\/+$/, "")}/${item.name}`.replace(/\/+/g, "/"),
    ),
  );
  const isMultiple = $derived(items.length > 1);
  const itemType = $derived(items[0]?.isDir ? "directory" : "file");
  let deletingName = $state("");
  async function deleteAction() {
    deleting = true;
    for (const entry of items) {
      deletingName = `${entry.isDir ? "Directory" : "File"} "${entry.name}"`;
      const { error } = await (async () => {
        if (entry.isDir) {
          return await NewRequestClient().DELETE(
            "/api/files/remove-dir-in-server/{id}",
            {
              body: {
                path: `${dir.replace(/\/+$/, "")}/${entry.name}`.replace(
                  /\/+/g,
                  "/",
                ),
              },
              params: {
                path: {
                  id: serverId,
                },
              },
            },
          );
        } else {
          return await NewRequestClient().DELETE(
            "/api/files/remove-file-in-server/{id}",
            {
              body: {
                name: entry.name,
                path: dir,
              },
              params: {
                path: {
                  id: serverId,
                },
              },
            },
          );
        }
      })();

   
      if (error) {
        deleting = false;
        deleteError = unWrapErrorMessage(
          error,
          `Failed to delete ${entry.isDir ? "directory" : "file"}`,
        );
        return;
      }
      if (onDeleted) {
        // remove file from list
        items = items.filter((i) => i.name !== entry.name);
        onDeleted(entry.name);
      }
    }
    deleting = false;
  }
</script>

<AlertDialog.Root
  open={items.length > 0}
  onOpenChange={(open) => {
    if (!open) items = [];
  }}
>
  <AlertDialog.Content>
    {#if deleting}
      <Loading
        class="h-50"
        text={`Deleting ${deletingName}`}
        loading={deleting}
      />
    {:else if deleteError !== ""}
      <Error message={deleteError} />
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={deleteAction} variant="destructive"
          >Try Again</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {:else}
      <AlertDialog.Header>
        <AlertDialog.Title>
          {#if isMultiple}
            Are you absolutely sure you want to delete these {items.length}
            entries?
          {:else}
            Are you absolutely sure you want to delete this {itemType}?
          {/if}
        </AlertDialog.Title>
        <AlertDialog.Description>
          {#if isMultiple}
            This action cannot be undone. This will permanently delete the
            following {items.length} entries:
            <ul class="mt-2 ml-1 list-disc text-sm max-h-40 overflow-auto">
              {#each paths as path}
                <li class="truncate" title={path}>
                  {path}
                  <Badge variant="outline"
                    >{itemType === "directory" ? "Directory" : "File"}</Badge
                  >
                </li>
              {/each}
            </ul>
          {:else}
            This action cannot be undone. This will permanently delete:
            <div
              class="bg-muted p-2 flex gap-1 mt-2 ml-1 font-mono text-sm break-all"
            >
              {paths[0]}<Badge variant="outline"
                >{itemType === "directory" ? "Directory" : "File"}</Badge
              >
            </div>
          {/if}
        </AlertDialog.Description>
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel
          onclick={() => {
            items = [];
          }}>Cancel</AlertDialog.Cancel
        >
        <AlertDialog.Action onclick={deleteAction} variant="destructive"
          >Continue</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
