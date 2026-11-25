<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import appPaths from "#lib/appPaths.ts";
  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import Button from "#lib/components/ui/button/button.svelte";
  import { goto } from "$app/navigation";
  import LogOut from "@lucide/svelte/icons/log-out";
  import Trash from "@lucide/svelte/icons/trash";

  let {
    id,
    os,
    current = false,
    onRemoved,
  }: {
    id: number;
    os?: string;
    current?: boolean;
    onRemoved?: () => void;
  } = $props();

  let open = $state(false);
  let submitting = $state(false);
  let actionError = $state("");

  async function run() {
    submitting = true;
    actionError = "";
    const { error } = current
      ? await NewRequestClient().DELETE("/api/account/logout")
      : await NewRequestClient().DELETE("/api/account/sessions/remove/{id}", {
          params: {
            path: {
              id: String(id),
            },
          },
        });
    submitting = false;

    if (error) {
      actionError = unWrapErrorMessage(
        error,
        current
          ? "An error occurred while signing out"
          : "An error occurred while removing the session",
      );
      return;
    }
    open = false;

    if (current) {
      await goto(appPaths.account.login, { invalidateAll: true });
    } else {
      onRemoved?.();
    }
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Trigger>
    {#snippet child({ props })}
      <Button {...props} variant="destructive" size="sm" class="w-full">
        {#if current}
          <LogOut />
          Sign Out
        {:else}
          <Trash />
          Remove
        {/if}
      </Button>
    {/snippet}
  </AlertDialog.Trigger>
  <AlertDialog.Content>
    {#if submitting}
      <Loading
        class="h-50"
        text={current ? "Signing out..." : "Removing session..."}
        loading={submitting}
      />
    {:else if actionError !== ""}
      <Error message={actionError} />
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={run} variant="destructive"
          >Try Again</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {:else}
      <AlertDialog.Header>
        {#if current}
          <AlertDialog.Title>Sign out this device?</AlertDialog.Title>
          <AlertDialog.Description>
            You will need to enter your password again to access your resources.
          </AlertDialog.Description>
        {:else}
          <AlertDialog.Title>
            Remove {os ? `the ${os}` : "this"} session?
          </AlertDialog.Title>
          <AlertDialog.Description>
            This will sign out the device associated with this session. This
            action cannot be undone.
          </AlertDialog.Description>
        {/if}
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={run} variant="destructive">
          {current ? "Sign Out" : "Remove"}
        </AlertDialog.Action>
      </AlertDialog.Footer>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
