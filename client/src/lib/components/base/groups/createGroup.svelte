<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import * as Dialog from "#lib/components/ui/dialog/index.ts";
  import Input from "#lib/components/ui/input/input.svelte";
  import Loader from "@lucide/svelte/icons/loader";
  import { toast } from "svelte-sonner";
  import * as z from "zod";
  import { Button } from "../../ui/button";
  import * as Field from "../../ui/field";
  let { text } = $props();
  let open = $state(false);
  let creating = $state(false);
  let groupName = $state("");
  let formError = $state({
    name: "",
  });

  const formData = z.object({
    name: z
      .string()
      .trim() // trim first so min() sees the trimmed value
      .min(1, "Name is required")
      .max(100, "Name is too long")
      .regex(/^[\w .@-]+$/, "Use letters, numbers, spaces, . _ - @ only"),
  });

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    formError.name = "";

    const result = formData.safeParse({ name: groupName });
    if (!result.success) {
      if (result.error.issues[0]?.path[0] === "name") {
        formError.name = result.error.issues[0]?.message ?? "Invalid name";
      }
      return;
    }
    creating = true;
    const { error } = await NewRequestClient().POST("/api/group/create", {
      body: {
        name: groupName,
      },
    });
    creating = false;

    if (error) {
      toast.error(error.title ?? "An error occurred while creating the group");
      return;
    }
    toast.success("Group created successfully");
    groupName = "";
    open = false;
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger><Button>{text}</Button></Dialog.Trigger>
  <Dialog.Content class="z-60">
    <Dialog.Header>
      <Dialog.Title>Create new group</Dialog.Title>
    </Dialog.Header>
    <div>
      <form
        id="create-group-form"
        onsubmit={handleSubmit}
        class="flex flex-col gap-4"
      >
        <Field.Field data-invalid={!!formError.name}>
          <Field.Label for="name">Name</Field.Label>
          <Input
            id="name"
            type="text"
            bind:value={groupName}
            oninput={() => (formError.name = "")}
            required
            aria-invalid={formError.name ? "true" : "false"}
            disabled={creating}
          />

          {#if formError.name}
            <Field.Error>{formError.name}</Field.Error>
          {/if}
        </Field.Field>
      </form>
    </div>
    <Dialog.Footer>
      <Button type="submit" form="create-group-form" disabled={creating}>
        {#if creating}
          <Loader class="animate-spin" />
        {:else}
          Save
        {/if}
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
