<script lang="ts">
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import Input from "#lib/components/ui/input/input.svelte";
  import type { WSServicePassPhraseRequiredData } from "#lib/types/websockets/index.js";
  import { Eye, EyeOff } from "@lucide/svelte";
  import { Button } from "../../ui/button";
  let {
    data = $bindable([]),
  }: {
    data: WSServicePassPhraseRequiredData[];
  } = $props();
  let submitting = false;
  let showPassword = $state(false);
  let password = $state("");
</script>

{#if data.length > 0}
  <Dialog.Root open={true}>
    <Dialog.Content
      showCloseButton={false}
      onEscapeKeydown={(event) => event.preventDefault()}
      onInteractOutside={(event) => event.preventDefault()}
    >
      <Dialog.Header>
        <Dialog.Title
          >SSH passphrase is needed to connect to server {data[0]
            .name}</Dialog.Title
        >
        <Dialog.Description>
          Enter the SSH passphrase for the server to continue.
        </Dialog.Description>
      </Dialog.Header>
      <div class="relative">
        <Input
          placeholder="**********"
          disabled={submitting}
          type={showPassword ? "text" : "password"}
          autocomplete="new-password"
          bind:value={password}
        />
        <Button
          type="button"
          variant="ghost"
          size="icon"
          class="absolute top-0 right-0 h-full px-3 hover:bg-transparent"
          onclick={() => (showPassword = !showPassword)}
        >
          {#if showPassword}
            <EyeOff class="h-4 w-4" />
          {:else}
            <Eye class="h-4 w-4" />
          {/if}
        </Button>
      </div>
      <Dialog.Footer>
        <Button
          disabled={submitting}
          onclick={() => {
            password = "";
            data = data.slice(1);
          }}>Cancel</Button
        >
        <Button disabled={submitting}>Submit</Button>
      </Dialog.Footer>
    </Dialog.Content>
  </Dialog.Root>
{/if}
