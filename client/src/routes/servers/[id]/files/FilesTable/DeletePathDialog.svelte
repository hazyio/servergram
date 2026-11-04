<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import * as ContextMenu from "#lib/components/ui/context-menu/index.js";

  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import { buttonVariants } from "#lib/components/ui/button/index.js";
  import { Trash, TrashIcon } from "@lucide/svelte";
  import type { Snippet } from "svelte";
  let deleting = $state(false);
  let deleteError = $state("");
  let open = $state(false);
  let {
    id,
    name,
    onDeleted,
  }: {
    id: number;
    name: string;
    onDeleted?: () => void;
  } = $props();
  async function deleteServer() {
    deleting = true;
    const { error } = await NewRequestClient().DELETE(
      `/api/group/delete/{id}`,
      {
        params: {
          path: {
            id: id,
          },
        },
      },
    );
    deleting = false;

    if (error) {
      deleteError = unWrapErrorMessage(
        error,
        `An error occurred while deleting the path`,
      );
      return;
    }
    open = false;
    onDeleted?.();
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Trigger class="w-full">
    <ContextMenu.Item
      variant="destructive"
      class={`gap-2 w-full ${buttonVariants({ variant: "destructive" })}`}
    >
      <TrashIcon class="text-destructive! hover:text-destructive!" />
      <span class="text-destructive!">Delete</span>
    </ContextMenu.Item>
  </AlertDialog.Trigger>
  <AlertDialog.Content>
    {#if deleting}
      <Loading class="h-50" text={`Deleting Path ...`} loading={deleting} />
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
        <AlertDialog.Title
          >Are you absolutely sure you want to delete the path "{name}"?</AlertDialog.Title
        >
        <AlertDialog.Description>
          This action cannot be undone. This will permanently delete your path "{name}".
        </AlertDialog.Description>
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={deleteServer} variant="destructive"
          >Continue</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
