<script lang="ts">
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import type { components } from "#lib/types/request-schema.d.ts";
  import { untrack } from "svelte";

  import * as ContextMenu from "#lib/components/ui/context-menu/index.js";
  import {
    ChevronLeft,
    ChevronRight,
    Eye,
    EyeOff,
    FileIcon,
    FolderIcon,
    Loader,
    RefreshCcw,
    Upload,
  } from "@lucide/svelte";
  import type { PageProps } from "./$types";
  import Actions from "./Actions.svelte";
  import DeleteEntriesDialog, {
    type DeleteEntryProps,
  } from "./DeleteEntriesDialog.svelte";
  import { isDir, type FileEntry } from "./FilesTable/filesColumns";
  import FilesTable from "./FilesTable/FilesTable.svelte";
  import NewEntry, { type NewEntryType } from "./NewEntry.svelte";
  import SettingOption, {
    type SettingOptionProps,
  } from "./SettingOption.svelte";
  import { severPageInfo } from "#lib/states/severPageInfo.svelte.ts";
  import RenameEntry from "./RenameEntry.svelte";

  type Chunk = components["schemas"]["ListFilesChunkResponse"];
  type WatchEvent = components["schemas"]["WatchDirectoryResponse"];

  let { params }: PageProps = $props();
  let settings = $state<SettingOptionProps>({
    showHidden: false,
    dirsFirst: true,
    columns: {
      size: true,
      modified: true,
      mode: true,
      owner: true,
    },
    askBeforeDelete: true,
  });
  let search = $state("");
  let dir = $state(severPageInfo.lastVisitedDir);
  let history = $state([""]);
  let historyIndex = $state(0);
  let loading = $state(true);
  let updating = $state(false);
  let loadError = $state<string | null>(null);
  let fileStreamError = $state<string | null>(null);
  let files: FileEntry[] = $state([]);
  let eventSource: EventSource | undefined;
  let watchSource: EventSource | undefined;
  let watching = $state(false);
  let deleteEntries = $state<DeleteEntryProps[]>([]);
  let newEntryType = $state<NewEntryType>("none");
  let renameOldName = $state<null | string>(null);

  function stop() {
    stopWatch();
    eventSource?.close();
    eventSource = undefined;
    updating = false;
  }

  function stopWatch() {
    watchSource?.close();
    watchSource = undefined;
    watching = false;
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
          // update dir and last visited dir
          dir = chunk.resolved_path;
          severPageInfo.lastVisitedDir = chunk.resolved_path;
          history[historyIndex] = dir;
        }
        loading = false;
        if (chunk.done) {
          stop();
          // start watching for changes after the initial load is complete,
          watchFiles();
        } else {
          updating = true;
        }
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
      fileStreamError ??= "File stream closed, an unexpected error occurred";
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

  function watchFiles() {
    if (watching) return;

    stopWatch();
    watching = true;

    const query = new URLSearchParams({
      dir,
      include_hidden: String(settings.showHidden),
    });
    const source = new EventSource(
      `/api/files/watch-dir-in-server/${params.id}?${query}`,
      { withCredentials: true },
    );
    watchSource = source;

    source.addEventListener("data", (event) => {
      if (source !== watchSource) return;
      try {
        const watchEvent = JSON.parse(
          (event as MessageEvent<string>).data,
        ) as WatchEvent;
        handleWatchEvent(watchEvent);
      } catch {
        console.error("Failed to parse watch event");
      }
    });

    source.addEventListener("error", (event) => {
      if (source !== watchSource || !(event instanceof MessageEvent)) return;
      try {
        const response = JSON.parse(event.data) as { message?: string };
        console.error(
          "Watch error:",
          response.message ?? "Failed to watch directory",
        );
      } catch {
        console.error("Watch error: Failed to watch directory");
      }
      stopWatch();
    });

    source.onerror = () => {
      if (source !== watchSource || source.readyState !== EventSource.CLOSED) {
        return;
      }
      stopWatch();
      console.warn("Watch stream closed");
    };
  }
  function removeFile(pathName: string) {
    files = files.filter((f) => f.name !== pathName);
  }
  function handleWatchEvent(event: WatchEvent) {
    const {
      kind,
      path,
      entry,
    }: components["schemas"]["WatchDirectoryResponse"] = event;

    switch (kind) {
      case 0: // MODIFIED
        console.warn("Unimplemented");
        break;
      case 1: // CREATED
        // Add the new file/directory to the list
        if (entry) {
          files = [...files, entry];
        } else {
          loadFiles();
        }
        break;
      case 2: // DELETED
        // Remove the file/directory from the list
        if (path) {
          removeFile(path);
        }
        break;
    }
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
    placeholder="Search here..."
    autocomplete="off"
    type="search"
    class="shrink"
    oninput={(e) => (search = e.currentTarget.value)}
  />
  <SettingOption
    bind:showHidden={settings.showHidden}
    bind:dirsFirst={settings.dirsFirst}
    bind:columns={settings.columns}
    bind:askBeforeDelete={settings.askBeforeDelete}
  />
  <Button variant="outline" onclick={loadFiles}><RefreshCcw /></Button>
  <Actions bind:entryType={newEntryType} />

  <Button
    variant={watching ? "default" : "outline"}
    onclick={watching ? stopWatch : watchFiles}
    class="flex items-center gap-1"
    title={watching ? "Stop watching for changes" : "Watch for changes"}
  >
    {#if watching}
      <Eye class="size-4" />
    {:else}
      <EyeOff class="size-4" />
    {/if}
    <span class="hidden sm:inline">{watching ? "Watching" : "Watch"}</span>
  </Button>
</div>

<DeleteEntriesDialog
  serverId={Number(params.id)}
  {dir}
  bind:items={deleteEntries}
  onDeleted={(name) => {
    // client-side update, backend will send an update but we cannot depend on that since watching can be turned off
    removeFile(name);
  }}
/>
<NewEntry serverId={Number(params.id)} bind:type={newEntryType} {dir} />
<RenameEntry serverId={Number(params.id)} bind:oldName={renameOldName} {dir} />

{#if loading}
  <Loading loading class="mt-20 w-12 h-12" />
{:else if loadError}
  <ErrorSection message={loadError} onRetry={loadFiles} class="my-7" />
{:else}
  {#if updating}
    <div class="mt-2 flex items-center gap-2 text-sm text-muted-foreground">
      <Loader class="size-4 animate-spin" />
      Updating...
    </div>
  {/if}

  <ContextMenu.Root>
    <ContextMenu.Trigger
      ><div class="p-4 mb-10 px-12 min-h-[70vh]">
        <FilesTable
          {files}
          {search}
          dirsFirst={settings.dirsFirst}
          visibleColumns={settings.columns}
          onOpen={openEntry}
          onDelete={(entry) => {
            deleteEntries = [entry];
          }}
          onRenameClicked={(entry) => {
            renameOldName = entry.name;
          }}
        />
      </div>
    </ContextMenu.Trigger>
    <ContextMenu.Content class="min-h-30 w-50">
      <ContextMenu.Group>
        <ContextMenu.Label
          >Actions on <span class="truncate" title={dir}>{dir}</span
          ></ContextMenu.Label
        >
        <ContextMenu.Item
          class="py-2"
          onclick={() => {
            newEntryType = "file";
          }}
        >
          <FileIcon />
          New File
        </ContextMenu.Item>
        <ContextMenu.Item
          class="py-2"
          onclick={() => {
            newEntryType = "directory";
          }}
        >
          <FolderIcon />
          New Directory
        </ContextMenu.Item>
        <ContextMenu.Item class="py-2">
          <Upload />
          Upload
        </ContextMenu.Item>
        <ContextMenu.Item class="py-2" onclick={loadFiles}>
          <RefreshCcw />
          Refresh
        </ContextMenu.Item>
      </ContextMenu.Group>
    </ContextMenu.Content>
  </ContextMenu.Root>
{/if}
