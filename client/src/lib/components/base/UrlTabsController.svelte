<script lang="ts">
  import { page } from "$app/state";
  import { cn } from "#lib/utils.js";

  let {
    items,
  }: {
    items: { label: string; href: string }[];
  } = $props();

  let currentPath = $derived(page.url.pathname.replace(/\/+$/, ""));
  const isActive = (href: string) => {
    const path = href.replace(/\/+$/, "");
    return currentPath === path || currentPath.startsWith(`${path}/`);
  };
</script>

<nav
  aria-label="Section navigation"
  class="mx-2 mt-2 flex w-fit items-center gap-1 rounded-lg border border-border/50 bg-muted/50 p-1"
>
  {#each items as item (item.href)}
    {@const active = isActive(item.href)}
    <a
      href={item.href}
      aria-current={active ? "page" : undefined}
      class={cn(
        "rounded-md px-3 py-1.5 text-sm font-medium whitespace-nowrap text-muted-foreground transition-colors duration-150",
        "hover:bg-background/70 hover:text-foreground",
        "focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:outline-none",
        active && "bg-background text-foreground shadow-sm",
      )}
    >
      {item.label}
    </a>
  {/each}
</nav>
