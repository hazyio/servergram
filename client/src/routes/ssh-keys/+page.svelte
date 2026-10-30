<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { NewRequestClient } from "#lib/api/request.ts";
  import type { components } from "#lib/types/request-schema.d.ts";
  import Error from "#lib/components/base/error.svelte";
  import CreateSSHKey from "#lib/components/base/ssh/createSHHKey.svelte";
  import Loading from "#lib/components/base/loading.svelte";
  import NoItemsFound from "#lib/components/base/noItemsFound.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import RefreshCw from "@lucide/svelte/icons/refresh-cw";
  import { onMount } from "svelte";
  import appPaths from "#lib/appPaths.ts";
  import ViewKeyInfo from "./viewKeyInfo.svelte";
  import DeleteRequestDialog from "#lib/components/base/DeleteRequestDialog.svelte";
  let loading = $state(true);
  let sshKeys: components["schemas"]["SSHKeyListResponseBody"] | undefined =
    $state();
  let search = $state(page.url.searchParams.get("search") ?? "");
  let filteredSSHKeys = $derived(
    sshKeys?.items?.filter((key) =>
      key.name.toLowerCase().includes(search.trim().toLowerCase()),
    ) ?? [],
  );
  let loadError = $state(false);
  function updateSearch(event: Event) {
    const value = (event.currentTarget as HTMLInputElement).value;
    const url = new URL(page.url.href);
    if (value) {
      url.searchParams.set("search", value);
    } else {
      url.searchParams.delete("search");
    }
    goto(url, { state: page.state, shallow: true, replace: true });
  }
  async function loadSSHKeys() {
    loading = true;
    const { data, error } = await NewRequestClient().GET("/api/ssh-key/list");
    loading = false;
    if (error) {
      loadError = true;
      return;
    }
    loadError = false;
    sshKeys = data;
  }
  onMount(loadSSHKeys);
</script>

<svelte:head>
  <title>SSH Keys | ServerGram</title>
  <meta
    name="description"
    content="Manage SSH keys used to connect to your ServerGram servers."
  />
</svelte:head>

<div class="p-2">
  <div class="flex gap-2 justify-center items-center">
    <Input
      placeholder="Search SSH keys by name..."
      bind:value={search}
      oninput={updateSearch}
    />
    <Button variant="secondary" onclick={loadSSHKeys}><RefreshCw /></Button>
    <CreateSSHKey text="Add SSH Key" onCreated={loadSSHKeys} />
  </div>
  <Loading class="mt-20" {loading} text="Loading..." />
  {#if (loadError || !sshKeys) && !loading}
    <Error
      title="Failed to load SSH keys"
      message="An error occurred while fetching the SSH key list."
      onRetry={loadSSHKeys}
    />
  {:else if sshKeys && !loading}
    {#if sshKeys.items?.length && filteredSSHKeys.length}
      <ul class="mt-3 grid md:grid-cols-2 lg:grid-cols-3 gap-2">
        {#each filteredSSHKeys as key (key.id)}
          <li class="bg-secondary rounded-md mt-2 py-3 px-2 flex flex-row">
            <div class="flex gap-2 flex-col grow">
              <span class="font-bold">{key.name}</span>
              <small
                >Encrypted: <span
                  class={key.encrypted ? "text-green-500" : "text-amber-500"}
                  >{key.encrypted ? "Yes" : "No"}</span
                ></small
              >
              <small>
                <a href={appPaths.sshKeys.resources(key.id)}>
                  Resources: {key.resources_attached}</a
                ></small
              >
              <small
                >Created at: {new Date(key.created_at).toLocaleString()}</small
              >
            </div>
            <div class="flex justify-start gap-2 flex-col">
              <ViewKeyInfo id={key.id} />
              <DeleteRequestDialog
                id={key.id}
                name={key.name}
                onDeleted={loadSSHKeys}
                type="ssh-key"
              />
            </div>
          </li>
        {/each}
      </ul>
    {:else if sshKeys.items?.length}
      <NoItemsFound
        title="No matching SSH keys"
        message="Try a different name."
      />
    {:else}
      <NoItemsFound
        title="No SSH keys yet"
        message="Add an SSH key to manage your servers."
      />
    {/if}
  {/if}
</div>
