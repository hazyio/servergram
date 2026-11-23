<script lang="ts">
  import type { PageProps } from "./$types";

  let { params }: PageProps = $props();
</script>

<svelte:head>
  <title>Run #{params.runId} log | ServerGram</title>
</svelte:head>
