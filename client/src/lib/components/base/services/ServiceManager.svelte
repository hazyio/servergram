<script lang="ts">
  import * as Popover from "#lib/components/ui/popover/index.js";
  import * as Tabs from "#lib/components/ui/tabs/index.js";
  import TrustHostFingerPrint from "./TrustHostFingerPrint.svelte";

  import { onMount } from "svelte";
  import { Button } from "../../ui/button/index.ts";
  import Error from "../error.svelte";
  import Loading from "../loading.svelte";
  import NoItemsFound from "../noItemsFound.svelte";
  import CollectPassPhrase from "./CollectPassPhrase.svelte";
  import ServiceItem from "./ServiceItem.svelte";
  import {
    checkWs,
    servicesState,
    startWebSocket,
    stopWebSocket,
    ConnectedServices,
    toastError,
  } from "./services.svelte.ts";
  import CompareTrustFingerPrint from "./CompareTrustFingerPrint.svelte";
  let open = $state(false);

  function startWebSocketWrapper() {
    startWebSocket(() => {
      if (servicesState.error || servicesState.services.length > 0) {
        open = true;
        return;
      }
    });
  }
  onMount(() => {
    startWebSocketWrapper();
    return stopWebSocket;
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
      onclick={() => {
        if (!checkWs()) {
          if (servicesState.tries >= 5) {
            toastError("Maximum connection attempts reached");
          } else {
            startWebSocketWrapper();
          }
          return;
        }
      }}
    >
      {#if servicesState.connecting}
        <Loading loading={true} />
      {:else}
        {servicesState.services.length}
      {/if}
    </Button></Popover.Trigger
  >

  <Popover.Content align="end" class="max-h-120 min-h-90">
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
