<script lang="ts">
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";
  import { isCronDateSet } from "../cronJobs.ts";

  let {
    datetime,
    empty = "-",
  }: {
    datetime: string | null;
    empty?: string;
  } = $props();
</script>

{#if isCronDateSet(datetime)}
  <TimeAgo {datetime} />
{:else}
  <span class="text-muted-foreground">{empty}</span>
{/if}
