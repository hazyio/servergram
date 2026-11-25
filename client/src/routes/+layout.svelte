<script lang="ts">
  import appPaths from "#lib/appPaths.ts";
  import favicon from "#lib/assets/favicon.svg";
  import SideBarContainer from "#lib/components/base/sidebar/SideBarContainer.svelte";
  import OfflineNotice from "#lib/components/base/pwa/OfflineNotice.svelte";
  import { Toaster } from "#lib/components/ui/sonner/index.ts";
  import { page } from "$app/state";
  import type { LayoutProps } from "./$types";
  import "./layout.css";
  import {
    severPageInfo,
    cronJobPageInfo,
  } from "#lib/states/severPageInfo.svelte.ts";

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
  const severCronBreadcrumbs = $derived([
    ...serverBreadcrumbs,
    {
      title: connectedServerInfo.name,
      url: appPaths.servers.overview(connectedServerInfo.id),
    },
    {
      title: "Cron Jobs",
      url: appPaths.servers.cron(connectedServerInfo.id),
    },
  ]);
  const severAddCronBreadcrumbs = $derived([
    ...serverBreadcrumbs,
    {
      title: connectedServerInfo.name,
      url: appPaths.servers.overview(connectedServerInfo.id),
    },
    {
      title: "Cron Jobs",
      url: appPaths.servers.cron(connectedServerInfo.id),
    },
    {
      title: "Add Cron Job",
      url: appPaths.servers.cronAdd(connectedServerInfo.id),
    },
  ]);
  const isCronJobPage = $derived(
    page.route.id === "/servers/[id]/cron/[jobId]",
  );
  const severCronJobBreadcrumbs = $derived([
    ...serverBreadcrumbs,
    {
      title: connectedServerInfo.name,
      url: appPaths.servers.overview(connectedServerInfo.id),
    },
    {
      title: "Cron Jobs",
      url: appPaths.servers.cron(connectedServerInfo.id),
    },
    {
      title: cronJobPageInfo.name || "Job runs",
      url: isCronJobPage
        ? page.url.pathname
        : appPaths.servers.cronJob(connectedServerInfo.id, cronJobPageInfo.id),
    },
  ]);
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>
<Toaster />
<OfflineNotice />

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
      severAddCronBreadcrumbs,
      severCronJobBreadcrumbs,
      severCronBreadcrumbs,
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
