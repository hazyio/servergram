<script lang="ts">
  import { onMount } from "svelte";
  import ErrorSection from "../ErrorSection.svelte";

  let online = $state(true);
  let mounted = $state(false);

  function syncConnection() {
    online = navigator.onLine;
  }

  function retry() {
    syncConnection();
    if (online) {
      location.reload();
    }
  }

  onMount(() => {
    mounted = true;
    syncConnection();

    window.addEventListener("online", syncConnection);
    window.addEventListener("offline", syncConnection);

    return () => {
      window.removeEventListener("online", syncConnection);
      window.removeEventListener("offline", syncConnection);
    };
  });
</script>

{#if mounted && !online}
  <div
    class="fixed inset-x-0 bottom-0 z-50 border-t border-border bg-background shadow-[0_-12px_40px_rgba(0,0,0,0.35)]"
    aria-live="assertive"
  >
    <ErrorSection
      title="You're offline"
      message="ServerGram can't reach the server right now. Cached pages still work, and the app will reconnect automatically when your connection returns."
      retryText="Retry connection"
      onRetry={retry}
      class="border-b-0"
    />
  </div>
{/if}
