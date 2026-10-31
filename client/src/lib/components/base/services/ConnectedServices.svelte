<script lang="ts">
  import * as Popover from "#lib/components/ui/popover/index.js";
  import * as Tabs from "#lib/components/ui/tabs/index.js";
  import {
    ConnectToSSH,
    PassPhraseRequired,
    Error as ServerWSError,
    type PassPhraseRequiredWSData,
    type ServiceWSMessage,
    type WSError,
  } from "#lib/types/websockets/index.ts";
  import { toast } from "svelte-sonner";

  import { onMount } from "svelte";
  import { Button } from "../../ui/button/index.ts";
  import Error from "../error.svelte";
  import Loading from "../loading.svelte";
  import NoItemsFound from "../noItemsFound.svelte";
  import CollectPassPhrase from "./CollectPassPhrase.svelte";
  import ServiceItem, { type ServiceItemProps } from "./ServiceItem.svelte";
  import {
    servicesState,
    startWebSocket,
    stopWebSocket,
    checkWs,
  } from "./services.svelte.ts";

  let open = $state(false);
  let collectPassPhraseFor = $state<PassPhraseRequiredWSData[]>([]);
  let services: ServiceItemProps[] = $state([]);
  function handleWebSocketMessage(message: ServiceWSMessage) {
    switch (message.action) {
      case ConnectToSSH:
        break;
      case PassPhraseRequired:
        let findIndex = collectPassPhraseFor.findIndex(
          (item) => item.id === (message.data as PassPhraseRequiredWSData).id,
        );
        if (findIndex === -1) {
          // Add the service to the list of services requiring a pass phrase
          collectPassPhraseFor.push(message.data as PassPhraseRequiredWSData);
        }
        break;
      case ServerWSError:
        // Handle server error message
        toast.error((message.data as WSError).message);
        break;
      default:
        console.warn("Unhandled WebSocket message:", message);
        break;
    }
  }
  function startWebSocketWrapper() {
    startWebSocket(() => {
      if (servicesState.error || servicesState.services.length > 0) {
        open = true;
        return;
      }
    }, handleWebSocketMessage);
  }
  onMount(() => {
    startWebSocketWrapper();
    return stopWebSocket;
  });
</script>

<CollectPassPhrase bind:data={collectPassPhraseFor} />
<div
  class=" fixed bottom-2 right-2 rounded-full p-3 h-15 w-15 flex justify-center items-center"
>
  <Popover.Root
    {open}
    onOpenChange={(value) => {
      open = value;
    }}
  >
    <Popover.Trigger
      ><Button
        title={servicesState.error ? servicesState.error : "connected services"}
        variant={servicesState.error ? "destructive" : "outline"}
        class="h-15 w-15 rounded-full"
        onclick={() => {
          checkWs();
        }}
      >
        {#if servicesState.connecting}
          <Loading loading={true} />
        {:else}
          {services.length}
        {/if}
      </Button></Popover.Trigger
    >

    <Popover.Content align="end" class="max-h-120 min-h-50">
      <Popover.Header>
        <h3 class="text-lg font-medium">Connected Services</h3>
      </Popover.Header>
      {#if servicesState.connecting}
        <Loading loading={true} text="Connecting..." class="mt-5" />
      {:else if servicesState.error !== null}
        <Error
          message={servicesState.error}
          onRetry={() => startWebSocketWrapper()}
        />
      {:else if servicesState.services.length === 0}
        <NoItemsFound
          title="No connected services"
          message="Connect a service to get started."
        />
      {:else}
        <Tabs.Root value="all" class="w-full">
          <Tabs.List class="w-full">
            <Tabs.Trigger value="all">All</Tabs.Trigger>
            <Tabs.Trigger value="connected">Connected</Tabs.Trigger>
            <Tabs.Trigger value="pending">Pending</Tabs.Trigger>
          </Tabs.List>
          <Tabs.Content value="all">
            {#each servicesState.services as service (service.id)}
              <ServiceItem
                id={service.id}
                name={service.name}
                status={service.status}
              />
            {/each}
          </Tabs.Content>
          <Tabs.Content value="pending">
            {#each servicesState.services.filter((s) => s.status === "pending") as service (service.id)}
              <ServiceItem
                id={service.id}
                name={service.name}
                status={service.status}
              />
            {/each}
          </Tabs.Content>
          <Tabs.Content value="connected">
            {#each servicesState.services.filter((s) => s.status === "connected") as service (service.id)}
              <ServiceItem
                id={service.id}
                name={service.name}
                status={service.status}
              />
            {/each}
          </Tabs.Content>
        </Tabs.Root>
      {/if}
    </Popover.Content>
  </Popover.Root>
</div>
