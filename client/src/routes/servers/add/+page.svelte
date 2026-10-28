<script lang="ts">
  import SelectGroup from "#lib/components/base/groups/selectGroup.svelte";
  import SelectSSHKey from "#lib/components/base/ssh/selectSSHKey.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import * as Field from "#lib/components/ui/field/index.ts";
  import { Input } from "#lib/components/ui/input/index.ts";
  import * as z from "zod";

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
    group: z.number().min(1, "Group is required"),
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
  <div class="container mx-auto p-4">
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
              <Button type="submit" form="create-server-form">Save</Button>

      </form>
    </div>
