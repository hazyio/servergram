<script lang="ts">
  import CircleAlert from "@lucide/svelte/icons/circle-alert";
  import RotateCcw from "@lucide/svelte/icons/rotate-ccw";
  import Button from "#lib/components/ui/button/button.svelte";
  import { cn } from "#lib/utils.js";
  import type { HTMLAttributes } from "svelte/elements";

  let {
    class: className = "",
    title = "Something went wrong",
    message = "An unexpected error occurred. Please try again.",
    onRetry,
    ...restProps
  }: HTMLAttributes<HTMLDivElement> & {
    title?: string;
    message?: string;
    onRetry?: () => void;
  } = $props();
</script>

<div
  role="alert"
  aria-live="assertive"
  class={cn(
    "relative flex items-start gap-3 overflow-hidden border-y border-border bg-[linear-gradient(90deg,color-mix(in_oklch,var(--primary)_4%,transparent)_1px,transparent_1px)] bg-size-[3rem_100%] p-4",
    className,
  )}
  {...restProps}
>
  <span
    class="grid size-9 shrink-0 place-items-center border border-chart-3/35 text-chart-3"
    aria-hidden="true"
  >
    <CircleAlert class="size-5" />
  </span>
  <div class="min-w-0 flex-1 pt-0.5">
    <p class="font-mono text-sm font-semibold text-foreground">{title}</p>
    <p class="mt-1 wrap-break-word text-sm text-muted-foreground">{message}</p>
    {#if onRetry}
      <Button variant="outline" size="sm" class="mt-3" onclick={onRetry}>
        <RotateCcw />
        Try again
      </Button>
    {/if}
  </div>
  <span
    class="absolute right-0 bottom-0 h-0.5 w-12 bg-primary"
    aria-hidden="true"
  ></span>
</div>
