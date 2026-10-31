<script lang="ts">
  import DeleteRequestDialog from "#lib/components/base/DeleteRequestDialog.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import { startServerConnection } from "#lib/components/base/services/services.svelte.ts";

  let {
    id,
    name,
    onDeleted,
  }: {
    id: number;
    name: string;
    onDeleted: () => void;
  } = $props();
</script>

<div class="flex items-center justify-end gap-2">
  <Button variant="secondary" onclick={() => startServerConnection(id, name)}
    >Connect</Button
  >
  <DeleteRequestDialog {id} {name} {onDeleted} type="server" useIcon />
</div>
