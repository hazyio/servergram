<script lang="ts">
  import Inbox from "@lucide/svelte/icons/inbox";
  import { cn } from "#lib/utils.js";
  import type { HTMLAttributes } from "svelte/elements";

  let {
    class: className = "",
    title = "No items found",
    message = "There is nothing to show here yet.",
    ...restProps
  }: HTMLAttributes<HTMLDivElement> & {
    title?: string;
    message?: string;
  } = $props();
</script>

<div
  role="status"
  class={cn(
    "flex min-h-48 flex-col items-center justify-center px-6 py-10 text-center",
    className,
  )}
  {...restProps}
>
  <span
    class="mb-4 flex size-12 items-center justify-center rounded-full border border-border bg-muted text-muted-foreground"
    aria-hidden="true"
  >
    <Inbox class="size-5" />
  </span>
  <h2 class="text-base font-semibold text-foreground">{title}</h2>
  <p class="mt-1 max-w-sm text-sm text-muted-foreground">{message}</p>
</div>
