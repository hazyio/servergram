<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import ErrorSection from "#lib/components/base/ErrorSection.svelte";
  import Highlighter from "#lib/components/base/HighLighter.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import { Badge } from "#lib/components/ui/badge/index.ts";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import { Copy, Loader } from "@lucide/svelte";
  import { toast } from "svelte-sonner";
  import {
    cronJobDisplayName,
    getCronRunLog,
    type CronJob,
    type CronRun,
  } from "../cronJobs.ts";

  let {
    open = $bindable(false),
    serverId,
    job,
    run,
  }: {
    open: boolean;
    serverId: number;
    job: CronJob;
    run: CronRun | null;
  } = $props();

  let loading = $state(false);
  let loaded = $state(false);
  let loadError = $state<string | null>(null);
  let content = $state("");

  async function load() {
    if (!run) return;

    loading = true;
    loaded = false;
    loadError = null;
    content = "";

    const { data, error } = await getCronRunLog(serverId, job.id, run.id);
    loading = false;

    if (error || !data) {
      loadError = error
        ? unWrapErrorMessage(error, "Failed to load run logs")
        : "Failed to load run logs";
      return;
    }

    content = data.content;
    loaded = true;
  }

  function close() {
    open = false;
  }

  async function copy() {
    if (!content) return;

    try {
      await navigator.clipboard.writeText(content);
      toast.success("Logs copied to clipboard");
    } catch {
      toast.error("Could not copy the logs");
    }
  }

  $effect(() => {
    if (open && run) load();
  });

  $effect(() => {
    if (!open) {
      loading = false;
      loaded = false;
      loadError = null;
      content = "";
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content
    class="sm:max-w-3xl"
    showCloseButton={false}
    onEscapeKeydown={() => close()}
    onInteractOutside={(event) => event.preventDefault()}
  >
    <Dialog.Header>
      <Dialog.Title>
        <span class="truncate font-medium">Run #{run?.id} logs</span>
      </Dialog.Title>
      <div class="flex items-center gap-2">
        <p
          class="truncate font-mono text-xs text-muted-foreground"
          title={job.command}
        >
          {cronJobDisplayName(job)} · {job.command}
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
        <ErrorSection
          title="Failed to load logs"
          message={loadError}
          onRetry={load}
        />
      {:else if loading && !content}
        <Loading loading containerClass="h-full" />
      {:else}
        <div
          class="max-h-[60vh] min-h-[50vh] overflow-y-auto overflow-x-hidden rounded-md border text-xs leading-relaxed"
        >
          <Highlighter
            class="min-h-[50vh]"
            code={content}
            fileName={`run-${run?.id}.log`}
          />
        </div>
      {/if}
    </div>

    <Dialog.Footer>
      <Button variant="outline" onclick={close}>Close</Button>
      <Button disabled={!content || loading} onclick={copy}>
        <Copy class="size-4" />
        Copy
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
