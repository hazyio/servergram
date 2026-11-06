<script lang="ts">
  import { Badge } from "#lib/components/ui/badge/index.ts";
  import { Button } from "#lib/components/ui/button/index.ts";

  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { joinPath } from "#lib/utils.ts";
  import { Download, Loader } from "@lucide/svelte";

  // Files larger than this are rejected by the agent instead of being streamed.
  const MAX_PREVIEW_BYTES = 5 * 1024 * 1024;

  // Go encodes []byte as base64 in JSON.
  type ReadFileChunk = components["schemas"]["ReadFileResponse"];
  type ReadFileError = components["schemas"]["SSEError"];

  let {
    serverId,
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
  let eventSource: EventSource | undefined;

  const path = $derived(joinPath(dir, entryName ?? ""));

  function stopStream() {
    eventSource?.close();
    eventSource = undefined;
  }

  function decodeBase64(value: string) {
    const binary = atob(value);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    return bytes;
  }

  function loadFile() {
    if (entryName === null) return;

    stopStream();
    loading = true;
    loaded = false;
    loadError = null;
    content = "";

    const query = new URLSearchParams({
      path: dir,
      name: entryName,
      max_bytes_size: String(MAX_PREVIEW_BYTES),
    });
    const source = new EventSource(
      `/api/files/read-file-in-server/${serverId}?${query}`,
      { withCredentials: true },
    );
    eventSource = source;

    const decoder = new TextDecoder();

    function handleChunk(event: MessageEvent<string>) {
      if (source !== eventSource) return;
      try {
        const chunk = JSON.parse(event.data) as ReadFileChunk;
        if (chunk.data) {
          // Stream-decode so multi-byte characters split across chunks survive.
          content += decoder.decode(decodeBase64(chunk.data), {
            stream: true,
          });
        }
        if (chunk.done) {
          content += decoder.decode();
          loaded = true;
          loading = false;
          stopStream();
        }
      } catch {
        loadError = "Failed to parse file data";
        loading = false;
        stopStream();
      }
    }

    source.addEventListener("data", (event) =>
      handleChunk(event as MessageEvent<string>),
    );
    // Huma sends named "data" events; also accept the default "message" event.
    source.onmessage = handleChunk;

    // Server-sent "error" events carry data; connection errors do not.
    source.addEventListener("error", (event) => {
      if (source !== eventSource || !(event instanceof MessageEvent)) return;
      try {
        const response = JSON.parse(event.data) as ReadFileError;
        loadError = response.message ?? "Failed to read file";
      } catch {
        loadError = "Failed to read file";
      }
      loading = false;
      stopStream();
    });

    source.onerror = () => {
      if (source !== eventSource || source.readyState !== EventSource.CLOSED) {
        return;
      }
      loading = false;
      if (!loaded && !loadError) {
        loadError = "File stream closed unexpectedly";
      }
      stopStream();
    };
  }

  function download() {
    if (!loaded) return;
    // TODO: use a dedicated download endpoint for binary/large files.
    const url = URL.createObjectURL(
      new Blob([content], { type: "text/plain" }),
    );
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = entryName ?? "file.txt";
    anchor.click();
    URL.revokeObjectURL(url);
  }

  function close() {
    stopStream();
    entryName = null;
    loaded = false;
    loadError = null;
    content = "";
  }

  $effect(() => {
    if (entryName === null) return;
    loadFile();
    return stopStream;
  });
</script>

<Dialog.Root open={!!entryName}>
  <Dialog.Content
    class="sm:max-w-3xl"
    showCloseButton={false}
    onEscapeKeydown={(event) => {
      close();
    }}
    onInteractOutside={(event) => event.preventDefault()}
  >
    <Dialog.Header>
      <Dialog.Title>
        <span class="truncate font-medium" title={entryName ?? ""}
          >{entryName}</span
        >
      </Dialog.Title>
      <div class="flex items-center gap-2">
        <p
          class="truncate font-mono text-xs text-muted-foreground"
          title={path}
        >
          {path}
        </p>
        {#if loading}
          <Badge variant="secondary" class="shrink-0 gap-1">
            <Loader class="size-3 animate-spin" />
            Receiving
          </Badge>
        {:else if loaded}
          <Badge variant="outline" class="shrink-0">Loaded</Badge>
        {/if}
      </div>
    </Dialog.Header>

    <div class="min-h-50">
      {#if loadError}
        <ErrorSection title="Failed to load file" message={loadError} />
      {:else if loading && !content}
        <Loading loading containerClass="h-full" />
      {:else}
        <pre
          class="min-h-[50vh] max-h-[60vh] overflow-auto whitespace-pre-wrap wrap-break-word rounded-md border bg-muted/40 p-4 text-xs leading-relaxed"><code
            >{content}</code
          ></pre>
      {/if}
    </div>

    <Dialog.Footer>
      <Button variant="outline" onclick={close}>Close</Button>
      <Button disabled={loading || !content} onclick={download}>
        <Download class="size-4" />
        Download
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
