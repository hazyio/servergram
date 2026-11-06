<script lang="ts">
  import { Button } from "#lib/components/ui/button/index.ts";

  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import { Download } from "@lucide/svelte";

  let {
    dir,
    entryName = $bindable(),
  }: {
    serverId: number;
    dir: string;
    entryName: null | string;
  } = $props();
  let loading = $state(false);
  let loaded = $state(false);
  let loadError = $state<string | null>(null);
  let content = $state("");

  const path = $derived(`${dir.replace(/\/+$/, "")}/${entryName ?? ""}`);

  async function loadFile() {
    loading = true;
    loadError = null;
    try {
      // TODO: replace with a request for the file's contents
      content = `// Stub preview of ${entryName}
export function greet(name: string) {
  return \`Hello, \${name}!\`;
}
`;
      loaded = true;
    } catch (error) {
      loadError = error instanceof Error ? error.message : String(error);
    } finally {
      loading = false;
    }
  }

  function download() {
    if (!content) return;
    // TODO: replace with a request to the file download endpoint
    const url = URL.createObjectURL(new Blob([content], { type: "text/plain" }));
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = entryName ?? "file.txt";
    anchor.click();
    URL.revokeObjectURL(url);
  }

  function close() {
    entryName = null;
    loaded = false;
    loadError = null;
    content = "";
  }

  $effect(() => {
    if (entryName !== null && !loaded) {
      loadFile();
    }
  });
</script>

<Dialog.Root open={!!entryName}>
  <Dialog.Content
    class="sm:max-w-3xl"
    showCloseButton={false}
    onEscapeKeydown={(event) => event.preventDefault()}
    onInteractOutside={(event) => event.preventDefault()}
  >
    <Dialog.Header>
      <Dialog.Title>
        {#if loading || !loaded}
          Loading {entryName}
        {:else}
          <span class="truncate font-medium" title={entryName ?? ""}
            >{entryName}</span
          >
        {/if}
      </Dialog.Title>
      <p class="truncate font-mono text-xs text-muted-foreground" title={path}>
        {path}
      </p>
    </Dialog.Header>

    <div class="min-h-50">
      {#if loading}
        <Loading loading containerClass="h-full" />
      {:else if loadError}
        <ErrorSection title="Failed to load file" message={loadError} />
      {:else}
        <pre
          class="max-h-[60vh] overflow-auto rounded-md border bg-muted/40 p-4 text-xs leading-relaxed"><code
            >{content}</code
          ></pre>
      {/if}
    </div>

    <Dialog.Footer>
      <Button variant="outline" disabled={loading} onclick={close}>
        Close
      </Button>
      <Button disabled={loading || !content} onclick={download}>
        <Download class="size-4" />
        Download
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
