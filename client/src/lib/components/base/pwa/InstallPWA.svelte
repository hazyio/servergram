<script lang="ts">
  import { onMount } from "svelte";
  import Download from "@lucide/svelte/icons/download";
  import Share from "@lucide/svelte/icons/share";
  import Button from "#lib/components/ui/button/button.svelte";
  import * as Popover from "#lib/components/ui/popover/index.js";

  interface BeforeInstallPromptEvent extends Event {
    prompt: () => Promise<void>;
    userChoice: Promise<{
      outcome: "accepted" | "dismissed";
      platform: string;
    }>;
  }

  let mounted = $state(false);
  let installed = $state(false);
  let isIOS = $state(false);
  let deferredPrompt = $state<BeforeInstallPromptEvent | null>(null);

  const canPrompt = $derived(deferredPrompt !== null);
  const visible = $derived(mounted && !installed && (canPrompt || isIOS));

  function isStandalone() {
    return (
      window.matchMedia("(display-mode: standalone)").matches ||
      (window.navigator as Navigator & { standalone?: boolean }).standalone ===
        true
    );
  }

  function detectIOS() {
    const userAgent = window.navigator.userAgent;
    return (
      /iphone|ipad|ipod/i.test(userAgent) ||
      (/macintosh/i.test(userAgent) && window.navigator.maxTouchPoints > 1)
    );
  }

  async function install() {
    if (!deferredPrompt) return;

    await deferredPrompt.prompt();
    const choice = await deferredPrompt.userChoice;

    if (choice.outcome === "accepted") {
      installed = true;
    }
    deferredPrompt = null;
  }

  onMount(() => {
    mounted = true;
    installed = isStandalone();
    isIOS = detectIOS();

    const onBeforeInstallPrompt = (event: Event) => {
      event.preventDefault();
      deferredPrompt = event as BeforeInstallPromptEvent;
    };

    const onInstalled = () => {
      installed = true;
      deferredPrompt = null;
    };

    window.addEventListener("beforeinstallprompt", onBeforeInstallPrompt);
    window.addEventListener("appinstalled", onInstalled);

    return () => {
      window.removeEventListener("beforeinstallprompt", onBeforeInstallPrompt);
      window.removeEventListener("appinstalled", onInstalled);
    };
  });
</script>

{#if visible}
  {#if canPrompt}
    <Button
      variant="ghost"
      size="icon"
      class="size-9 rounded-md text-muted-foreground hover:text-foreground"
      title="Install ServerGram"
      aria-label="Install ServerGram"
      onclick={install}
    >
      <Download class="size-4" />
    </Button>
  {:else if isIOS}
    <Popover.Root>
      <Popover.Trigger>
        <Button
          variant="ghost"
          size="icon"
          class="size-9 rounded-md text-muted-foreground hover:text-foreground"
          title="Install ServerGram"
          aria-label="Install ServerGram"
        >
          <Download class="size-4" />
        </Button>
      </Popover.Trigger>
      <Popover.Content align="end" class="w-80">
        <Popover.Header>
          <Popover.Title>Install ServerGram</Popover.Title>
          <Popover.Description>
            Keep ServerGram one tap away and open it in its own window.
          </Popover.Description>
        </Popover.Header>
        <p class="mt-1 text-sm leading-6 text-muted-foreground">
          In Safari, tap <Share class="inline size-4 align-text-bottom" /> then
          <span class="font-medium text-foreground">Add to Home Screen</span>.
        </p>
      </Popover.Content>
    </Popover.Root>
  {/if}
{/if}
