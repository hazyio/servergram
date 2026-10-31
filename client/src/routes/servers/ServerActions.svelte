<script lang="ts">
  import DeleteRequestDialog from "#lib/components/base/DeleteRequestDialog.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import { ConnectedServices } from "#lib/components/base/services/services.svelte.ts";

  let {
    id,
    name,
    onDeleted,
  }: {
    id: number;
    name: string;
    onDeleted: () => void;
  } = $props();
  let serverConnected = $derived(ConnectedServices.findServer(id));
</script>

<div class="flex items-center justify-end gap-2">
  {#if serverConnected}
    <span>Connected</span>
  {/if}
  <Button
    variant="secondary"
    onclick={() => ConnectedServices.startServerConnection(id, name)}
    >Connect</Button
  >
  <DeleteRequestDialog {id} {name} {onDeleted} type="server" useIcon />
</div>
