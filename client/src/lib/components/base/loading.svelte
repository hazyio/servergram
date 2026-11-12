<script lang="ts">
  import Loader from "@lucide/svelte/icons/loader";
  import { cn } from "cn";
  import type { HTMLAttributes } from "svelte/elements";

  let {
    containerClass = "",
    class: className = "",
    text = "",
    loading,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> & {
    containerClass?: string;
    loading: boolean;
    text?: string;
  } = $props();
</script>

{#if loading}
  <div
    class={cn("flex flex-col justify-center items-center", containerClass)}
    {...restProps}
  >
    <Loader class={cn("h-1/10 w-1/10 animate-spin", className)} />
    {#if text}
      <span class="mt-5">{text}</span>
    {/if}
  </div>
{/if}
