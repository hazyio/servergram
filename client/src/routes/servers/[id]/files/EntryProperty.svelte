<script lang="ts">
  import { Badge } from "#lib/components/ui/badge/index.ts";
  import { Button } from "#lib/components/ui/button/index.ts";
  import { Separator } from "#lib/components/ui/separator/index.js";

  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import { formatBytes, formatMode } from "#lib/utils.ts";
  import type { EntryKind } from "#lib/types/index.ts";

  interface EntryProperties {
    name: string;
    path: string;
    kind: EntryKind;
    size: number;
    modifiedUnixNs: number;
    mode: number;
    uid: number;
    gid: number;
    hidden: boolean;
    symlinkTarget?: string;
  }

  const kindLabels: Record<EntryKind, string> = {
    file: "File",
    directory: "Directory",
    symlink: "Symlink",
  };

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
  let entry = $state<EntryProperties | null>(null);

  async function loadProperties() {
    loading = true;
    loadError = null;
    try {
      // TODO: replace with a request for the entry's real properties
      entry = {
        name: entryName ?? "",
        path: `${dir.replace(/\/+$/, "")}/${entryName ?? ""}`,
        kind: "file",
        size: 24_576,
        modifiedUnixNs: Date.now() * 1e6,
        mode: 0o644,
        uid: 1000,
        gid: 1000,
        hidden: entryName?.startsWith(".") ?? false,
      };
      loaded = true;
    } catch (error) {
      loadError = error instanceof Error ? error.message : String(error);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (entryName !== null && !loaded) {
      loadProperties();
    }
  });
</script>

<Dialog.Root open={!!entryName}>
  <Dialog.Content
    showCloseButton={false}
    onEscapeKeydown={(event) => event.preventDefault()}
    onInteractOutside={(event) => event.preventDefault()}
  >
    <Dialog.Header>
      <Dialog.Title>
        {#if loading || entry === null}
          Loading {entryName} Properties
        {:else}
          <span class="truncate text-base font-medium" title={entry.name}
            >{entry.name}</span
          >
          <Badge variant="outline">{kindLabels[entry.kind]}</Badge>
          {#if entry.hidden}
            <Badge variant="secondary">Hidden</Badge>
          {/if}
        {/if}
      </Dialog.Title>
    </Dialog.Header>

    <div class="min-h-50">
      {#if loading}
        <Loading loading containerClass="h-full" />
      {:else if loadError}
        <ErrorSection
          title="Failed to load entry properties"
          message={loadError}
        />
      {:else if entry}
        <div class="space-y-4">
          <Separator />

          <dl class="grid grid-cols-[6.5rem_1fr] gap-x-4 gap-y-3 text-sm">
            <dt class="text-muted-foreground">Path</dt>
            <dd class="break-all font-mono text-xs">{entry.path}</dd>

            <dt class="text-muted-foreground">Type</dt>
            <dd>{kindLabels[entry.kind]}</dd>

            <dt class="text-muted-foreground">Size</dt>
            <dd>{formatBytes(entry.size)}</dd>

            <dt class="text-muted-foreground">Modified</dt>
            <dd>{new Date(entry.modifiedUnixNs / 1e6).toLocaleString()}</dd>

            <dt class="text-muted-foreground">Mode</dt>
            <dd class="font-mono text-xs">
              {formatMode(entry.mode, entry.kind)}
            </dd>

            <dt class="text-muted-foreground">Owner</dt>
            <dd class="font-mono text-xs">{entry.uid}:{entry.gid}</dd>

            {#if entry.symlinkTarget}
              <dt class="text-muted-foreground">Target</dt>
              <dd class="break-all font-mono text-xs">
                {entry.symlinkTarget}
              </dd>
            {/if}
          </dl>
        </div>
      {/if}
    </div>
    <Dialog.Footer>
      <Button
        disabled={loading}
        onclick={() => {
          entryName = null;
          loaded = false;
          loadError = null;
        }}>Close</Button
      >
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
