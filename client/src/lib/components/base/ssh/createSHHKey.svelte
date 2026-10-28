<script lang="ts">
  import * as Dialog from "#lib/components/ui/dialog/index.ts";
  import Textarea from "#lib/components/ui/textarea/textarea.svelte";
  import * as z from "zod";
  import { Button } from "../../ui/button";
  import * as Field from "../../ui/field";
  import Input from "#lib/components/ui/input/input.svelte";
  import { NewRequestClient } from "#lib/api/request.ts";
  import { toast } from "svelte-sonner";
  import Loader from "@lucide/svelte/icons/loader";
  let {
    text,
    onCreated,
  }: {
    text: string;
    onCreated?: () => void;
  } = $props();
  let open = $state(false);
  let creating = $state(false);
  let sshName = $state("");
  let sshKey = $state("");
  let formError = $state({
    key: "",
    name: "",
  });
  const PRIVATE_KEY_RE =
    /^-----BEGIN (OPENSSH|RSA|EC|DSA)? ?PRIVATE KEY-----[\s\S]+-----END \1? ?PRIVATE KEY-----$/;

  const formData = z.object({
    name: z
      .string()
      .trim() // trim first so min() sees the trimmed value
      .min(1, "Name is required")
      .max(100, "Name is too long")
      .regex(/^[\w .@-]+$/, "Use letters, numbers, spaces, . _ - @ only"),

    sshKey: z
      .string()
      .trim()
      .min(1, "SSH key is required")
      .refine((k) => !/^(ssh-|ecdsa-|sk-)/.test(k), {
        message: "That looks like a public key. Paste the private key.",
      })
      .refine((k) => PRIVATE_KEY_RE.test(k.replace(/\r\n/g, "\n")), {
        message: "Not a valid private key (expected a PEM/OpenSSH block)",
      }),
  });

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    formError.key = "";
    formError.name = "";

    const result = formData.safeParse({ name: sshName, sshKey });
    if (!result.success) {
      if (result.error.issues[0]?.path[0] === "sshKey") {
        formError.key = result.error.issues[0]?.message ?? "Invalid key";
      }
      if (result.error.issues[0]?.path[0] === "name") {
        formError.name = result.error.issues[0]?.message ?? "Invalid name";
      }
      return;
    }
    creating = true;
    const { error } = await NewRequestClient().POST("/api/ssh/add", {
      body: {
        name: sshName,
        key: sshKey,
      },
    });
    creating = false;

    if (error) {
      toast.error(
        error.detail ??
          error.title ??
          "An error occurred while adding the SSH key",
      );
      return;
    }
    toast.success("SSH key added successfully");
    sshName = "";
    sshKey = "";
    open = false;
    onCreated?.();
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger><Button>{text}</Button></Dialog.Trigger>
  <Dialog.Content class="z-60">
    <Dialog.Header>
      <Dialog.Title>Add SSH Key</Dialog.Title>
    </Dialog.Header>
    <div>
      <form
        id="create-ssh-key-form"
        onsubmit={handleSubmit}
        class="flex flex-col gap-4"
      >
        <Field.Field data-invalid={!!formError.name}>
          <Field.Label for="name">Name</Field.Label>
          <Input
            id="name"
            type="text"
            bind:value={sshName}
            oninput={() => (formError.name = "")}
            required
            aria-invalid={formError.name ? "true" : "false"}
            disabled={creating}
          />

          {#if formError.name}
            <Field.Error>{formError.name}</Field.Error>
          {/if}
        </Field.Field>
        <Field.Field data-invalid={!!formError.key}>
          <Field.Label for="key">Private Key</Field.Label>
          <Textarea
            id="key"
            class="min-h-70 max-h-[50vh] w-full resize-none field-sizing-fixed overflow-y-auto"
            bind:value={sshKey}
            oninput={() => (formError.key = "")}
            required
            aria-invalid={formError.key ? "true" : "false"}
            disabled={creating}
          />

          {#if formError.key}
            <Field.Error>{formError.key}</Field.Error>
          {/if}
        </Field.Field>
      </form>
    </div>
    <Dialog.Footer>
      <Button type="submit" form="create-ssh-key-form" disabled={creating}>
        {#if creating}
          <Loader class="animate-spin" />
        {:else}
          Save
        {/if}
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
