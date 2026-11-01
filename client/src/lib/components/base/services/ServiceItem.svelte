<script lang="ts">
  import { type WSServiceConnectedListItem } from "#lib/types/websockets/index.ts";
  import { Database, Server } from "@lucide/svelte";
  import Loading from "../loading.svelte";
  import DisconnectService from "./DisconnectService.svelte";

  let { id, name, status, type }: WSServiceConnectedListItem = $props();
</script>

<div
  class="mt-2 bg-muted p-2 rounded-md flex items-center gap-2 justify-center"
>
  {#if type === "server"}
    <Server class="h-4 w-4" />
  {:else}
    <Database class="h-4 w-4" />
  {/if}
  <span class="grow truncate" title={name}> {name} </span>
  {#if status === "connected"}
    <DisconnectService {id} {name} {type} />
  {:else if status === "connecting"}
    <Loading loading={true} class="h-5 w-5 text-amber-500" />
  {/if}
</div>
