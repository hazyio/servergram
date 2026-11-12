<script lang="ts">
  import appPaths from "#lib/appPaths.ts";
  import favicon from "#lib/assets/favicon.svg";
  import SideBarContainer from "#lib/components/base/sidebar/SideBarContainer.svelte";
  import { Toaster } from "#lib/components/ui/sonner/index.ts";
  import { page } from "$app/state";
  import type { LayoutProps } from "./$types";
  import "./layout.css";
  import { severPageInfo } from "#lib/states/severPageInfo.svelte.ts";

  let { children }: LayoutProps = $props();
  let connectedServerInfo = $derived(severPageInfo);
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
  const severPageBreadcrumbs = $derived([
    ...serverBreadcrumbs,
    {
      title: connectedServerInfo.name,
      url: appPaths.servers.overview(connectedServerInfo.id),
    },
  ]);
  const severOverviewBreadcrumbs = $derived([
    ...serverBreadcrumbs,
    {
      title: connectedServerInfo.name,
      url: appPaths.servers.overview(connectedServerInfo.id),
    },
    {
      title: "Overview",
      url: appPaths.servers.overview(connectedServerInfo.id),
    },
  ]);
  const severFilesBreadcrumbs = $derived([
    ...serverBreadcrumbs,
    {
      title: connectedServerInfo.name,
      url: appPaths.servers.overview(connectedServerInfo.id),
    },
    {
      title: "Files",
      url: appPaths.servers.files(connectedServerInfo.id),
    },
  ]);
  const severProcessesBreadcrumbs = $derived([
    ...serverBreadcrumbs,
    {
      title: connectedServerInfo.name,
      url: appPaths.servers.overview(connectedServerInfo.id),
    },
    {
      title: "Processes",
      url: appPaths.servers.processes(connectedServerInfo.id),
    },
  ]);
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
      severOverviewBreadcrumbs,
      severPageBreadcrumbs,
      severFilesBreadcrumbs,
      severProcessesBreadcrumbs,
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
            title: "Servers",
            url: appPaths.servers.index,
            breadcrumb: [
              ...homeBreadcrumbs,
              { title: "Servers", url: appPaths.servers.index },
            ],
          },
          {
            title: "Databases",
            url: appPaths.databases.index,
            breadcrumb: [
              ...homeBreadcrumbs,
              { title: "Databases", url: appPaths.databases.index },
            ],
          },
          {
            title: "Cron Jobs",
            url: appPaths.cronJobs.index,
            breadcrumb: [
              ...homeBreadcrumbs,
              { title: "Cron Jobs", url: appPaths.cronJobs.index },
            ],
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
          {
            title: "Logs",
            url: appPaths.logs.index,
            breadcrumb: [
              ...homeBreadcrumbs,
              { title: "Logs", url: appPaths.logs.index },
            ],
          },
        ],
      },
    ]}
  >
    {@render children()}
  </SideBarContainer>
{/if}
