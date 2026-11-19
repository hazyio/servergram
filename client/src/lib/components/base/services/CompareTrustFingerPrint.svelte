<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import * as Tabs from "#lib/components/ui/tabs/index.js";
  import type { WSServiceHostFingerPrintMismatchData } from "#lib/types/websockets/index.js";
  import { toast } from "svelte-sonner";
  import { Button } from "../../ui/button";
  import { ConnectedServices } from "./services.svelte";

  let {
    data = $bindable([]),
  }: {
    data: WSServiceHostFingerPrintMismatchData[];
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
          fingerprint: edit.new_finger_print,
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
        <Dialog.Title>Host fingerprint mismatch</Dialog.Title>
        <Dialog.Description>
          Host fingerprint mismatch for server {edit.name}. Do you want to trust
          the new fingerprint?
        </Dialog.Description>
      </Dialog.Header>
      <Tabs.Root value="old" class="w-full">
        <Tabs.List class="w-full">
          <Tabs.Trigger value="old">Old</Tabs.Trigger>
          <Tabs.Trigger value="new">New</Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content value="old">
          <pre
            class="max-h-80 overflow-auto rounded-md bg-muted p-3 text-xs break-all whitespace-pre-wrap">{edit.old_finger_print}</pre>
        </Tabs.Content>
        <Tabs.Content value="new">
          <pre
            class="max-h-80 overflow-auto rounded-md bg-muted p-3 text-xs break-all whitespace-pre-wrap">{edit.new_finger_print}</pre>
        </Tabs.Content>
      </Tabs.Root>

      <Dialog.Footer>
        <Button variant="outline" onclick={removed} disabled={submitting}
          >Don't trust</Button
        >
        <Button onclick={trustHost} disabled={submitting}>Trust</Button>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
{/if}
