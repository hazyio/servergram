<script lang="ts">
  import * as Dialog from "#lib/components/ui/dialog/index.ts";
  import * as z from "zod";
  import { Button } from "../../ui/button";
  import * as Field from "../../ui/field";
  import { Input } from "../../ui/input";
  import SelectGroup from "../groups/selectGroup.svelte";
  import SelectSSHKey from "../ssh/selectSSHKey.svelte";

  let serverName = $state("");
  let group = $state(0);
  let sshKey = $state(0);
  let formError = $state({
    serverName: "",
    group: "",
    sshKey: "",
  });
  const formData = z.object({
    serverName: z
      .string()
      .min(1, "Server name is required")
      .max(100, "Server name is too long")
      .trim(),
    group: z.string().min(1, "Group is required"),
    sshKey: z.number().min(1, "SSH key is required"),
  });

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    formError.serverName = "";
    formError.group = "";
    formError.sshKey = "";

    const result = formData.safeParse({ serverName, group, sshKey });
    if (!result.success) {
      if (result.error.issues[0]?.path[0] === "serverName") {
        formError.serverName =
          result.error.issues[0]?.message ?? "Invalid server name";
      }
      if (result.error.issues[0]?.path[0] === "group") {
        formError.group = result.error.issues[0]?.message ?? "Invalid group";
      }
      if (result.error.issues[0]?.path[0] === "sshKey") {
        formError.sshKey = result.error.issues[0]?.message ?? "Invalid SSH key";
      }
      return;
    }
  }
</script>

<Dialog.Root>
  <Dialog.Trigger><Button>Add Server</Button></Dialog.Trigger>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Add a New Server</Dialog.Title>
    </Dialog.Header>
    <div>
      <form
        id="create-server-form"
        onsubmit={handleSubmit}
        class="flex flex-col gap-4"
      >
        <Field.Field data-invalid={!!formError.serverName}>
          <Field.Label for="server-name">Server Name</Field.Label>
          <Input
            id="server-name"
            type="text"
            bind:value={serverName}
            oninput={() => (formError.serverName = "")}
            required
            aria-invalid={formError.serverName ? "true" : "false"}
          />

          {#if formError.serverName}
            <Field.Error>{formError.serverName}</Field.Error>
          {/if}
        </Field.Field>
        <Field.Field data-invalid={!!formError.group}>
          <Field.Label for="group">Group</Field.Label>
          <SelectGroup
            value={group}
            onSelect={(val) => {
              group = val;
              formError.group = "";
            }}
          />
          <input type="hidden" name="group" value={group} />

          {#if formError.group}
            <Field.Error>{formError.group}</Field.Error>
          {/if}
        </Field.Field>
        <Field.Field data-invalid={!!formError.sshKey}>
          <Field.Label for="ssh-key">SSH Key</Field.Label>
          <SelectSSHKey
            value={sshKey}
            onSelect={(val) => {
              sshKey = val;
              formError.sshKey = "";
            }}
          />
          <input type="hidden" name="sshKey" value={sshKey} />

          {#if formError.sshKey}
            <Field.Error>{formError.sshKey}</Field.Error>
          {/if}
        </Field.Field>
      </form>
    </div>
    <Dialog.Footer>
      <Button type="submit" form="create-server-form">Save</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
