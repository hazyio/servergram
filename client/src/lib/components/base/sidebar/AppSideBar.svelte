<script lang="ts">
  import { page } from "$app/state";
  import * as Sidebar from "#lib/components/ui/sidebar/index.js";
  import type { ComponentProps } from "svelte";
  import ThemeChanger from "./ThemeChanger.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  export interface SideBarItem {
    title: string;
    url: string;
    breadcrumb?: {
      title: string;
      url: string;
    }[];
    items?: {
      title: string;
      url: string;
      breadcrumb?: {
        title: string;
        url: string;
      }[];
    }[];
  }
  let {
    ref = $bindable(null),
    items = [],
    ...restProps
  }: ComponentProps<typeof Sidebar.Root> & { items: SideBarItem[] } = $props();
</script>

<Sidebar.Root {...restProps} bind:ref>
  <Sidebar.Header class="flex flex-row items-center justify-between">
    <span>Server Gram</span>
    <ThemeChanger />
  </Sidebar.Header>
  <Sidebar.Content>
    <!-- We create a Sidebar.Group for each parent. -->
    {#each items as group (group.title)}
      <Sidebar.Group>
        <Sidebar.GroupLabel>{group.title}</Sidebar.GroupLabel>
        <Sidebar.GroupContent>
          <Sidebar.Menu>
            {#each group.items as item (item.title)}
              <Sidebar.MenuItem class="mt-0.5">
                <Sidebar.MenuButton
                  isActive={new URL(item.url, page.url.origin).pathname ===
                    page.url.pathname}
                >
                  {#snippet child({ props })}
                    <a href={item.url} {...props}>{item.title}</a>
                  {/snippet}
                </Sidebar.MenuButton>
              </Sidebar.MenuItem>
            {/each}
          </Sidebar.Menu>
        </Sidebar.GroupContent>
      </Sidebar.Group>
    {/each}
  </Sidebar.Content>
  <Sidebar.Rail />
  <Sidebar.Footer>
    <Button variant="destructive" class="grow">Sign Out</Button>
  </Sidebar.Footer>
</Sidebar.Root>
