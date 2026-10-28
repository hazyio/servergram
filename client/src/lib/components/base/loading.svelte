<script lang="ts">
  import Loader from "@lucide/svelte/icons/loader";
  import { cn } from "cn";
  import type { HTMLAttributes } from "svelte/elements";

  let {
    class: className = "",
    text = "",
    loading,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> & {
    loading: boolean;
    text?: string;
  } = $props();
</script>

{#if loading}
  <div class={cn("flex flex-col justify-center items-center", className)} {...restProps}>
    <Loader class="h-1/10 w-1/10 animate-spin" />
    {#if text}
      <span class="ml-2">{text}</span>
    {/if}
  </div>
{/if}
