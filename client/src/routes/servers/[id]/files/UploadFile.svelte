<script lang="ts">
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Field from "#lib/components/ui/field/index.ts";
  import * as z from "zod";

  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import Loading from "#lib/components/base/loading.svelte";
  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import Input from "#lib/components/ui/input/input.svelte";
  let {
    serverId,
    dir,
    open = $bindable(false),
  }: {
    serverId: number;
    dir: string;
    open: boolean;
  } = $props();
  let submitting = $state(false);
  let formError = $state({
    name: "",
  });
  let entryName = $state("");

  const entrySchema = z.object({
    name: z
      .string()
      .trim() // trim first so min() sees the trimmed value
      .min(1, "Name is required")
      .max(500, "Name is too long")
      .regex(/^[\w .@-]+$/, "Use letters, numbers, spaces, . _ - @ only"),
  });
  async function handleSubmit() {
    formError.name = "";

    const result = entrySchema.safeParse({ name: entryName });
    if (!result.success) {
      if (result.error.issues[0]?.path[0] === "name") {
        formError.name = result.error.issues[0]?.message ?? "Invalid name";
      }
      return;
    }
    submitting = true;
    const { error } = await NewRequestClient().POST(
      "/api/files/upload-file-to-server/{id}",
      {
        body: {
          path: dir,
          new_name: entryName,
        },
        params: {
          path: {
            id: serverId,
          },
        },
      },
    );

    submitting = false;
    if (error) {
      formError.name = unWrapErrorMessage(error, "Failed to rename entry");
      return;
    }
    entryName = "";
  }
  function close() {
    entryName = "";
    formError.name = "";
    open = false;
  }
</script>

<Dialog.Root open={open}>
  <Dialog.Content
    showCloseButton={false}
    onEscapeKeydown={() => {
      close();
    }}
    onInteractOutside={(event) => event.preventDefault()}
  >
    <Dialog.Header>
      <Dialog.Title>Rename Entry</Dialog.Title>
    </Dialog.Header>
    <div>
      <form
        id="create-emtry-form"
        onsubmit={handleSubmit}
        class="flex flex-col gap-4"
      >
        <Field.Field data-invalid={!!formError.name}>
          <Field.Label for="name">Name</Field.Label>
          <Input
            id="name"
            placeholder="Test"
            disabled={submitting}
            type={"text"}
            bind:value={entryName}
            aria-invalid={formError.name ? "true" : "false"}
            oninput={() => {
              formError.name = "";
            }}
          />

          {#if formError.name}
            <Field.Error>{formError.name}</Field.Error>
          {/if}
        </Field.Field>
      </form>
    </div>

    <Dialog.Footer>
      <Button disabled={submitting} onclick={close}>Cancel</Button>
      <Button
        disabled={submitting}
        onclick={() => {
          handleSubmit();
        }}
      >
        {#if submitting}
          <Loading loading />
        {:else} 
          Rename
        {/if}
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
