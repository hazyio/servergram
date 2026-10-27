<script lang="ts">
  import * as Field from "#lib/components/ui/field/index.ts";
  import Card from "#lib/components/ui/card/card.svelte";
  import Input from "#lib/components/ui/input/input.svelte";
  import Button from "#lib/components/ui/button/button.svelte";
  import { NewRequestClient } from "#lib/api/request.ts";
  import Loader from "@lucide/svelte/icons/loader";

  let password = $state("");
  let error = $state("");
  let submitting = $state(false);

  async function handleLogin(event: SubmitEvent) {
    event.preventDefault();
    error = "";
    submitting = true;
    try {
      const { data, requestError } = await NewRequestClient().POST(
        "/api/account/login",
        { password },
      );
    } catch {
      error = "The password was not accepted.";
    } finally {
      submitting = false;
    }
  }
</script>

<div class="h-dvh flex justify-center items-center">
  <Card class="max-w-sm mt-12 md:mt-20 w-[90%] md:min-w-120 p-4 px-5">
    <form onsubmit={handleLogin}>
      <div class="mb-4 text-lg font-semibold">Sign in</div>
      <Field.Field>
        <Field.Label for="password">Password</Field.Label>
        <Input
          disabled={submitting}
          id="password"
          type="password"
          autocomplete="current-password"
          bind:value={password}
          required
        />
        {#if error}
          <Field.Error>{error}</Field.Error>
        {/if}
      </Field.Field>
      <Button class="mt-4 w-full" type="submit" disabled={submitting}>
        {#if submitting}
          <Loader class="animate-spin" />
        {:else}
          Sign In
        {/if}
      </Button>
    </form>
  </Card>
</div>
