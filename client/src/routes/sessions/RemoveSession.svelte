<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import Button from "#lib/components/ui/button/button.svelte";
  import { Trash } from "@lucide/svelte";

  let {
    id,
    os,
    onDeleted,
  }: {
    id: number;
    os?: string;
    onDeleted: () => void;
  } = $props();

  let open = $state(false);
  let removing = $state(false);
  let removeError = $state("");

  async function removeSession() {
    removing = true;
    removeError = "";
    const { error } = await NewRequestClient().DELETE(
      "/api/account/sessions/remove/{id}",
      {
        params: {
          path: {
            id: String(id),
          },
        },
      },
    );
    removing = false;

    if (error) {
      removeError = unWrapErrorMessage(
        error,
        "An error occurred while removing the session",
      );
      return;
    }
    open = false;
    onDeleted();
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Trigger>
    {#snippet child({ props })}
      <Button
        {...props}
        variant="destructive"
        size="icon"
        title="Remove session"
      >
        <Trash />
        <span class="sr-only">Remove</span>
      </Button>
    {/snippet}
  </AlertDialog.Trigger>
  <AlertDialog.Content>
    {#if removing}
      <Loading class="h-50" text="Removing session..." loading={removing} />
    {:else if removeError !== ""}
      <Error message={removeError} />
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={removeSession} variant="destructive"
          >Try Again</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {:else}
      <AlertDialog.Header>
        <AlertDialog.Title>
          Are you sure you want to remove {os ? `the ${os}` : "this"} session?
        </AlertDialog.Title>
        <AlertDialog.Description>
          This will sign out the device associated with this session. This
          action cannot be undone.
        </AlertDialog.Description>
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={removeSession} variant="destructive"
          >Remove</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
