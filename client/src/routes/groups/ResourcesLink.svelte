<script lang="ts">
  import appPaths from "#lib/appPaths.ts";

  let { id, count }: { id: number; count: number } = $props();
</script>

<a href={appPaths.groups.resources(id)}>
  {count}
</a>
