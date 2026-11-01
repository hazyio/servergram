<script lang="ts">
  import DeleteRequestDialog from "#lib/components/base/DeleteRequestDialog.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import { ConnectedServices } from "#lib/components/base/services/services.svelte.ts";
  import Loading from "#lib/components/base/loading.svelte";
  import DisconnectService from "#lib/components/base/services/DisconnectService.svelte";

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

<div class="flex items-center justify-end gap-3">
  {#if serverConnected}
    {#if serverConnected.status === "connecting"}
      <Loading loading class="w-4 h-4" containerClass="w-12" />
    {:else}
      <DisconnectService {id} {name} type="server" />
    {/if}
  {:else}
    <Button
      variant="secondary"
      onclick={() => ConnectedServices.startServerConnection(id, name)}
      >Connect</Button
    >
    <DeleteRequestDialog {id} {name} {onDeleted} type="server" useIcon />
  {/if}
</div>
