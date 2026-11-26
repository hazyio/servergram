<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import { ChevronDown, Download, Loader } from "@lucide/svelte";
  import { untrack } from "svelte";
  import { toast } from "svelte-sonner";
  import {
    CRON_RUN_LOG_PAGE_SIZE,
    getCronRunLog,
    type CronRun,
  } from "../cronJobs.ts";

  let {
    open = $bindable(false),
    serverId,
    jobName,
    run,
  }: {
    open: boolean;
    serverId: number;
    jobName: string;
    run: CronRun | null;
  } = $props();

  let content = $state("");
  let hasMore = $state(false);
  let offset = $state(0);
  let loading = $state(false);
  let loadingMore = $state(false);
  let loadError = $state("");
  let requestToken = 0;

  const title = $derived(
    run ? `${jobName ? `${jobName} · ` : ""}Run #${run.id}` : "Run log",
  );
  const canDownload = $derived(!!content);

  function reset() {
    content = "";
    hasMore = false;
    offset = 0;
    loading = false;
    loadingMore = false;
    loadError = "";
  }

  async function load() {
    if (!run) return;

    const token = ++requestToken;
    reset();
    loading = true;

    const { data, error } = await getCronRunLog(serverId, run.id, 0);

    if (token !== requestToken) return;

    loading = false;

    if (error || !data) {
      loadError = error
        ? unWrapErrorMessage(error, "Failed to load run log")
        : "Failed to load run log";
      return;
    }

    content = data.content;
    // The log endpoint has no "more" flag. Keep the button available until a
    // request comes back empty, which means the end of the log was reached.
    hasMore = content.length > 0;
    offset = CRON_RUN_LOG_PAGE_SIZE;
  }

  async function loadMore() {
    if (!run || loadingMore || !hasMore) return;

    const token = requestToken;
    const limit = CRON_RUN_LOG_PAGE_SIZE;
    loadingMore = true;

    const { data, error } = await getCronRunLog(
      serverId,
      run.id,
      offset,
      limit,
    );

    if (token !== requestToken) return;

    loadingMore = false;

    if (error || !data) {
      // Keep the output we already have and leave the button in place so the
      // user can retry.
      toast.error(
        error
          ? unWrapErrorMessage(error, "Failed to load more run log")
          : "Failed to load more run log",
      );
      return;
    }

    if (!data.content) {
      // Empty log means there is nothing left to load.
      hasMore = false;
      return;
    }

    content += data.content;
    offset += limit;
  }

  function download() {
    if (!content) return;

    const url = URL.createObjectURL(
      new Blob([content], { type: "text/plain" }),
    );
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = `cron-run-${run?.id ?? "log"}.log`;
    anchor.click();
    URL.revokeObjectURL(url);
  }

  $effect.pre(() => {
    if (!open) {
      requestToken++;
      reset();
      return;
    }

    untrack(() => {
      void load();
    });
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="sm:max-w-3xl">
    <Dialog.Header>
      <Dialog.Title>{title}</Dialog.Title>
      <Dialog.Description>
        Standard output and errors captured while this run was executing.
      </Dialog.Description>
    </Dialog.Header>

    <div class="min-h-40">
      {#if loading}
        <Loading
          loading
          containerClass="h-40"
          class="h-8 w-8"
          text="Loading log..."
        />
      {:else if loadError}
        <ErrorSection
          title="Failed to load run log"
          message={loadError}
          onRetry={load}
        />
      {:else}
        <div
          class="max-h-[60vh] overflow-auto rounded-md border bg-muted/30 p-3 font-mono text-xs leading-relaxed whitespace-pre-wrap break-words"
        >
          {#if content}
            {content}
          {:else}
            <span class="text-muted-foreground">
              No output was recorded for this run.
            </span>
          {/if}
        </div>
        {#if hasMore || loadingMore}
          <div class="mt-3 flex justify-center">
            <Button
              variant="outline"
              size="sm"
              class="gap-1.5"
              disabled={loadingMore}
              onclick={loadMore}
            >
              {#if loadingMore}
                <Loader class="size-4 animate-spin" />
                Loading more...
              {:else}
                <ChevronDown class="size-4" />
                Load more
              {/if}
            </Button>
          </div>
        {/if}
      {/if}
    </div>

    <Dialog.Footer>
      <Button variant="outline" onclick={() => (open = false)}>Close</Button>
      <Button onclick={download} disabled={!canDownload}>
        <Download class="size-4" />
        Download
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
