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
    oldName = $bindable(),
  }: {
    serverId: number;
    dir: string;
    oldName: null | string;
  } = $props();
  let submitting = $state(false);
  let entryName = $state(oldName ?? "");
  let formError = $state({
    name: "",
  });
  const entrySchema = z.object({
    name: z
      .string()
      .trim() // trim first so min() sees the trimmed value
      .min(1, "Name is required")
      .max(500, "Name is too long")
      .regex(/^[\w .@-]+$/, "Use letters, numbers, spaces, . _ - @ only")
      .refine((val) => val !== oldName, {
        message: "New name must be different from current name",
      }),
  });
  async function handleSubmit() {
    if (!oldName) {
      formError.name = "No entry selected";
      return;
    }
    formError.name = "";

    const result = entrySchema.safeParse({ name: entryName });
    if (!result.success) {
      if (result.error.issues[0]?.path[0] === "name") {
        formError.name = result.error.issues[0]?.message ?? "Invalid name";
      }
      return;
    }
    submitting = true;
    const { error } = await NewRequestClient().PATCH(
      "/api/files/rename-path-in-server/{id}",
      {
        body: {
          path: dir,
          old_name: oldName,
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
    oldName = null;
  }
  $effect(() => {
    entryName = oldName ?? "";
  });
</script>

<Dialog.Root open={!!oldName}>
  <Dialog.Content
    showCloseButton={false}
    onEscapeKeydown={(event) => event.preventDefault()}
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
      <Button
        disabled={submitting}
        onclick={() => {
          entryName = "";
          formError.name = "";
          oldName = null;
        }}>Cancel</Button
      >
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
