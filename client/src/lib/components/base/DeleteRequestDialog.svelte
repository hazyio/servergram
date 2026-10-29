<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import Error from "#lib/components/base/error.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import { buttonVariants } from "#lib/components/ui/button/index.js";
  let deleting = $state(false);
  let deleteError = $state("");
  let open = $state(false);
  let { id, name, onDeleted, type } = $props();
  async function deleteServer() {
    deleting = true;
    const { error } = await NewRequestClient().DELETE(
      `/api/${type}/delete/{id}` as any,
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
        `An error occurred while deleting the ${type}`,
      );
      return;
    }
    open = false;
    onDeleted?.();
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Trigger class={buttonVariants({ variant: "destructive" })}>
    Delete
  </AlertDialog.Trigger>
  <AlertDialog.Content>
    {#if deleting}
      <Loading class="h-50" text={`Deleting ${type}...`} loading={deleting} />
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
          >Are you absolutely sure you want to delete {type} "{name}"?</AlertDialog.Title
        >
        <AlertDialog.Description>
          This action cannot be undone. This will permanently delete your {type} "{name}".
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
