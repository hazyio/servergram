<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import type { WSServiceTrustHostFingerPrintData } from "#lib/types/websockets/index.js";
  import { Button } from "../../ui/button";
  import { toast } from "svelte-sonner";
  import { ConnectedServices } from "./services.svelte";

  let {
    data = $bindable([]),
  }: {
    data: WSServiceTrustHostFingerPrintData[];
  } = $props();
  let submitting = $state(false);
  let edit = $derived(data[0]);
  async function trustHost() {
    if (!edit) return;
    submitting = true;
    const { error } = await NewRequestClient().PUT(
      `/api/server/update-finger-print/{id}`,
      {
        params: {
          path: {
            id: edit.id,
          },
        },
        body: {
          fingerprint: edit.finger_print,
        },
      },
    );
    submitting = false;
    if (error) {
      toast.error(
        unWrapErrorMessage(error, "Failed to trust host fingerprint"),
      );
    } else {
      toast.success(
        `Successfully trusted host fingerprint for server ${edit.name}`,
      );
            ConnectedServices.startServerConnection(edit.id, edit.name);

      removed();
    }
  }
  function removed() {
    data = data.slice(1);
  }
</script>

{#if edit !== undefined}
  <Dialog.Root open={true}>
    <Dialog.Content
      showCloseButton={false}
      onEscapeKeydown={(event) => event.preventDefault()}
      onInteractOutside={(event) => event.preventDefault()}
    >
      <Dialog.Header>
        <Dialog.Title>Trust the SSH host fingerprint?</Dialog.Title>
        <Dialog.Description>
          Do you want to trust the host fingerprint for server {edit.name}?
        </Dialog.Description>
      </Dialog.Header>
      <pre
        class="max-h-80 overflow-auto rounded-md bg-muted p-3 text-xs break-all whitespace-pre-wrap">{edit.finger_print}</pre>
      <Dialog.Footer>
        <Button variant="outline" onclick={removed} disabled={submitting}>Don't trust</Button>
        <Button onclick={trustHost} disabled={submitting}>Trust</Button>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
{/if}
