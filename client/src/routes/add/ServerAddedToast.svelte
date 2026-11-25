<script lang="ts">
  import Button from "#lib/components/ui/button/button.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import { CircleCheck, Monitor, X } from "@lucide/svelte";

  interface ToastProps {
    onConnect?: () => void;
    onSeeServer?: () => void;
    onDismiss?: () => void;
  }

  let { onConnect, onSeeServer, onDismiss }: ToastProps = $props();

  function handleConnect() {
    onConnect?.();
  }

  function handleSeeServer() {
    onSeeServer?.();
  }

  function handleDismiss() {
    onDismiss?.();
  }
</script>

<Card class="w-full max-w-sm shadow-lg border-border/50">
  <div class="flex items-start gap-3 p-4">
    <div class="shrink-0 mt-0.5">
      <CircleCheck
        class="size-5 text-green-500 dark:text-green-400"
        aria-hidden="true"
      />
    </div>
    <div class="flex-1 min-w-0">
      <p class="font-medium text-base text-foreground">
        Server added successfully
      </p>
      <p class="text-sm text-muted-foreground mt-1">
        Your server has been configured and is ready to use.
      </p>
      <div class="flex flex-wrap gap-2 mt-4">
        <Button variant="default" size="sm" onclick={handleConnect}>
          <Monitor class="size-3.5 mr-1.5" aria-hidden="true" />
          Connect
        </Button>
        <Button variant="outline" size="sm" onclick={handleSeeServer}>
          See Server
        </Button>
        <Button variant="ghost" size="sm" onclick={handleDismiss}>
          <X class="size-3.5" aria-hidden="true" />
          <span class="sr-only">Dismiss</span>
        </Button>
      </div>
    </div>
  </div>
</Card>
