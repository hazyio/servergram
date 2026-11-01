<script lang="ts">
  import type { ComponentProps } from "svelte";
  import ArrowUpDownIcon from "@lucide/svelte/icons/arrow-up-down";
  import { Button } from "#lib/components/ui/button/index.js";

  let {
    variant = "ghost",
    text,
    children,
    ...restProps
  }: ComponentProps<typeof Button> & {
    text?: string;
  } = $props();
</script>

<Button {variant} {...restProps}>
  {#if children}
    {@render children()}
  {:else}
    {text}
  {/if}
  <ArrowUpDownIcon class="ms-2 size-4" />
</Button>
