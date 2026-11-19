<script lang="ts">
  import * as Breadcrumb from "#lib/components/ui/breadcrumb/index.js";
  import { Separator } from "#lib/components/ui/separator/index.js";
  import * as Sidebar from "#lib/components/ui/sidebar/index.js";
  import { page } from "$app/state";
  import type { Snippet } from "svelte";
  import ConnectedServices from "../services/ServiceManager.svelte";
  import type { SideBarItem } from "./AppSideBar.svelte";
  import AppSidebar from "./AppSideBar.svelte";
  import Notifications from "./Notifications.svelte";

  let {
    children,
    sideBarItems,
    unTetheredBreadcrumbs,
  }: {
    children?: Snippet;
    sideBarItems: SideBarItem[];
    unTetheredBreadcrumbs: SideBarItem[][];
  } = $props();
  let combinedBreadcrumbs = $derived([
    ...unTetheredBreadcrumbs.flat(),
    ...sideBarItems.flatMap((group) => group.items ?? []),
  ]);
  let activeItem = $derived(
    combinedBreadcrumbs.find(
      (item) =>
        new URL(item.url, page.url.origin).pathname === page.url.pathname,
    ),
  );
  let breadcrumbItems = $derived.by(() => {
    if (!activeItem?.url) return [];
    if (activeItem.breadcrumb) return activeItem.breadcrumb;

    const pathBreadcrumbs = unTetheredBreadcrumbs.find((breadcrumbs) =>
      breadcrumbs.some(
        (item) =>
          new URL(item.url, page.url.origin).pathname === page.url.pathname,
      ),
    );
    if (!pathBreadcrumbs) {
      return [{ title: activeItem.title, url: activeItem.url }];
    }

    const activeIndex = pathBreadcrumbs.findIndex(
      (item) =>
        new URL(item.url, page.url.origin).pathname === page.url.pathname,
    );
    return pathBreadcrumbs.slice(0, activeIndex + 1);
  });
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
      <div class="ms-auto gap-3 flex items-center">
        <ConnectedServices />


        <Notifications />
      </div>
    </header>
    {@render children?.()}
  </Sidebar.Inset>
</Sidebar.Provider>
