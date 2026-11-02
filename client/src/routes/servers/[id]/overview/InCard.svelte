<script lang="ts">
  import type { Snippet } from "svelte";
  import { cn } from "#lib/utils.ts";

  let {
    children,
    class: className,
  }: {
    children: Snippet;
    class?: string;
  } = $props();
</script>

<div class={cn("bg-muted rounded-md p-2 ", className)}>
  {@render children()}
</div>
