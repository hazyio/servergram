<script lang="ts">
  import appPaths from "#lib/appPaths.ts";

  let { id, count } = $props();
</script>

<a href={appPaths.sshKeys.resources(id)}>
  {count}
</a>
