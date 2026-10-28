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
    "flex items-start gap-3 rounded-md border border-destructive/25 bg-destructive/5 p-4",
    className,
  )}
  {...restProps}
>
  <span
    class="flex size-9 shrink-0 items-center justify-center rounded-md border border-destructive/20 bg-destructive/10 text-destructive"
    aria-hidden="true"
  >
    <CircleAlert class="size-5" />
  </span>
  <div class="min-w-0 flex-1 pt-0.5">
    <p class="text-sm font-semibold text-foreground">{title}</p>
    <p class="mt-1 wrap-break-word text-sm text-muted-foreground">{message}</p>
    {#if onRetry}
      <Button variant="outline" size="sm" class="mt-3" onclick={onRetry}>
        <RotateCcw />
        Try again
      </Button>
    {/if}
  </div>
</div>
