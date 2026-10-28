<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import Error from "#lib/components/base/error.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.ts";
  import { Textarea } from "#lib/components/ui/textarea/index.ts";
  import { buttonVariants } from "#lib/components/ui/button/index.js";

  let loading = $state(false);
  let sshKey = $state("");
  let loadError = $state("");
  let { id } = $props();
  async function loadKey() {
    loading = true;
    const { data, error } = await NewRequestClient().GET(
      `/api/ssh/view-key/{id}`,
      {
        params: {
          path: {
            id: id,
          },
        },
      },
    );
    loading = false;
    if (error) {
      loadError =
        error.detail ??
        error.title ??
        "An error occurred while loading the SSH key";
      return;
    }
    sshKey = data.private;
  }
</script>

<Dialog.Root>
  <Dialog.Trigger
    class={buttonVariants({ variant: "outline" })}
    onclick={loadKey}>View Key</Dialog.Trigger
  >
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>{loading ? "Loading..." : "SSH Key"}</Dialog.Title>
    </Dialog.Header>
    <Loading class="mt-4 h-50 w-full" {loading} />
    {#if !loading && loadError}
      <Error message={loadError} />
    {/if}
    {#if !loading && sshKey}
      <Textarea
        id="key"
        class="min-h-70 max-h-[70vh] w-full resize-none field-sizing-fixed overflow-y-auto"
        bind:value={sshKey}
        readonly
      />
    {/if}
  </Dialog.Content>
</Dialog.Root>
