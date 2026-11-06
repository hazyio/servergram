<script lang="ts" module>
  export type NewEntryType = "none" | "file" | "directory";
</script>

<script lang="ts">
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as z from "zod";
  import * as Field from "#lib/components/ui/field/index.ts";

  import * as Dialog from "#lib/components/ui/dialog/index.js";
  import Input from "#lib/components/ui/input/input.svelte";
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import Loading from "#lib/components/base/loading.svelte";
  let {
    serverId,
    dir,
    type = $bindable(),
  }: {
    serverId: number;
    dir: string;
    type: NewEntryType;
  } = $props();
  let submitting = $state(false);
  let entryName = $state("");
  let formError = $state({
    name: "",
  });
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
    const { data, error, response } = await NewRequestClient().PUT(
      "/api/files/create-dir-in-server/{id}",
      {
        body: {
          name: entryName,
          path: dir,
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
      formError.name = unWrapErrorMessage(
        error,
        `Failed to create ${type === "directory" ? "directory" : "file"}`,
      );
      return;
    }
    entryName = "";
    type = "none";
  }
</script>

<Dialog.Root open={type !== "none"}>
  <Dialog.Content
    showCloseButton={false}
    onEscapeKeydown={(event) => event.preventDefault()}
    onInteractOutside={(event) => event.preventDefault()}
  >
    <Dialog.Header>
      <Dialog.Title
        >Create New {type === "directory" ? "Directory" : "File"}</Dialog.Title
      >
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
          type = "none";
          formError.name = "";
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
          Create {type === "directory" ? "Directory" : "File"}
        {/if}
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
