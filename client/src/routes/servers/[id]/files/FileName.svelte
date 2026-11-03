<script lang="ts">
  import { getFileIcon } from "#lib/fileIcons";

  let { name, isDir }: { name: string; isDir: boolean } = $props();
  const icon = $derived(getFileIcon(name, isDir));
</script>

<span class="flex items-center gap-2">
  <img src={icon} alt="" class="size-4 shrink-0" />
  {name}
</span>
