<script lang="ts">
  import appPaths from "#lib/appPaths.ts";
  import { Button } from "#lib/components/ui/button/index.ts";
  import { PlusIcon } from "@lucide/svelte";

  let { serverId }: { serverId: number } = $props();
</script>

<Button href={appPaths.servers.cronAdd(serverId)}>
  <PlusIcon />
  Add Job
</Button>
