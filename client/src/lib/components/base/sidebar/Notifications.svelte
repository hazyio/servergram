<script lang="ts">
  import * as Popover from "#lib/components/ui/popover/index.js";
  import * as Tabs from "#lib/components/ui/tabs/index.js";
  import { Bell, BellDot } from "@lucide/svelte";
  import Error from "../ErrorSection.svelte";
  import Loading from "../loading.svelte";
  import NoItemsFound from "../noItemsFound.svelte";
  import NotificationItem from "./NotificationItem.svelte";
  import { Button } from "../../ui/button";
  import { onDestroy, onMount } from "svelte";
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import type { components } from "#lib/types/request-schema.js";

  type Notification = components["schemas"]["NotificationListItem"];

  let open = $state(false);
  let unreadNotificationsAvailable = $state(false);
  let gettingNotifications = $state(false);
  let notificationError = $state<string | null>(null);
  let statusFilter = $state<"unread" | "read">("unread");
  let notifications = $state<Notification[]>([]);
  let notificationRequest: AbortController | undefined;
  let availabilityRequest: AbortController | undefined;
  let availabilityTimer: ReturnType<typeof setTimeout> | undefined;
  let stopped = false;

  function cancelNotificationRequest() {
    notificationRequest?.abort();
    notificationRequest = undefined;
    gettingNotifications = false;
  }

  async function getNotifications() {
    notificationRequest?.abort();
    availabilityRequest?.abort();
    const controller = new AbortController();
    notificationRequest = controller;
    const requestedFilter = statusFilter;
    gettingNotifications = true;
    notificationError = null;

    const endpoint =
      requestedFilter === "read"
        ? "/api/notification/list/read"
        : "/api/notification/list/unread";

    try {
      const { data, error } = await NewRequestClient().GET(endpoint, {
        signal: controller.signal,
      });
      if (controller.signal.aborted) return;
      if (error) {
        notificationError = unWrapErrorMessage(
          error,
          "Could not load notifications.",
        );
        return;
      }

      notifications = data ?? [];
      if (requestedFilter === "unread") {
        unreadNotificationsAvailable = notifications.length > 0;
      }
    } catch (error) {
      if (!controller.signal.aborted) {
        notificationError =
          error instanceof globalThis.Error
            ? error.message
            : "Could not load notifications.";
      }
    } finally {
      if (notificationRequest === controller) {
        notificationRequest = undefined;
        gettingNotifications = false;
      }
    }
  }

  async function checkUnreadNotifications() {
    if (!gettingNotifications) {
      const controller = new AbortController();
      availabilityRequest = controller;
      try {
        const { data } = await NewRequestClient().GET(
          "/api/notification/available",
          { signal: controller.signal },
        );
        if (!controller.signal.aborted && data !== undefined) {
          unreadNotificationsAvailable = data.value;
        }
      } catch {
        // A later poll can retry if the availability check fails.
      } finally {
        if (availabilityRequest === controller) {
          availabilityRequest = undefined;
        }
      }
    }

    if (!stopped) {
      // TODO: Add timeout for polling
      // availabilityTimer = setTimeout(checkUnreadNotifications, 120_000);
    }
  }

  onMount(() => {
    void checkUnreadNotifications();
  });
  onDestroy(() => {
    stopped = true;
    if (availabilityTimer) clearTimeout(availabilityTimer);
    availabilityRequest?.abort();
    notificationRequest?.abort();
  });
</script>

<Popover.Root
  bind:open
  onOpenChange={(value) => {
    if (value) {
      statusFilter = "unread";
      void getNotifications();
    } else {
      cancelNotificationRequest();
    }
  }}
>
  <Popover.Trigger>
    <Button variant={"outline"} class="h-9 w-9 rounded-md">
      {#if unreadNotificationsAvailable}
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
    <Tabs.Root value={statusFilter} class="w-full">
      <Tabs.List class="w-full">
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
      <Tabs.Content value={statusFilter}>
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
            title={`No ${statusFilter} notifications`}
            message="Notifications will appear here when available."
          />
        {:else}
          {#each notifications as notification (notification.id)}
            <NotificationItem
              id={notification.id}
              title={notification.title}
              message={notification.message}
            />
          {/each}
        {/if}
      </Tabs.Content>
    </Tabs.Root>
  </Popover.Content>
</Popover.Root>
