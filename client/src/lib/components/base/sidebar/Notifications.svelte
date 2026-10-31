<script lang="ts">
  import * as Popover from "#lib/components/ui/popover/index.js";
  import * as Tabs from "#lib/components/ui/tabs/index.js";
  import { Bell, BellDot } from "@lucide/svelte";
  import Error from "../error.svelte";
  import Loading from "../loading.svelte";
  import NoItemsFound from "../noItemsFound.svelte";
  import NotificationItem from "./NotificationItem.svelte";
  import { Button } from "../../ui/button";
  import { onMount } from "svelte";
  import { NewRequestClient } from "#lib/api/request.ts";
  let newNotificationsNotificationAvailable = $state(false);
  let gettingNotifications = $state(false);
  let notificationError = $state(null);
  let statusFilter = $state<"new" | "unread" | "read">("new");
  let notifications = $state([
    {
      id: 1,
      title: "Welcome",
      message: "You have successfully connected a service.",
    },
  ]);
  function getNotifications() {
    gettingNotifications = true;
    newNotificationsNotificationAvailable = false;
    notificationError = null;
  }
  async function checkNewNotifications() {
    if (!gettingNotifications && !newNotificationsNotificationAvailable) {
      const { data } = await NewRequestClient().GET(
        "/api/notification/available",
      );
      if (data !== undefined) {
        newNotificationsNotificationAvailable = data.value;
      }
    }
    setTimeout(() => {
      checkNewNotifications();
    }, 50000);
  }
  onMount(() => {
    checkNewNotifications();
  });
</script>

<Popover.Root>
  <Popover.Trigger>
    <Button variant={"outline"} class="h-9 w-9 rounded-md">
      {#if newNotificationsNotificationAvailable}
        <BellDot />
      {:else}
        <Bell />
      {/if}
    </Button>
  </Popover.Trigger>

  <Popover.Content align="end" class="min-h-[90vh]">
    <Popover.Header>
      <h3 class="text-lg font-medium">Notifications</h3>
    </Popover.Header>
    {#if gettingNotifications}
      <Loading
        loading={true}
        text={`Getting ${statusFilter} notifications`}
        class="mt-30 h-12"
      />
    {:else if notificationError !== null}
      <Error message={notificationError} />
    {:else if notifications.length === 0}
      <NoItemsFound
        title="No unread notifications"
        message="Connect a service to get started."
      />
    {:else}
      <Tabs.Root value={statusFilter} class="w-full">
        <Tabs.List class="w-full">
          <Tabs.Trigger
            value="new"
            onclick={() => {
              statusFilter = "new";
              getNotifications();
            }}>New</Tabs.Trigger
          >
          <Tabs.Trigger
            value="unread"
            onclick={() => {
              statusFilter = "unread";
              getNotifications();
            }}>Unread</Tabs.Trigger
          >
          <Tabs.Trigger
            value="read"
            onclick={() => {
              statusFilter = "read";
              getNotifications();
            }}>Read</Tabs.Trigger
          >
        </Tabs.List>
        <Tabs.Content value="new">
          {#each notifications as notification (notification.id)}
            <NotificationItem
              id={notification.id}
              title={notification.title}
              message={notification.message}
            />
          {/each}
        </Tabs.Content>
        <Tabs.Content value="read">
          {#each notifications as notification (notification.id)}
            <NotificationItem
              id={notification.id}
              title={notification.title}
              message={notification.message}
            />
          {/each}
        </Tabs.Content>
        <Tabs.Content value="unread">
          {#each notifications as notification (notification.id)}
            <NotificationItem
              id={notification.id}
              title={notification.title}
              message={notification.message}
            />
          {/each}
        </Tabs.Content>
      </Tabs.Root>
    {/if}
  </Popover.Content>
</Popover.Root>
