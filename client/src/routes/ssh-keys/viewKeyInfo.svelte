<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import Error from "#lib/components/base/error.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import { buttonVariants } from "#lib/components/ui/button/index.js";
  import * as Dialog from "#lib/components/ui/dialog/index.ts";
  import { Textarea } from "#lib/components/ui/textarea/index.ts";
  import { Eye } from "@lucide/svelte";

  let loading = $state(false);
  let sshKey = $state("");
  let loadError = $state("");
  let { id } = $props();
  async function loadKey() {
    loading = true;
    const { data, error } = await NewRequestClient().GET(
      `/api/ssh-key/view-key/{id}`,
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
    onclick={loadKey}
    title="View SSH key"
    ><Eye />
    <span class="sr-only">View SSH key</span>
  </Dialog.Trigger>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>{loading ? "Loading..." : "SSH Key"}</Dialog.Title>
    </Dialog.Header>
    <Loading containerClass="mt-4 w-full h-50" {loading} />
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
