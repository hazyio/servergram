<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import { Download } from "@lucide/svelte";
  import { untrack } from "svelte";
  import { getCronRunLog, type CronRun, type CronRunLog } from "../cronJobs.ts";

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

  let log = $state<CronRunLog | null>(null);
  let loading = $state(false);
  let loadError = $state("");
  let requestToken = 0;

  const title = $derived(
    run ? `${jobName ? `${jobName} · ` : ""}Run #${run.id}` : "Run log",
  );
  const canDownload = $derived(!!log?.content);

  async function load() {
    if (!run) return;

    const token = ++requestToken;
    loading = true;
    loadError = "";
    log = null;

    const { data, error } = await getCronRunLog(serverId, run.job_id, run.id);

    if (token !== requestToken) return;

    loading = false;

    if (error || !data) {
      loadError = error
        ? unWrapErrorMessage(error, "Failed to load run log")
        : "Failed to load run log";
      return;
    }

    log = data;
  }

  function download() {
    if (!log?.content) return;

    const url = URL.createObjectURL(
      new Blob([log.content], { type: "text/plain" }),
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
      loading = false;
      loadError = "";
      log = null;
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
          {#if log?.content}
            {log.content}
          {:else}
            <span class="text-muted-foreground">
              No output was recorded for this run.
            </span>
          {/if}
        </div>
        {#if log?.truncated}
          <p class="mt-2 text-xs text-muted-foreground">
            Output was truncated to the last part of the log.
          </p>
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
