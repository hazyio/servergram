<script lang="ts">
  import "./layout.css";
  import favicon from "#lib/assets/favicon.svg";
  import type { LayoutProps } from "./$types";
  import { Toaster } from "#lib/components/ui/sonner/index.ts";
  import SideBarContainer from "#lib/components/base/sidebar/sideBarContainer.svelte";
  import appPaths from "#lib/appPaths.ts";
  import { page } from "$app/state";

  let { children }: LayoutProps = $props();
  const homeBreadcrumbs = [
    {
      title: "Home",
      url: appPaths.home,
    },
  ];
  const serverBreadcrumbs = [
    ...homeBreadcrumbs,
    {
      title: "Servers",
      url: appPaths.servers.index,
    },
  ];
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>
<Toaster />

{#if page.url.pathname === appPaths.account.login}
  {@render children()}
{:else}
  <SideBarContainer
    unTetheredBreadcrumbs={[
      [
        ...serverBreadcrumbs,
        {
          title: "Add New Server",
          url: appPaths.servers.add,
        },
      ],
    ]}
    sideBarItems={[
      {
        title: "Main",
        url: "#",
        items: [
          {
            title: "Home",
            url: appPaths.home,
            breadcrumb: [{ title: "Home", url: appPaths.home }],
          },
          {
            title: "Services",
            url: appPaths.home,
            breadcrumb: [{ title: "Home", url: appPaths.home }],
          },
          {
            title: "Databases",
            url: appPaths.home,
            breadcrumb: [{ title: "Home", url: appPaths.home }],
          },
          {
            title: "Cron Jobs",
            url: appPaths.home,
            breadcrumb: [{ title: "Home", url: appPaths.home }],
          },
        ],
      },
      {
        title: "Subs",
        url: "#",
        items: [
          {
            title: "Groups",
            url: appPaths.groups.index,
            breadcrumb: [
              { title: "Home", url: appPaths.home },
              { title: "Groups", url: appPaths.groups.index },
            ],
          },
          {
            title: "SSH Keys",
            url: appPaths.sshKeys.index,
            breadcrumb: [
              { title: "Home", url: appPaths.home },
              { title: "SSH Keys", url: appPaths.sshKeys.index },
            ],
          },
        ],
      },
    ]}
  >
    {@render children()}
  </SideBarContainer>
{/if}
