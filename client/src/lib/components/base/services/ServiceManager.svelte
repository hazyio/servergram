<script lang="ts">
  import * as Popover from "#lib/components/ui/popover/index.js";
  import * as Tabs from "#lib/components/ui/tabs/index.js";
  import TrustHostFingerPrint from "./TrustHostFingerPrint.svelte";

  import { onMount } from "svelte";
  import { Button } from "../../ui/button/index.ts";
  import Error from "../ErrorSection.svelte";
  import Loading from "../loading.svelte";
  import NoItemsFound from "../noItemsFound.svelte";
  import CollectPassPhrase from "./CollectPassPhrase.svelte";
  import CompareTrustFingerPrint from "./CompareTrustFingerPrint.svelte";
  import ServiceItem from "./ServiceItem.svelte";

  import { RefreshCcw } from "@lucide/svelte";
  import {
    checkWs,
    servicesState,
    startWebSocket,
    stopWebSocket,
    toastError,
  } from "./services.svelte.ts";
  import { GetConnectedServices } from "./webSocketMessages";
  let open = $state(false);
  let intervalId = $state<ReturnType<typeof setInterval> | null>(null);
  function retryWebSocketConnection() {
    if (servicesState.tries >= 5) {
      toastError("Maximum connection attempts reached");
    } else {
      startWebSocketWrapper();
    }
  }
  function startWebSocketWrapper() {
    if (intervalId !== null) {
      clearInterval(intervalId);
    }
    intervalId = setInterval(() => {
      if (!checkWs()) {
        retryWebSocketConnection();
        return;
      }
    }, 10000);
    startWebSocket(() => {
      if (servicesState.error) {
        open = true;
        return;
      }
    });
  }
  onMount(() => {
    startWebSocketWrapper();
    return () => {
      if (intervalId !== null) {
        clearInterval(intervalId);
      }
      stopWebSocket();
    };
  });
</script>

<CollectPassPhrase bind:data={servicesState.pendingPassPhraseCollection} />
<TrustHostFingerPrint bind:data={servicesState.pendingHostTrustCollection} />
<CompareTrustFingerPrint
  bind:data={servicesState.pendingHostTrustMismatchCollection}
/>
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
      class="h-9 w-9 rounded-md"
      onclick={retryWebSocketConnection}
    >
      {#if servicesState.connecting || servicesState.gettingConnected}
        <Loading loading={true} />
      {:else}
        {servicesState.services.length}
      {/if}
    </Button></Popover.Trigger
  >

  <Popover.Content align="end" class="max-h-120 min-h-90">
    <Popover.Header>
      <h3 class="text-lg font-medium flex">
        <span class="grow">Connected Services</span>
        {#if !servicesState.connecting && servicesState.error === null}
          <Button
            size="sm"
            variant="outline"
            disabled={servicesState.gettingConnected}
            onclick={GetConnectedServices}
            title="Reload connected services"
            ><RefreshCcw
              class={servicesState.gettingConnected ? "animate-spin" : ""}
            /></Button
          >
        {/if}
      </h3>
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
          <Tabs.Trigger value="connecting">Connecting</Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content value="all">
          {#each servicesState.services as service (service.id)}
            <ServiceItem {...service} />
          {/each}
        </Tabs.Content>
        <Tabs.Content value="connecting">
          {#each servicesState.services.filter((s) => s.status === "connecting") as service (service.id)}
            <ServiceItem {...service} />
          {/each}
        </Tabs.Content>
        <Tabs.Content value="connected">
          {#each servicesState.services.filter((s) => s.status === "connected") as service (service.id)}
            <ServiceItem {...service} />
          {/each}
        </Tabs.Content>
      </Tabs.Root>
    {/if}
  </Popover.Content>
</Popover.Root>
