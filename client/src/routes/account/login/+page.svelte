<script lang="ts">
  import { NewRequestClient } from "#lib/api/request.ts";
  import Button from "#lib/components/ui/button/button.svelte";
  import Card from "#lib/components/ui/card/card.svelte";
  import * as Field from "#lib/components/ui/field/index.ts";
  import Input from "#lib/components/ui/input/input.svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import Eye from "@lucide/svelte/icons/eye";
  import EyeOff from "@lucide/svelte/icons/eye-off";
  import Loader from "@lucide/svelte/icons/loader";

  let password = $state("");
  let formError = $state("");
  let submitting = $state(false);
  let showPassword = $state(false);

  async function handleLogin(event: SubmitEvent) {
    event.preventDefault();
    formError = "";
    submitting = true;
    const { data, error } = await NewRequestClient().POST(
      "/api/account/login",
      {
        body: {
          password,
        },
      },
    );
    submitting = false;
    if (error) {
      formError = error.detail ?? "An unknown error occurred";
      return;
    }
    // Only same-origin paths, to avoid open redirects.
    const redirect = page.url.searchParams.get("redirect");
    const target =
      redirect && redirect.startsWith("/") && !redirect.startsWith("//")
        ? redirect
        : "/";
    await goto(target);
  }
</script>

<div class="h-dvh flex justify-center items-center">
  <Card class="max-w-sm w-[90%] md:min-w-120 p-4 px-5">
    <form onsubmit={handleLogin}>
      <div class="mb-4 text-lg font-semibold">
        This resource is protected, enter password to unlock
      </div>
      <Field.Field>
        <Field.Label for="password">Password</Field.Label>
        <div class="relative">
          <Input
            disabled={submitting}
            id="password"
            type={showPassword ? "text" : "password"}
            autocomplete="current-password"
            bind:value={password}
            required
          />
          <Button
            type="button"
            variant="ghost"
            size="icon"
            class="absolute top-0 right-0 h-full px-3 hover:bg-transparent"
            onclick={() => (showPassword = !showPassword)}
          >
            {#if showPassword}
              <EyeOff class="h-4 w-4" />
            {:else}
              <Eye class="h-4 w-4" />
            {/if}
          </Button>
        </div>

        {#if formError}
          <Field.Error>{formError}</Field.Error>
        {/if}
      </Field.Field>
      <Button class="mt-6 w-full" type="submit" disabled={submitting}>
        {#if submitting}
          <Loader class="animate-spin" />
        {:else}
          Sign In
        {/if}
      </Button>
    </form>
  </Card>
</div>
