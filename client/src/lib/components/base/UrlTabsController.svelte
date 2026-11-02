<script lang="ts">
  import { page } from "$app/state";

  let {
    items,
  }: {
    items: { label: string; href: string }[];
  } = $props();
  let url = page.url.pathname;
</script>

<div class="p-1 bg-muted/60 flex justify-between px-3">
  {#each items as item}
    <a
      class={url === item.href
        ? "bg-background/80 p-1 px-3 rounded-md"
        : "p-1 px-3"}
      href={item.href}>{item.label}</a
    >
  {/each}
</div>
