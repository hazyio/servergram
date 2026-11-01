<script lang="ts">
  import appPaths from "#lib/appPaths.ts";
  import {
    ConnectedServices
  } from "#lib/components/base/services/services.svelte.ts";
  let { id, name } = $props();
  let server = $derived(ConnectedServices.findServer(id));
</script>

{#if server}
  <a href={appPaths.servers.show(id)}>
    {name}
  </a>
{:else}
  <span>
    {name}
  </span>
{/if}
