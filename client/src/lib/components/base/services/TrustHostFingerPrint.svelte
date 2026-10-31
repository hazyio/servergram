<script lang="ts">
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import type { WSServiceTrustHostFingerPrintData } from "#lib/types/websockets/index.js";
  import { Button } from "../../ui/button";
  let {
    data = $bindable([]),
  }: {
    data: WSServiceTrustHostFingerPrintData[];
  } = $props();
</script>

{#if data.length > 0}
  <Dialog.Root open={true}>
    <Dialog.Content
      showCloseButton={false}
      onEscapeKeydown={(event) => event.preventDefault()}
      onInteractOutside={(event) => event.preventDefault()}
    >
      <Dialog.Header>
        <Dialog.Title>Trust the SSH host fingerprint?</Dialog.Title>
        <Dialog.Description>
          Do you want to trust the host fingerprint for server {data[0].name}?
        </Dialog.Description>
      </Dialog.Header>
      <pre
        class="max-h-80 overflow-auto rounded-md bg-muted p-3 text-xs break-all whitespace-pre-wrap">{data[0]
          .finger_print}</pre>
      <Dialog.Footer>
        <Button
          onclick={() => {
            data = data.slice(1);
          }}>Don't trust</Button
        >
        <Button onclick={() => (data = data.slice(1))}>Trust host</Button>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
{/if}
