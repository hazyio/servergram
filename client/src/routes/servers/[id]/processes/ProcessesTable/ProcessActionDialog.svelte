<script lang="ts">
  import { unWrapErrorMessage } from "#lib/api/request.ts";
  import Loading from "#lib/components/base/loading.svelte";
  import * as AlertDialog from "#lib/components/ui/alert-dialog/index.js";
  import { Button } from "#lib/components/ui/button/index.js";
  import { untrack } from "svelte";
  import { toast } from "svelte-sonner";
  import { sendProcessSignal, type ProcessAction } from "./processActions";

  let {
    open = $bindable(false),
    serverId,
    pid,
    processName,
    action,
  }: {
    open: boolean;
    serverId: number;
    pid: number;
    processName: string;
    action: ProcessAction | null;
  } = $props();

  let loading = $state(false);
  let error = $state("");
  let requestToken = 0;

  const loadingText = $derived(
    action
      ? `Sending ${action.signal} to ${processName} (PID ${pid})...`
      : "Sending signal...",
  );

  async function run(currentAction: ProcessAction) {
    const token = ++requestToken;
    loading = true;
    error = "";

    const { error: requestError } = await sendProcessSignal({
      serverId,
      pid,
      signal: currentAction.signal,
    });

    // Ignore stale responses, for example if the dialog was closed meanwhile.
    if (token !== requestToken) return;

    loading = false;
    if (requestError) {
      error = unWrapErrorMessage(
        requestError,
        `Failed to send ${currentAction.signal}`,
      );
      return;
    }

    toast.success(
      `${currentAction.signal} sent to ${processName} (PID ${pid})`,
    );
    open = false;
  }

  $effect(() => {
    const currentAction = action;
    if (!open || !currentAction) return;
    untrack(() => {
      void run(currentAction);
    });
  });

  function reset() {
    requestToken++;
    loading = false;
    error = "";
  }

  function retry() {
    if (!action) return;
    void run(action);
  }
</script>

<AlertDialog.Root
  bind:open
  onOpenChange={(value) => {
    if (!value) reset();
  }}
>
  <AlertDialog.Content>
    {#if loading}
      <AlertDialog.Header>
        <AlertDialog.Title class="sr-only">{loadingText}</AlertDialog.Title>
      </AlertDialog.Header>
      <Loading class="h-50" text={loadingText} loading={loading} />
    {:else if error !== ""}
      <AlertDialog.Header>
        <AlertDialog.Title>
          {action?.label ?? "Action"} failed
        </AlertDialog.Title>
        <AlertDialog.Description>
          {error}
        </AlertDialog.Description>
      </AlertDialog.Header>
      <AlertDialog.Footer>
        <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
        <Button onclick={retry}>Try again</Button>
      </AlertDialog.Footer>
    {:else}
      <AlertDialog.Header>
        <AlertDialog.Title>{action?.label ?? "Process action"}</AlertDialog.Title>
        <AlertDialog.Description>
          Ready to send {action?.signal} to {processName} (PID {pid}).
        </AlertDialog.Description>
      </AlertDialog.Header>
    {/if}
  </AlertDialog.Content>
</AlertDialog.Root>
