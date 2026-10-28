<script lang="ts">
  import { page } from "$app/state";
  import * as Breadcrumb from "#lib/components/ui/breadcrumb/index.js";
  import * as Sidebar from "#lib/components/ui/sidebar/index.js";
  import { Separator } from "#lib/components/ui/separator/index.js";
  import AppSidebar from "./appSideBar.svelte";
  import type { Snippet } from "svelte";
  import type { SideBarItem } from "./appSideBar.svelte";

  let {
    children,
    sideBarItems,
  }: { children?: Snippet; sideBarItems: SideBarItem[] } = $props();
  let activeItem = $derived(
    sideBarItems
      .flatMap((group) => group.items ?? [])
      .find(
        (item) =>
          new URL(item.url, page.url.origin).pathname === page.url.pathname,
      ),
  );
  let breadcrumbItems = $derived(
    activeItem?.breadcrumb ??
      (activeItem ? [{ title: activeItem.title, url: activeItem.url }] : []),
  );
</script>

<Sidebar.Provider>
  <AppSidebar items={sideBarItems} />
  <Sidebar.Inset>
    <header class="flex h-16 shrink-0 items-center gap-2 border-b px-4">
      <Sidebar.Trigger class="-ms-1" />
      <Separator
        orientation="vertical"
        class="me-2 data-vertical:h-4 data-vertical:self-auto"
      />
      {#if breadcrumbItems.length > 0}
        <Breadcrumb.Root>
          <Breadcrumb.List>
            {#each breadcrumbItems as item, index (item.title)}
              <Breadcrumb.Item
                class={index < breadcrumbItems.length - 1
                  ? "hidden md:block"
                  : ""}
              >
                {#if index < breadcrumbItems.length - 1}
                  <Breadcrumb.Link href={item.url}>{item.title}</Breadcrumb.Link
                  >
                {:else}
                  <Breadcrumb.Page>{item.title}</Breadcrumb.Page>
                {/if}
              </Breadcrumb.Item>
              {#if index < breadcrumbItems.length - 1}
                <Breadcrumb.Separator class="hidden md:block" />
              {/if}
            {/each}
          </Breadcrumb.List>
        </Breadcrumb.Root>
      {/if}
    </header>
    {@render children?.()}
  </Sidebar.Inset>
</Sidebar.Provider>
