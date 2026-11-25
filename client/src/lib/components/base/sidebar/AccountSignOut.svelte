<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import appPaths from "#lib/appPaths.ts";
  import Error from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import Button from "#lib/components/ui/button/button.svelte";
  import { goto } from "$app/navigation";

  let open = $state(false);
  let signingOut = $state(false);
  let signOutError = $state("");

  async function signOut() {
    signingOut = true;
    signOutError = "";
    const { error } = await NewRequestClient().DELETE("/api/account/logout");
    signingOut = false;

    if (error) {
      signOutError = unWrapErrorMessage(
        error,
        "An error occurred while signing out",
      );
      return;
    }
    open = false;
    await goto(appPaths.account.login, { invalidateAll: true });
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Trigger>
    {#snippet child({ props })}
      <Button {...props} variant="destructive" class="grow">Sign Out</Button>
    {/snippet}
  </AlertDialog.Trigger>
  <AlertDialog.Content>
    {#if signingOut}
      <Loading class="h-50" text="Signing out..." loading={signingOut} />
    {:else if signOutError !== ""}
      <Error message={signOutError} />
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={signOut} variant="destructive"
          >Try Again</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {:else}
      <AlertDialog.Header>
        <AlertDialog.Title>Are you sure you want to sign out?</AlertDialog.Title
        >
        <AlertDialog.Description>
          You will need to enter your password again to access your resources.
        </AlertDialog.Description>
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <AlertDialog.Action onclick={signOut} variant="destructive"
          >Sign Out</AlertDialog.Action
        >
      </AlertDialog.Footer>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
