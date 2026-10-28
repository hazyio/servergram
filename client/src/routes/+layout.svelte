<script lang="ts">
  import "./layout.css";
  import favicon from "#lib/assets/favicon.svg";
  import type { LayoutProps } from "./$types";
  import { Toaster } from "#lib/components/ui/sonner/index.ts";
  import SideBarContainer from "#lib/components/base/sidebar/sideBarContainer.svelte";
  import appPaths from "#lib/appPaths.ts";

  let { children }: LayoutProps = $props();
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>
<Toaster />

<SideBarContainer
  sideBarItems={[
    {
      title: "Getting Started",
      url: "#",
      items: [
        {
          title: "Home",
          url: appPaths.home,
          breadcrumb: [{ title: "Home", url: appPaths.home }],
        },
        {
          title: "Groups",
          url: appPaths.groups,
          breadcrumb: [
            { title: "Home", url: appPaths.home },
            { title: "Groups", url: appPaths.groups },
          ],
        },
        {
          title: "SSH Keys",
          url: appPaths.sshKeys,
          breadcrumb: [
            { title: "Home", url: appPaths.home },
            { title: "SSH Keys", url: appPaths.sshKeys },
          ],
        },
      ],
    },
  ]}
>
  {@render children()}
</SideBarContainer>
