<script lang="ts">
  import TimeAgo from "#lib/components/base/TimeAgo.svelte";

  let {
    datetime,
    empty = "-",
  }: {
    datetime: string | null;
    empty?: string;
  } = $props();
</script>

{#if datetime}
  <TimeAgo {datetime} />
{:else}
  <span class="text-muted-foreground">{empty}</span>
{/if}
