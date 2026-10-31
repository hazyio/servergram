<script lang="ts">
  import { type WSServiceConnectedListItem } from "#lib/types/websockets/index.ts";
  import Loading from "../loading.svelte";

  let { id, name, status }: WSServiceConnectedListItem = $props();
</script>

<div class="mt-2 bg-muted p-2 rounded-md flex gap-2 justify-center">
  <span class="grow"> {name} </span>
  {#if status === "connected"}
    <span class="text-green-500">Connected</span>
  {:else if status === "connecting"}
    <Loading loading={true} class="h-5 w-5 text-amber-500" />
  {/if}
</div>
