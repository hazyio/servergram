<script lang="ts" module>
  export interface ServiceItemProps {
    id: number;
    name: string;
    status: "connected" | "pending";
  }
</script>

<script lang="ts">
  import Loading from "../loading.svelte";

  let { id, name, status }: ServiceItemProps = $props();
</script>

<div class="mt-2 bg-muted p-2 rounded-md flex gap-2 justify-center">
  <span class="grow"> {name} </span>
  {#if status === "connected"}
    <span class="text-green-500">Connected</span>
  {:else if status === "pending"}
    <Loading loading={true} class="h-5 w-5 text-amber-500" />
  {/if}
</div>
