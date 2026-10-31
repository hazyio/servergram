<script lang="ts">
  import DeleteRequestDialog from "#lib/components/base/DeleteRequestDialog.svelte";
  import ViewKeyInfo from "./viewKeyInfo.svelte";

  let {
    id,
    onDeleted,
    name,
  }: {
    id: number;
    onDeleted: () => void;
    name: string;
  } = $props();
</script>

<div class="flex items-center justify-end gap-2">


  <ViewKeyInfo {id} />
  <DeleteRequestDialog {id} {name} type="ssh-key" useIcon={true} onDeleted={onDeleted} /> 
</div>
