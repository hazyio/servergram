<script lang="ts">
  import { NewRequestClient, unWrapErrorMessage } from "#lib/api/request.ts";
  import SelectGroup from "#lib/components/base/groups/selectGroup.svelte";
  import SelectSSHKey from "#lib/components/base/ssh/selectSSHKey.svelte";
  import { Button } from "#lib/components/ui/button/index.ts";
  import Card from "#lib/components/ui/card/card.svelte";
  import * as Field from "#lib/components/ui/field/index.ts";
  import { Input } from "#lib/components/ui/input/index.ts";
  import Loader from "@lucide/svelte/icons/loader";
  import { toast } from "svelte-sonner";
  import * as z from "zod";

  let serverName = $state("");
  let host = $state("");
  let username = $state("");
  let port = $state(22);
  let group = $state(0);
  let sshKey = $state(0);
  let formError = $state({
    serverName: "",
    host: "",
    username: "",
    port: "",
    group: "",
    sshKey: "",
  });
  let submitting = $state(false);
  const hostnamePattern =
    /^(?=.{1,253}\.?$)(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)(?:\.(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?))*\.?$/;

  function isValidHost(value: string) {
    const normalized =
      value.startsWith("[") && value.endsWith("]") ? value.slice(1, -1) : value;
    if (/^(?:\d{1,3}\.){3}\d{1,3}$/.test(normalized)) {
      return z.ipv4().safeParse(normalized).success;
    }
    return (
      z.ipv6().safeParse(normalized).success || hostnamePattern.test(normalized)
    );
  }

  const formData = z.object({
    serverName: z
      .string()
      .min(1, "Server name is required")
      .max(100, "Server name is too long")
      .trim(),
    host: z
      .string()
      .trim()
      .min(1, "Host is required")
      .refine(isValidHost, "Enter a valid hostname or IP address"),
    username: z
      .string()
      .trim()
      .min(1, "SSH username is required")
      .max(255, "SSH username is too long"),
    port: z
      .number()
      .int("Port must be a whole number")
      .min(1, "Port must be at least 1")
      .max(65535, "Port must not exceed 65535"),
    group: z.number().min(1, "Group is required"),
    sshKey: z.number().min(1, "SSH key is required"),
  });

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    formError.serverName = "";
    formError.host = "";
    formError.username = "";
    formError.port = "";
    formError.group = "";
    formError.sshKey = "";

    const result = formData.safeParse({
      serverName,
      host,
      username,
      port,
      group,
      sshKey,
    });
    if (!result.success) {
      for (const issue of result.error.issues) {
        const field = issue.path[0];
        if (field === "serverName") formError.serverName = issue.message;
        if (field === "host") formError.host = issue.message;
        if (field === "username") formError.username = issue.message;
        if (field === "port") formError.port = issue.message;
        if (field === "group") formError.group = issue.message;
        if (field === "sshKey") formError.sshKey = issue.message;
      }
      return;
    }

    submitting = true;
    const { data, error } = await NewRequestClient().POST("/api/servers/add", {
      body: {
        group_id: group,
        host: host,
        name: serverName,
        port: port,
        ssh_key_id: sshKey,
        username: username,
      },
    });
    submitting = false;
    if (error) {
      toast.error(unWrapErrorMessage(error, "Failed to add server"));
      return;
    }
    toast.success("Server added successfully");
    submitting = false;
    serverName = "";
    host = "";
    username = "";
    port = 22;
    group = 0;
    sshKey = 0;
  }
</script>

<svelte:head>
  <title>Add Server | ServerGram</title>
  <meta
    name="description"
    content="Add a server with its hostname, SSH username, port, group, and SSH key."
  />
</svelte:head>

<div class="container mx-auto p-3">
  <Card class="mx-auto w-full mt-2 max-w-xl p-4">
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
      <Field.Field data-invalid={!!formError.host}>
        <Field.Label for="server-host">Host</Field.Label>
        <Input
          id="server-host"
          type="text"
          bind:value={host}
          oninput={() => (formError.host = "")}
          placeholder="example.com or IP address"
          required
          aria-invalid={formError.host ? "true" : "false"}
        />
        {#if formError.host}
          <Field.Error>{formError.host}</Field.Error>
        {/if}
      </Field.Field>
      <Field.Field data-invalid={!!formError.username}>
        <Field.Label for="ssh-username">SSH Username</Field.Label>
        <Input
          disabled={submitting}
          id="ssh-username"
          type="text"
          bind:value={username}
          oninput={() => (formError.username = "")}
          autocomplete="username"
          required
          aria-invalid={formError.username ? "true" : "false"}
        />
        {#if formError.username}
          <Field.Error>{formError.username}</Field.Error>
        {/if}
      </Field.Field>
      <Field.Field data-invalid={!!formError.port}>
        <Field.Label for="ssh-port">SSH Port</Field.Label>
        <Input
          disabled={submitting}
          id="ssh-port"
          type="number"
          bind:value={port}
          oninput={() => (formError.port = "")}
          min="1"
          max="65535"
          step="1"
          required
          aria-invalid={formError.port ? "true" : "false"}
        />
        {#if formError.port}
          <Field.Error>{formError.port}</Field.Error>
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

      <Button type="submit" form="create-server-form" disabled={submitting}>
        {#if submitting}
          <Loader class="animate-spin" />
        {:else}
          Save
        {/if}
      </Button>
    </form>
  </Card>
</div>
