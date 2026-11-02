<script lang="ts">
  import { type WSServiceConnectedListItem } from "#lib/types/websockets/index.ts";
  import { Database, Server } from "@lucide/svelte";
  import Loading from "../loading.svelte";
  import DisconnectService from "./DisconnectService.svelte";
  import appPaths from "#lib/appPaths.ts";

  let { id, name, status, type }: WSServiceConnectedListItem = $props();
</script>

<div
  class="mt-2 bg-muted p-2 rounded-md flex items-center gap-2 justify-center"
>
  {#if type === "server"}
    <Server class="h-4 w-4" />
    <a class="grow truncate" href={appPaths.servers.overview(id)}>{name}</a>
  {:else}
    <Database class="h-4 w-4" />
    <a class="grow truncate" href={appPaths.databases.add}>{name}</a>
  {/if}
  {#if status === "connected"}
    <DisconnectService {id} {name} {type} />
  {:else if status === "connecting"}
    <Loading loading={true} class="h-5 w-5 text-amber-500" />
  {/if}
</div>
