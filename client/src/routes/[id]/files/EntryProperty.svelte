<script lang="ts">
  import { Badge } from "#lib/components/ui/badge/index.ts";
  import { Button } from "#lib/components/ui/button/index.ts";
  import { Separator } from "#lib/components/ui/separator/index.js";

  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import { formatBytes, formatMode } from "#lib/utils.ts";
  import type { EntryKind, FileEntry } from "#lib/types/index.ts";
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";

  const kindLabels: Record<number, string> = {
    0: "File",
    1: "Directory",
    2: "Symlink",
  };

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
  let entry = $state<FileEntry | null>(null);

  async function loadProperties() {
    if (!entryName) return;
    loading = true;
    loadError = null;
    const { data, error } = await NewRequestClient().GET(
      "/api/files/get-path-properties-in-server/{id}",
      {
        params: {
          path: {
            id: serverId,
          },
          query: {
            name: entryName,
            path: dir,
          },
        },
      },
    );
    loading = false;

    if (error) {
      loadError = unWrapErrorMessage(error);
      return;
    }
    entry = data.Enity;
  }

  $effect(() => {
    if (entryName !== null && !loaded) {
      loadProperties();
    }
  });
</script>

<Dialog.Root
  open={!!entryName}
  onOpenChange={(open) => {
    if (!open) {
      entryName = null;
      loaded = false;
      loadError = null;
    }
  }}
>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>
        {#if loading || entry === null}
          Loading {entryName} Properties
        {:else}
          <span class="truncate text-base font-medium" title={entry.name}
            >{entry.name}</span
          >
          <Badge variant="outline">{kindLabels[entry.kind ?? 0]}</Badge>
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
            <dd class="break-all font-mono">{dir}</dd>

            <dt class="text-muted-foreground">Type</dt>
            <dd>{kindLabels[entry.kind ?? 0]}</dd>

            <dt class="text-muted-foreground">Size</dt>
            <dd>{formatBytes(entry.size ?? 0)}</dd>

            <dt class="text-muted-foreground">Modified</dt>
            <dd class="break-all font-mono text-xs">
              {new Date(
                entry.modified_unix_ns ? entry.modified_unix_ns / 1e6 : 0,
              ).toLocaleString()}
            </dd>

            <dt class="text-muted-foreground">Mode</dt>
            <dd class="font-mono text-xs">
              {formatMode(entry)}
            </dd>

            <dt class="text-muted-foreground">Owner</dt>
            <dd class="font-mono text-xs">{entry.uid}:{entry.gid}</dd>

            <!-- {#if entry.symlinkTarget}
              <dt class="text-muted-foreground">Target</dt>
              <dd class="break-all font-mono text-xs">
                {entry.symlinkTarget}
              </dd>
            {/if} -->
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
