<script lang="ts">
  import * as Dialog from "#lib/components/ui/dialog/index.ts";
  import { Unplug } from "@lucide/svelte";
  import { Button } from "../../ui/button";
  import { ConnectedServices } from "./services.svelte.ts";
  let {
    id,
    name,
    type,
  }: {
    id: number;
    name: string;
    type: "server" | "database";
  } = $props();
</script>

<Dialog.Root>
  <Dialog.Trigger
    ><Button size="sm" variant="outline" title={`Disconnect ${name}`}
      ><Unplug />
    </Button></Dialog.Trigger
  >
  <Dialog.Content class="z-60">
    <Dialog.Header>
      <Dialog.Title>Disconnect {name}?</Dialog.Title>
    </Dialog.Header>
    <div>Are you sure you want to disconnect {name}?</div>
    <Dialog.Footer>
      <Dialog.Close><Button variant="secondary">No</Button></Dialog.Close>
      <Dialog.Close>
        <Button
          variant="destructive"
          onclick={() => {
            if (type === "server") {
              ConnectedServices.disconnectServer(id);
            } else if (type === "database") {
              ConnectedServices.disconnectDatabase(id);
            }
          }}>Yes</Button
        >
      </Dialog.Close>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
