<script lang="ts">
  import * as ContextMenu from "#lib/components/ui/context-menu/index.js";
  import { CircleStop, Send, Skull } from "@lucide/svelte";
  import ProcessActionDialog from "./ProcessActionDialog.svelte";
  import {
    forceStopProcessAction,
    processSignalActions,
    stopProcessAction,
    type ProcessAction,
  } from "./processActions";

  let {
    name,
    pid,
    serverId,
  }: {
    name: string;
    pid: number;
    serverId: number;
  } = $props();

  let open = $state(false);
  let action = $state<ProcessAction | null>(null);

  function run(nextAction: ProcessAction) {
    action = nextAction;
    open = true;
  }
</script>

<ContextMenu.Root>
  <ContextMenu.Trigger>
    <div class="flex max-w-50 items-center gap-3" title={name}>
      <span class="truncate">{name}</span>
    </div>
  </ContextMenu.Trigger>
  <ContextMenu.Content>
    <ContextMenu.Group>
      <ContextMenu.Item
        class="py-2 my-1"
        onSelect={() => run(stopProcessAction)}
      >
        <CircleStop /> Stop
      </ContextMenu.Item>
      <ContextMenu.Item
        class="py-2 my-1"
        variant="destructive"
        onSelect={() => run(forceStopProcessAction)}
      >
        <Skull /> Force stop
      </ContextMenu.Item>
      <ContextMenu.Sub>
        <ContextMenu.SubTrigger class="py-2 my-1">
          <Send /> Send signal
        </ContextMenu.SubTrigger>
        <ContextMenu.SubContent>
          {#each processSignalActions as signalAction (signalAction.signal)}
            <ContextMenu.Item
              class="py-2 my-1"
              variant={signalAction.variant}
              onSelect={() => run(signalAction)}
            >
              <span class="font-mono text-xs">{signalAction.signal}</span>
              <span class="text-muted-foreground">{signalAction.label}</span>
            </ContextMenu.Item>
          {/each}
        </ContextMenu.SubContent>
      </ContextMenu.Sub>
    </ContextMenu.Group>
  </ContextMenu.Content>
</ContextMenu.Root>

<ProcessActionDialog bind:open {serverId} {pid} processName={name} {action} />
