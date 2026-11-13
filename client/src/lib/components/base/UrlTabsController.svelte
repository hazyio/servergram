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
  class="flex items-center gap-1 overflow-x-none border-b border-border/60 px-2"
>
  {#each items as item (item.href)}
    {@const active = isActive(item.href)}
    <a
      href={item.href}
      aria-current={active ? "page" : undefined}
      class={cn(
        "-mb-px border-b-2 border-transparent px-3 py-2.5 text-sm font-medium whitespace-nowrap text-muted-foreground transition-colors duration-150",
        "hover:border-border hover:text-foreground",
        "focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:outline-none",
        active && "border-primary text-foreground",
      )}
    >
      {item.label}
    </a>
  {/each}
</nav>
