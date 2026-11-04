<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { untrack } from "svelte";

  import {
    ChevronLeft,
    ChevronRight,
    FileIcon,
    FolderIcon,
    Loader,
    RefreshCcw,
    Upload,
  } from "@lucide/svelte";
  import type { PageProps } from "./$types";
  import { isDir, type FileEntry } from "./FilesTable/filesColumns";
  import FilesTable from "./FilesTable/FilesTable.svelte";
  import SettingOption, {
    type SettingOptionProps,
  } from "./SettingOption.svelte";
  import Actions from "./Actions.svelte";
  import * as ContextMenu from "#lib/components/ui/context-menu/index.js";

  type Chunk = components["schemas"]["ListFilesChunkResponse"];

  let { params }: PageProps = $props();
  let settings = $state<SettingOptionProps>({
    showHidden: false,
    dirsFirst: true,
    columns: { size: true, modified: true, mode: true, owner: true },
  });
  let search = $state("");
  let dir = $state("");
  let history = $state([""]);
  let historyIndex = $state(0);
  let loading = $state(true);
  let updating = $state(false);
  let loadError = $state<string | null>(null);
  let files: FileEntry[] = $state([]);
  let eventSource: EventSource | undefined;

  function stop() {
    eventSource?.close();
    eventSource = undefined;
    updating = false;
  }

  function loadFiles() {
    stop();
    loading = true;
    loadError = null;
    files = [];

    const query = new URLSearchParams({
      dir,
      include_hidden: String(settings.showHidden),
    });
    const source = new EventSource(
      `/api/files/list-in-server/${params.id}?${query}`,
      { withCredentials: true },
    );
    eventSource = source;

    source.addEventListener("data", (event) => {
      if (source !== eventSource) return;
      try {
        const chunk = JSON.parse((event as MessageEvent<string>).data) as Chunk;
        if (chunk.entries?.length) files = [...files, ...chunk.entries];
        if (chunk.resolved_path) {
          dir = chunk.resolved_path;
          history[historyIndex] = dir;
        }
        loading = false;
        if (chunk.done) stop();
        else updating = true;
      } catch {
        stop();
        loading = false;
        loadError = "Failed to parse file data";
      }
    });

    // Server-sent "error" events carry data; connection errors do not.
    source.addEventListener("error", (event) => {
      if (source !== eventSource || !(event instanceof MessageEvent)) return;
      try {
        const response = JSON.parse(event.data) as { message?: string };
        loadError = response.message ?? "Failed to load files";
      } catch {
        loadError = "Failed to load files";
      }
      stop();
      loading = false;
    });

    source.onerror = () => {
      if (source !== eventSource || source.readyState !== EventSource.CLOSED) {
        return;
      }
      stop();
      loading = false;
      loadError ??= "File stream closed, an unexpected error occurred";
    };
  }

  function navigate(path: string) {
    history = [...history.slice(0, historyIndex + 1), path];
    historyIndex = history.length - 1;
    dir = path;
    loadFiles();
  }

  function go(offset: -1 | 1) {
    historyIndex += offset;
    dir = history[historyIndex];
    loadFiles();
  }

  function openEntry(entry: FileEntry) {
    if (!isDir(entry) || !entry.name) return;
    navigate(`${dir.replace(/\/+$/, "")}/${entry.name}`);
  }

  $effect(() => {
    void settings.showHidden;
    untrack(loadFiles);
    return stop;
  });
</script>

<div class="flex gap-2">
  <div class="flex gap-2">
    <Button
      variant="outline"
      disabled={historyIndex <= 0}
      onclick={() => go(-1)}><ChevronLeft /></Button
    >
    <Button
      variant="outline"
      disabled={historyIndex >= history.length - 1}
      onclick={() => go(1)}><ChevronRight /></Button
    >
  </div>
  <Input
    bind:value={dir}
    class="grow-2"
    onkeydown={(e: KeyboardEvent) => e.key === "Enter" && navigate(dir)}
  />
  <Input
    name="filter-files"
    placeholder="Search files here..."
    autocomplete="off"
    type="search"
    class="shrink"
    oninput={(e) => (search = e.currentTarget.value)}
  />
  <Actions />

  <SettingOption
    bind:showHidden={settings.showHidden}
    bind:dirsFirst={settings.dirsFirst}
    bind:columns={settings.columns}
  />
  <Button variant="outline" onclick={loadFiles}><RefreshCcw /></Button>
</div>

{#if loading}
  <Loading loading class="my-7" />
{:else if loadError}
  <ErrorSection message={loadError} onRetry={loadFiles} class="my-7" />
{:else}
  {#if updating}
    <div class="mt-2 flex items-center gap-2 text-sm text-muted-foreground">
      <Loader class="size-4 animate-spin" />
      Updating...
    </div>
  {/if}
  <div class="p-4 mb-10 px-12">
    <ContextMenu.Root>
      <ContextMenu.Trigger>
        <FilesTable
          {files}
          {search}
          dirsFirst={settings.dirsFirst}
          visibleColumns={settings.columns}
          onOpen={openEntry}
        />
      </ContextMenu.Trigger>
      <ContextMenu.Content class="h-50 w-50">
        <ContextMenu.Group>
          <ContextMenu.Item>
            <FileIcon />
            New File
          </ContextMenu.Item>
          <ContextMenu.Item>
            <FolderIcon />
            New Folder
          </ContextMenu.Item>
          <ContextMenu.Item>
            <Upload />
            Upload File
          </ContextMenu.Item>
          <ContextMenu.Item onclick={loadFiles}>
            <RefreshCcw />
            Refresh
          </ContextMenu.Item>
        </ContextMenu.Group>
      </ContextMenu.Content>
    </ContextMenu.Root>
  </div>
{/if}
