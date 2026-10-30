<script lang="ts">
  import * as Popover from "#lib/components/ui/popover/index.js";
  import * as Tabs from "#lib/components/ui/tabs/index.js";
  import { Button } from "../../ui/button";
  import Loading from "../loading.svelte";
  import ServiceItem, { type ServiceItemProps } from "./ServiceItem.svelte";

  let open = $state(true);
  let keepOpened = $state(true);
  let services: ServiceItemProps[] = $state([
    {
      id: 1,
      name: "The main boy",
      status: "connected",
    },
    {
      id: 2,
      name: "The Fat man",
      status: "pending",
    },
  ]);
</script>

<div
  class=" fixed bottom-2 right-2 rounded-full p-3 h-15 w-15 flex justify-center items-center"
>
  <Popover.Root
    {open}
    onOpenChange={(value) => {
      if (keepOpened) {
        open = true;
        return;
      }
      open = value;
    }}
  >
    <Popover.Trigger
      ><Button
        title="connected services"
        variant="outline"
        class="h-15 w-15 rounded-full"
      >
        3
      </Button></Popover.Trigger
    >

    <Popover.Content align="end" class="max-h-120 min-h-50">
      <Popover.Header>
        <h3 class="text-lg font-medium">Connected Services</h3>
      </Popover.Header>
      <Tabs.Root value="all" class="w-full">
        <Tabs.List class="w-full">
          <Tabs.Trigger value="all">All</Tabs.Trigger>
          <Tabs.Trigger value="connected">Connected</Tabs.Trigger>
          <Tabs.Trigger value="pending">Pending</Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content value="all">
          {#each services as service (service.id)}
            <ServiceItem
              id={service.id}
              name={service.name}
              status={service.status}
            />
          {/each}
        </Tabs.Content>
        <Tabs.Content value="pending">
     
          {#each services.filter((s) => s.status === "pending") as service (service.id)}
            <ServiceItem
              id={service.id}
              name={service.name}
              status={service.status}
            />
          {/each}
        </Tabs.Content>
        <Tabs.Content value="connected">
          {#each services.filter((s) => s.status === "connected") as service (service.id)}
            <ServiceItem
              id={service.id}
              name={service.name}
              status={service.status}
            />
          {/each}
        </Tabs.Content>
      </Tabs.Root>
    </Popover.Content>
  </Popover.Root>
</div>
